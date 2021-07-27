package query

import (
	"fmt"
	"strings"

	"github.com/cpekyaman/goits/framework/orm/metadata"
	"github.com/cpekyaman/goits/framework/orm/domain"
)

var queryDefRegistry map[string]QueryDef

func init() {
	queryDefRegistry = make(map[string]QueryDef)
}

const (
	findAllTemplate                  = "select %s from %s order by %s"
	findAllPagedTemplate             = findAllTemplate + " limit %d offset %d"
	findOneByAttributeTemplate       = "select %s from %s where %s = $1"
	findAllByAttributesTemplate      = "select %s from %s where %s order by %s"
	findAllByAttributesPagedTemplate = findAllByAttributesTemplate + " limit %d offset %d"
)

// GetQueryDef returns an already registered QueryDef for the entity represented by provided metadata.
func GetQueryDef(ed metadata.EntityDef) (QueryDef, bool) {
	cm, found := queryDefRegistry[ed.Name()]
	return cm, found
}

// BuildQueryDef builds the default static sql statements that will be used by a repository.
func BuildQueryDef(introspect interface{}, ed metadata.EntityDef, cm metadata.ColumnMapper) QueryDef {
	selectColumns := strings.Join(cm.Columns(), ", ")

	qd := sqlQueryDef{
		selectColumns: selectColumns,
		findOne:       fmt.Sprintf(findOneByAttributeTemplate, selectColumns, ed.FullTableName(), ed.PKColumn()),
		findAll:       fmt.Sprintf(findAllTemplate, selectColumns, ed.FullTableName(), ed.DefaultSort()),
		insert:        generateInsertStatement(ed.Schema(), ed.Table(), cm),
		update:        generateUpdateStatement(ed.Schema(), ed.Table(), cm, introspect),
		delete:        generateDeleteStatement(ed, introspect),
	}

	queryDefRegistry[ed.Name()] = qd

	return qd
}

// generateInsertStatement creates the insert sql statement.
func generateInsertStatement(schema string, table string, cm metadata.ColumnMapper) string {
	var columns []string
	var values []string
	for _, f := range cm.Fields() {
		if domain.IsNonInsertableField(f) {
			continue
		}
		if cm.HasColumn(f) {
			col := cm.Column(f)

			columns = append(columns, col)
			values = append(values, ":"+col)
		}
	}

	columnsPart := strings.Join(columns, ", ")
	valuesPart := strings.Join(values, ", ")

	return fmt.Sprintf("insert into %s.%s(%s) values(%s)", schema, table, columnsPart, valuesPart)
}

// generateUpdateStatement creates appropriate update-all statement that updates all fields.
// The generated statement also contains versioning / timestamping if the target type supports those.
func generateUpdateStatement(schema string, table string, cm metadata.ColumnMapper, introspect interface{}) string {
	var stmt []string
	for _, f := range cm.Fields() {
		if domain.IsNonUpdatableField(f) {
			continue
		}
		
		if cm.HasColumn(f) {
			col := cm.Column(f)
			stmt = append(stmt, fmt.Sprintf("%s=:%s", col, col))
		}
	}

	where := " where id=:id"

	// versioned entities have both update part and where condition different
	_, ok := introspect.(domain.Versioned)
	if ok {
		where = " where id=:id AND version=:version"
		stmt = append(stmt, "version=version + 1")
	}

	// if timestamped, we also need to update modify time
	_, ok = introspect.(domain.Timestamped)
	if ok {
		stmt = append(stmt, "last_modified_time=now()")
	}

	return fmt.Sprintf("update %s.%s set %s %s", schema, table, strings.Join(stmt, ", "), where)
}

// generateDeleteStatement builds the default delete statement for the entity.
// It generates a delete or update statement depending on whether the entity uses soft delete or not.
func generateDeleteStatement(ed metadata.EntityDef, introspect interface{}) string {
	if ed.SoftDelete() {
		_, ok := introspect.(domain.Timestamped)
		if ok {
			return fmt.Sprintf("update %s.%s set deleted=true, last_modified_time=now() where %s = $1", ed.Schema(), ed.Table(), ed.PKColumn())
		} else {
			return fmt.Sprintf("update %s.%s set deleted=true where %s = $1", ed.Schema(), ed.Table(), ed.PKColumn())
		}
	} else {
		return fmt.Sprintf("delete from %s.%s where %s = $1", ed.Schema(), ed.Table(), ed.PKColumn())
	}
}

// QueryDef is the metadata of pre-generated queries for an entity.
type QueryDef interface {
	FindOne() string
	FindAll() string
	Insert() string
	Update() string
	Delete() string
	SelectColumns() string
}

// sqlQueryDef provides statically defined sql statements for a repository.
type sqlQueryDef struct {
	findOne       string
	findAll       string
	insert        string
	update        string
	delete        string
	selectColumns string
}

func (this sqlQueryDef) FindOne() string {
	return this.findOne
}
func (this sqlQueryDef) FindAll() string {
	return this.findAll
}
func (this sqlQueryDef) Insert() string {
	return this.insert
}
func (this sqlQueryDef) Update() string {
	return this.update
}
func (this sqlQueryDef) Delete() string {
	return this.delete
}
func (this sqlQueryDef) SelectColumns() string {
	return this.selectColumns
}