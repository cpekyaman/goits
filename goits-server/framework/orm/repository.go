package orm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/cpekyaman/goits/framework/monitoring"
)

var hist monitoring.HistogramBundle
var cnt monitoring.CounterBundle

func init() {
	labels := []string{"type", "query"}
	hist = monitoring.NewHistogram("db_query_duration_ms", labels)
	cnt = monitoring.NewCounter("db_query_executions", labels)
}

// ReaderRepository provides common querying functionality for db entities.
type ReaderRepository interface {
	// FindOneById finds a single entity by the given id.
	FindOneById(ctx context.Context, dest interface{}, id uint64) error

	// FindAll finds all entities.
	FindAll(ctx context.Context, dest interface{}) error

	// FindAllPaged finds all entities in the given page offset and limit.
	FindAllPaged(ctx context.Context, dest interface{}, limit uint, offset uint64) error

	// FindOneByAttribute is a generic finder to find a single entity by a unique attribute.
	FindOneByAttribute(ctx context.Context, dest interface{}, attr string, bindval interface{}) error

	// FindAllByAttributes finds all entities that match attrs criteria given.
	FindAllByAttributes(ctx context.Context, dest interface{}, attrs map[string]interface{}) error
}

// WriterRepository provides basic modify functionality for db entities.
type WriterRepository interface {
	Save(ctx context.Context, entity Entity) error
	Delete(ctx context.Context, id uint64) error
}

// Repository is a generic repository definition for standard crud operations.
type Repository interface {
	ReaderRepository
	WriterRepository
}

type queryDef struct {
	findOne string
	findAll string
	insert  string
	update  string
	delete  string
}

// SqlRepository is the Repository implementation for sql db.
type SqlRepository struct {
	orm  DBW
	ed   *EntityDef
	qd   queryDef
	fmap map[string]string
}

// NewRepository creates a new crud repository for given entity.
func NewRepository(ed *EntityDef, introspect interface{}) SqlRepository {
	v := reflect.ValueOf(introspect)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	fieldMap := make(map[string]string)
	if v.Kind() == reflect.Struct {
		buildFieldMap(fieldMap, v)
	}

	qd := queryDef{
		findOne: fmt.Sprintf("select * from %s where %s = $1", ed.TableName(), ed.PKColumn()),
		findAll: fmt.Sprintf("select * from %s order by %s", ed.TableName(), ed.DefaultSort()),
	}

	qd.insert = generateInsertStatement(ed.tableName, fieldMap)
	qd.update = generateUpdateStatement(ed.tableName, fieldMap, introspect)

	if ed.softDelete {
		_, ok := introspect.(Timestamped)
		if ok {
			qd.delete = fmt.Sprintf("update %s set deleted=true, last_modified_time=now() where %s = $1", ed.TableName(), ed.PKColumn())
		} else {
			qd.delete = fmt.Sprintf("update %s set deleted=true where %s = $1", ed.TableName(), ed.PKColumn())
		}
	} else {
		qd.delete = fmt.Sprintf("delete from %s where %s = $1", ed.TableName(), ed.PKColumn())
	}

	return SqlRepository{dbw, ed, qd, fieldMap}
}

// buildFieldMap builds a map of field names and their corresponding column names by reflecting on value.
func buildFieldMap(fieldMap map[string]string, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)

		switch f.Type.Kind() {
		case reflect.Struct:
			if f.Anonymous {
				buildFieldMap(fieldMap, v.Field(i))
			} else {
				addField(f, fieldMap)
			}
		case reflect.Ptr:
			monitoring.RootLogger().Info(fmt.Sprintf("field %s is pointer %v", f.Name, f))
		default:
			addField(f, fieldMap)
		}
	}
}

// addField adds the mapping of given field to fieldMap.
func addField(f reflect.StructField, fieldMap map[string]string) {
	fieldName := f.Name
	colName, found := f.Tag.Lookup("db")
	if found {
		fieldMap[fieldName] = colName
	} else {
		fieldMap[fieldName] = fieldName
	}
}

// generateInsertStatement creates the insert sql statement.
func generateInsertStatement(table string, fieldMap map[string]string) string {
	var columns []string
	var values []string
	for k, v := range fieldMap {
		if defaultNonInsertableColumns[k] {
			continue
		}
		columns = append(columns, v)
		values = append(values, ":"+v)
	}

	columnsPart := strings.Join(columns, ", ")
	valuesPart := strings.Join(values, ", ")

	return "insert into " + table + "(" + columnsPart + ")" + " values (" + valuesPart + ")"
}

// generateUpdateStatement creates appropriate update-all statement that updates all fields.
func generateUpdateStatement(table string, fieldMap map[string]string, introspect interface{}) string {
	var stmt []string
	for k, v := range fieldMap {
		if defaultNonUpdatableColumns[k] {
			continue
		}
		stmt = append(stmt, fmt.Sprintf("%s=:%s", v, v))
	}

	where := " where id=:id"

	// versioned entities have both update part and where condition different
	_, ok := introspect.(Versioned)
	if ok {
		where = " where id=:id AND version=:version"
		stmt = append(stmt, "version=version + 1")
	}

	// if timestamped, we also need to update modify time
	_, ok = introspect.(Timestamped)
	if ok {
		stmt = append(stmt, "last_modified_time=now()")
	}

	return "update " + table + " set " + strings.Join(stmt, ", ") + where
}

func (this SqlRepository) FindOneById(ctx context.Context, dest interface{}, id uint64) error {
	defer this.log(ctx, "FindOneById", time.Now())
	return this.orm.GetDB().GetContext(ctx, dest, this.qd.findOne, id)
}

func (this SqlRepository) FindAll(ctx context.Context, dest interface{}) error {
	defer this.log(ctx, "FindAll", time.Now())
	return this.orm.GetDB().SelectContext(ctx, dest, this.qd.findAll)
}

func (this SqlRepository) FindAllPaged(ctx context.Context, dest interface{}, limit uint, offset uint64) error {
	defer this.log(ctx, "FindAllPaged", time.Now())
	return this.orm.GetDB().SelectContext(ctx, dest, fmt.Sprintf("select * from %s order by %s limit %d offset %d", this.ed.TableName(), this.ed.DefaultSort(), limit, offset))
}

func (this SqlRepository) FindOneByAttribute(ctx context.Context, dest interface{}, attr string, bindval interface{}) error {
	defer this.log(ctx, "FindOneByAttribute", time.Now())
	return this.orm.GetDB().GetContext(ctx, dest, fmt.Sprintf("select * from %s where %s = $1", this.ed.TableName(), attr), bindval)
}

func (this SqlRepository) FindAllByAttributes(ctx context.Context, dest interface{}, attrs map[string]interface{}) error {
	defer this.log(ctx, "FindAllByAttributes", time.Now())

	var criteria []string
	var params []interface{}
	var idx = 1
	for k, v := range attrs {
		criteria = append(criteria, fmt.Sprintf("%s=$%d", k, idx))
		params = append(params, v)
		idx++
	}
	where := strings.Join(criteria, " AND ")

	return this.orm.GetDB().SelectContext(ctx, dest, fmt.Sprintf("select * from %s where %s", this.ed.TableName(), where), params)
}

func (this SqlRepository) Save(ctx context.Context, entity Entity) error {
	var q string
	var qt string
	if entity.GetId() > 0 {
		q = this.qd.update
		qt = "Update"
	} else {
		q = this.qd.insert
		qt = "Create"
	}

	defer this.log(ctx, qt, time.Now())

	_, err := this.orm.GetDB().NamedExecContext(ctx, q, entity)
	return err
}

func (this SqlRepository) Delete(ctx context.Context, id uint64) error {
	defer this.log(ctx, "Delete", time.Now())
	_, err := this.orm.GetDB().ExecContext(ctx, this.qd.delete, id)
	return err
}

func (this SqlRepository) log(ctx context.Context, query string, start time.Time) {
	mctx, ok := monitoring.GetMonitoringContext(ctx)
	if ok {
		duration := time.Since(start)

		mctx.Logger().
			With(monitoring.StrLogField(monitoring.LF_Type, this.ed.name),
				monitoring.StrLogField(monitoring.LF_Query, query),
				monitoring.Int64LogField(monitoring.LF_Ms, duration.Milliseconds())).
			Info("query executed")

		lv := map[string]string{
			"type":  this.ed.name,
			"query": query,
		}

		hist.With(lv).Record(float64(duration.Milliseconds()))
		cnt.With(lv).Incr()
	}
}
