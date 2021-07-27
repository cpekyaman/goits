package repository

import (
	"context"
	"time"

	"github.com/cpekyaman/goits/framework/monitoring"
	"github.com/cpekyaman/goits/framework/orm/domain"
	"github.com/cpekyaman/goits/framework/orm/metadata"
	"github.com/cpekyaman/goits/framework/orm/query"

	"github.com/jmoiron/sqlx"
)

var hist monitoring.HistogramBundle
var cnt monitoring.CounterBundle

func init() {
	labels := []string{"type", "query"}
	hist = monitoring.NewHistogram("db_query_duration_ms", labels)
	cnt = monitoring.NewCounter("db_query_executions", labels)
}

// SqlRepository is the Repository implementation for sql db.
type SqlRepository struct {
	db *sqlx.DB
	ed metadata.EntityDef
	qd query.QueryDef
	cm metadata.ColumnMapper
}

func (this SqlRepository) FindOneById(ctx context.Context, dest interface{}, id uint64) error {
	defer this.log(ctx, "FindOneById", time.Now())
	return this.db.GetContext(ctx, dest, this.qd.FindOne(), id)
}

func (this SqlRepository) FindAll(ctx context.Context, dest interface{}) error {
	defer this.log(ctx, "FindAll", time.Now())
	return this.db.SelectContext(ctx, dest, this.qd.FindAll())
}

func (this SqlRepository) FindAllPaged(ctx context.Context, dest interface{}, limit uint, offset uint64) error {
	defer this.log(ctx, "FindAllPaged", time.Now())
	return this.db.SelectContext(ctx, dest, query.BuildFindAllPagedQuery(this.ed, this.qd, limit, offset))
}

func (this SqlRepository) FindOneByAttribute(ctx context.Context, dest interface{}, attr string, bindval interface{}) error {
	defer this.log(ctx, "FindOneByAttribute", time.Now())

	return this.db.GetContext(ctx, dest, query.BuildFindOneQuery(this.ed, this.qd, this.cm, attr), bindval)
}

func (this SqlRepository) FindAllByAttributes(ctx context.Context, dest interface{}, attrs map[string]interface{}) error {
	defer this.log(ctx, "FindAllByAttributes", time.Now())

	q, params := query.BuildQueryByAttributes(this.ed, this.qd, this.cm, attrs, 0, 0)
	return this.db.SelectContext(ctx, dest, q, params)
}

func (this SqlRepository) FindAllByAttributesPaged(ctx context.Context, dest interface{}, attrs map[string]interface{}, limit uint, offset uint64) error {
	defer this.log(ctx, "FindAllByAttributesPaged", time.Now())

	q, params := query.BuildQueryByAttributes(this.ed, this.qd, this.cm, attrs, limit, offset)
	return this.db.SelectContext(ctx, dest, q, params)
}

func (this SqlRepository) Save(ctx context.Context, entity domain.Entity) error {
	var q string
	var qt string
	if entity.GetId() > 0 {
		q = this.qd.Update()
		qt = "Update"
	} else {
		q = this.qd.Insert()
		qt = "Create"
	}

	defer this.log(ctx, qt, time.Now())

	_, err := this.db.NamedExecContext(ctx, q, entity)
	return err
}

func (this SqlRepository) Delete(ctx context.Context, id uint64) error {
	defer this.log(ctx, "Delete", time.Now())
	_, err := this.db.ExecContext(ctx, this.qd.Delete(), id)
	return err
}

func (this SqlRepository) log(ctx context.Context, query string, start time.Time) {
	mctx, ok := monitoring.GetMonitoringContext(ctx)
	if ok {
		duration := time.Since(start)

		mctx.Logger().
			With(monitoring.StrLogField(monitoring.LF_Type, this.ed.Name()),
				monitoring.StrLogField(monitoring.LF_Query, query),
				monitoring.Int64LogField(monitoring.LF_Ms, duration.Milliseconds())).
			Info("query executed")

		lv := map[string]string{
			"type":  this.ed.Name(),
			"query": query,
		}

		hist.With(lv).Record(float64(duration.Milliseconds()))
		cnt.With(lv).Incr()
	}
}
