package orm

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type TestEntity struct {
	Id       uint64 `db:"entity_id"`
	Priority uint8  `db:"priority"`
}

const (
	findOneId uint64 = 5
)

var ed *EntityDef
var columns []string

func init() {
	ed = &EntityDef{
		name:        "TestEntity",
		pkColumn:    "entity_id",
		tableName:   "test_entity",
		defaultSort: "priority desc",
	}

	columns = []string{"entity_id", "priority"}
}

func TestFindById_DataFound(t *testing.T) {
	// given
	repo, mock := newRepoWithMock(t)

	_, rows := mockFindOneWithRows(mock)
	rows.AddRow(findOneId, 8)

	// when
	var te TestEntity
	err := repo.FindOneById(context.Background(), &te, findOneId)

	// then
	assert.Nil(t, err, err)
	assert.Equal(t, uint64(5), te.Id, "id is not correct")
	assert.Equal(t, uint8(8), te.Priority, "priority is not correct")
}

func TestFindById_NoDataFound(t *testing.T) {
	// given
	repo, mock := newRepoWithMock(t)

	mockFindOneWithRows(mock)

	// when
	var te TestEntity
	err := repo.FindOneById(context.Background(), &te, findOneId)

	// then
	assert.NotNil(t, err, "should get back an error")
	assert.True(t, strings.Contains(err.Error(), "no rows"), "unexpected error")
	assert.Equal(t, uint64(0), te.Id, "should not set id")
	assert.Equal(t, uint8(0), te.Priority, "should not set priority")
}

func TestFindById_Error(t *testing.T) {
	// given
	repo, mock := newRepoWithMock(t)

	eq, rows := mockFindOneWithRows(mock)
	eq.WillReturnError(fmt.Errorf("for test"))
	rows.AddRow(findOneId, 8)

	// when
	var te TestEntity
	err := repo.FindOneById(context.Background(), &te, findOneId)

	// then
	assert.NotNil(t, err, "should get back an error")
	assert.True(t, strings.Contains(err.Error(), "for test"), "unexpected error")
}

func TestFindAll_DataFound(t *testing.T) {
	// given
	repo, mock := newRepoWithMock(t)

	_, rows := mockFindAllWithRows(mock)
	rows.AddRow(1, 5)
	rows.AddRow(2, 6)
	rows.AddRow(3, 4)

	// when
	var tearr []TestEntity
	err := repo.FindAll(context.Background(), &tearr)

	// then
	assert.Nil(t, err, "should not get error")
	assert.Equal(t, 3, len(tearr), "not all rows are returned")
	for _, v := range tearr {
		assert.NotEmpty(t, v.Id, "has empty struct")
	}
}

func TestFindAll_NoDataFound(t *testing.T) {
	// given
	repo, mock := newRepoWithMock(t)

	mockFindAllWithRows(mock)

	// when
	var tearr []TestEntity
	err := repo.FindAll(context.Background(), &tearr)

	// then
	assert.Nil(t, err, "should not get an error")
	assert.Equal(t, 0, len(tearr), "no rows should be returned")
}

func TestFindAll_Error(t *testing.T) {
	// given
	repo, mock := newRepoWithMock(t)

	eq, rows := mockFindAllWithRows(mock)
	rows.AddRow(1, 5)
	rows.AddRow(2, 6)
	rows.AddRow(3, 4)
	eq.WillReturnError(fmt.Errorf("for test"))

	// when
	var tearr []TestEntity
	err := repo.FindAll(context.Background(), &tearr)

	// then
	assert.NotNil(t, err, "should not back get error")
	assert.True(t, strings.Contains(err.Error(), "for test"), "unexpected error")
	assert.Equal(t, 0, len(tearr), "no rows should be returned")
}

func mockFindOneWithRows(mock sqlmock.Sqlmock) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	eq, rows := mockFinderWithRows("select \\* from test_entity where entity_id = \\$1", mock)
	eq.WithArgs(findOneId)
	return eq, rows
}

func mockFindAllWithRows(mock sqlmock.Sqlmock) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	eq, rows := mockFinderWithRows("select \\* from test_entity order by priority desc", mock)
	return eq, rows
}

func mockFinderWithRows(queryRegex string, mock sqlmock.Sqlmock) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	rows := sqlmock.NewRows(columns)
	eq := mock.ExpectQuery(queryRegex).WillReturnRows(rows)
	return eq, rows
}

func newRepoWithMock(t *testing.T) (Repository, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err, "could not create mock db")
	WithDB(mockDB, "sqlmock")

	t.Cleanup(func() {
		mockDB.Close()
	})

	repo := NewRepository(ed, &TestEntity{})
	return repo, mock
}
