package testlib

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cpekyaman/goits/framework/orm"
	"github.com/stretchr/testify/assert"
)

func NewSqlMock(t *testing.T) sqlmock.Sqlmock {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err, "could not create mock db")
	orm.WithDB(mockDB, "sqlmock")

	t.Cleanup(func() {
		mockDB.Close()
	})

	return mock
}
