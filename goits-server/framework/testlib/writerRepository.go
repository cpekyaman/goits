package testlib

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cpekyaman/goits/framework/orm"
	"github.com/stretchr/testify/assert"
)

type WriterRepositoryTest struct {
	entityDef         orm.EntityDef
	repositoryFactory WriterRepositoryFactory
}

func (this *WriterRepositoryTest) WithEntity(ed orm.EntityDef) *WriterRepositoryTest {
	this.entityDef = ed
	return this
}

func (this *WriterRepositoryTest) WithInstanceFactory(rf func() orm.WriterRepository) *WriterRepositoryTest {
	this.repositoryFactory = WriterRepositoryFactoryFunc(rf)
	return this
}

// creates a new service test wrapper for given name
func NewWriterRepositoryTest() *WriterRepositoryTest {
	return &WriterRepositoryTest{}
}

type WriterRepositoryFactory interface {
	New() orm.WriterRepository
}

type WriterRepositoryFactoryFunc func() orm.WriterRepository

func (this WriterRepositoryFactoryFunc) New() orm.WriterRepository {
	return this()
}

func (this WriterRepositoryTest) Create_Error(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)
	entity, ok := tc.valueHolder.(orm.Entity)
	assert.True(t, ok, "value holder is not an entity")
	assert.Equal(t, uint64(0), entity.GetId(), "entity id is non-zero")

	mock.
		ExpectExec("insert into " + this.entityDef.TableName() + "(.*) values (.*)").
		WillReturnError(fmt.Errorf("could not insert"))

	// when
	err := repo.Save(context.Background(), entity)

	// then
	assert.NotNil(t, err, "expected error to be returned")
	assert.Equal(t, "could not insert", err.Error(), "error message is not correct")
}

func (this WriterRepositoryTest) Create_Success(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)
	entity, ok := tc.valueHolder.(orm.Entity)
	assert.True(t, ok, "value holder is not an entity")
	assert.Equal(t, uint64(0), entity.GetId(), "entity id is non-zero")

	mock.
		ExpectExec("insert into " + this.entityDef.TableName() + "(.*) values (.*)").
		WillReturnResult(sqlmock.NewResult(100, 1))

	// when
	err := repo.Save(context.Background(), entity)

	// then
	assert.Nil(t, err, "no error expected")
}

func (this WriterRepositoryTest) Update_Error(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)
	entity, ok := tc.valueHolder.(orm.Entity)
	assert.True(t, ok, "value holder is not an entity")
	assert.True(t, entity.GetId() > 0, "entity id is zero")

	var q string
	_, ok = tc.valueHolder.(orm.Versioned)
	if ok {
		q = "update " + this.entityDef.TableName() + " set .* where id=\\? AND version=\\?"
	} else {
		q = "update " + this.entityDef.TableName() + " set .* where id=\\?"
	}

	mock.
		ExpectExec(q).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// when
	err := repo.Save(context.Background(), entity)

	// then
	assert.Nil(t, err, "no error expected")
}

func (this WriterRepositoryTest) Update_Success(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)
	entity, ok := tc.valueHolder.(orm.Entity)
	assert.True(t, ok, "value holder is not an entity")
	assert.True(t, entity.GetId() > 0, "entity id is zero")

	var q string
	_, ok = tc.valueHolder.(orm.Versioned)
	if ok {
		q = "update " + this.entityDef.TableName() + " set .* where id=\\? AND version=\\?"
	} else {
		q = "update " + this.entityDef.TableName() + " set .* where id=\\?"
	}

	mock.
		ExpectExec(q).
		WillReturnError(fmt.Errorf("could not update"))

	// when
	err := repo.Save(context.Background(), entity)

	// then
	assert.NotNil(t, err, "expected error to be returned")
	assert.Equal(t, "could not update", err.Error(), "error message is not correct")
}

func (this WriterRepositoryTest) NewRepoWithMock(t *testing.T) (orm.WriterRepository, sqlmock.Sqlmock) {
	mock := NewSqlMock(t)
	repo := this.repositoryFactory.New()
	return repo, mock
}
