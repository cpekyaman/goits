package mocking

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cpekyaman/goits/framework/orm/db"
	"github.com/cpekyaman/goits/framework/orm/domain"
	"github.com/cpekyaman/goits/framework/orm/metadata"
	"github.com/cpekyaman/goits/framework/orm/query"
	"github.com/stretchr/testify/assert"
)

var SqlError = errors.New("sql: error")

// NewSqlMock creates a new Sqlmock instance for a test.
func NewSqlMock(t *testing.T) sqlmock.Sqlmock {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err, "could not create mock db")
	db.WithDB(mockDB, "sqlmock")

	t.Cleanup(func() {
		mockDB.Close()
	})

	return mock
}

// NewQueryMocker returns a new QueryMocker instance for the given entity.
func NewQueryMocker(ed metadata.EntityDef) QueryMocker {
	return QueryMocker{ed}
}

// QueryMocker provides certain helper methods to create mock queries and sql statements.
type QueryMocker struct {
	ed metadata.EntityDef
}

// ExpectFindOne creates an ExpectedQuery that expects a single row select by primary key.
func (this QueryMocker) ExpectFindOne(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	return mock.ExpectQuery("select .* from " + this.ed.FullTableName() + " where " + this.ed.PKColumn() + " = \\$1")
}

// ExpectFindOneByAttr creates an ExpectedQuery that expects a single row select by given attribute.
func (this QueryMocker) ExpectFindOneByAttr(mock sqlmock.Sqlmock, attr string) *sqlmock.ExpectedQuery {
	cm, found := metadata.GetColumnMapper(this.ed)
	if !found {
		return nil
	}
	return mock.ExpectQuery("select .* from " + this.ed.FullTableName() + " where " + cm.Column(attr) + " = \\$1")
}

// ExpectFindAll creates an ExpectedQuery that expects a select all with default ordering.
func (this QueryMocker) ExpectFindAll(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	return mock.ExpectQuery("select .* from " + this.ed.FullTableName() + " order by " + this.ed.DefaultSort())
}

// ExpectFindAllByAttributes creates and ExpectedQuery that expects a select by using given attributes as criteria.
func (this QueryMocker) ExpectFindAllByAttributes(mock sqlmock.Sqlmock, attrs map[string]interface{}) (*sqlmock.ExpectedQuery, []interface{}) {
	cm, found := metadata.GetColumnMapper(this.ed)
	if !found {
		return nil, nil
	}

	qd, found := query.GetQueryDef(this.ed)
	if ! found {
		return nil, nil
	}

	where, params := query.BuildCriteria(this.ed, qd, cm, attrs)
	return mock.ExpectQuery("select .* from " + this.ed.FullTableName() + " " + where + " order by " + this.ed.DefaultSort()), params
}

// ExpectQueryWithRows gets an existing ExpectedQuery and prepares it to return rows with given column structure.
func (this QueryMocker) ExpectQueryWithRows(eq *sqlmock.ExpectedQuery, columns []string) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	rows := sqlmock.NewRows(columns)
	q := eq.WillReturnRows(rows)
	return q, rows
}

// ExpectInsert creates an ExpectedExec that expects the default insert statement and completes successfully.
func (this QueryMocker) ExpectInsert(mock sqlmock.Sqlmock, expectedId int64) *sqlmock.ExpectedExec {
	return this.insertMock(mock).WillReturnResult(sqlmock.NewResult(expectedId, 1))
}

// ExpectInsertError creates an ExpectedExec that expects the default insert statement and fails with error.
func (this QueryMocker) ExpectInsertError(mock sqlmock.Sqlmock) {
	this.insertMock(mock).WillReturnError(SqlError)
}

func (this QueryMocker) insertMock(mock sqlmock.Sqlmock) *sqlmock.ExpectedExec {
	return mock.ExpectExec("insert into " + this.ed.FullTableName() + "(.*) values(.*)")
}

// ExpectUpdate creates an ExpectedExec that expects the default update statement and completes successfully.
func (this QueryMocker) ExpectUpdate(mock sqlmock.Sqlmock, value interface{}) *sqlmock.ExpectedExec {
	return this.updateMock(mock, value).WillReturnResult(sqlmock.NewResult(0, 1))
}

// ExpectUpdateError creates an ExpectedExec that expects the default update statement and fails with error.
func (this QueryMocker) ExpectUpdateError(mock sqlmock.Sqlmock, value interface{}) {
	this.updateMock(mock, value).WillReturnError(SqlError)
}

func (this QueryMocker) updateMock(mock sqlmock.Sqlmock, value interface{}) *sqlmock.ExpectedExec {
	var q string
	_, ok := value.(domain.Versioned)
	if ok {
		q = "update " + this.ed.FullTableName() + " set .* where id=\\? AND version=\\?"
	} else {
		q = "update " + this.ed.FullTableName() + " set .* where id=\\?"
	}

	return mock.ExpectExec(q)
}
