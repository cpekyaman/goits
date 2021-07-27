package testlib

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cpekyaman/goits/framework/orm/metadata"
	"github.com/cpekyaman/goits/framework/orm/repository"
	"github.com/cpekyaman/goits/framework/orm/domain"
	"github.com/cpekyaman/goits/framework/testlib/mocking"
	"github.com/stretchr/testify/assert"
)

type WriterRepositoryTest struct {
	entityDef         metadata.EntityDef
	qm                mocking.QueryMocker
	repositoryFactory WriterRepositoryFactory
}

func (this *WriterRepositoryTest) WithInstanceFactory(rf func() repository.WriterRepository) *WriterRepositoryTest {
	this.repositoryFactory = WriterRepositoryFactoryFunc(rf)
	return this
}

// creates a new service test wrapper for given name
func NewWriterRepositoryTest(ed metadata.EntityDef) *WriterRepositoryTest {
	return &WriterRepositoryTest{
		entityDef: ed,
		qm:        mocking.NewQueryMocker(ed),
	}
}

type WriterRepositoryFactory interface {
	New() repository.WriterRepository
}

type WriterRepositoryFactoryFunc func() repository.WriterRepository

func (this WriterRepositoryFactoryFunc) New() repository.WriterRepository {
	return this()
}

//////////////////////
// Tests For Create
//////////////////////

func (this WriterRepositoryTest) Create_Error(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)
	entity, ok := tc.valueHolder.(domain.Entity)
	assert.True(t, ok, "value holder is not an entity")
	assert.Equal(t, uint64(0), entity.GetId(), "entity id is non-zero")

	this.qm.ExpectInsertError(mock)

	// when
	err := repo.Save(context.Background(), entity)

	// then
	assert.Equal(t, mocking.SqlError, err, "error is not correct")
}

func (this WriterRepositoryTest) Create_Success(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)
	entity, ok := tc.valueHolder.(domain.Entity)
	assert.True(t, ok, "value holder is not an entity")
	assert.Equal(t, uint64(0), entity.GetId(), "entity id is non-zero")

	this.qm.ExpectInsert(mock, 100)

	// when
	err := repo.Save(context.Background(), entity)

	// then
	assert.Nil(t, err, "no error expected")
}

//////////////////////
// Tests For Update
//////////////////////

func (this WriterRepositoryTest) Update_Error(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)
	entity, ok := tc.valueHolder.(domain.Entity)
	assert.True(t, ok, "value holder is not an entity")
	assert.True(t, entity.GetId() > 0, "entity id is zero")

	this.qm.ExpectUpdateError(mock, tc.valueHolder)

	// when
	err := repo.Save(context.Background(), entity)

	// then
	assert.Equal(t, mocking.SqlError, err, "error is not correct")
}

func (this WriterRepositoryTest) Update_Success(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)
	entity, ok := tc.valueHolder.(domain.Entity)
	assert.True(t, ok, "value holder is not an entity")
	assert.True(t, entity.GetId() > 0, "entity id is zero")

	this.qm.ExpectUpdate(mock, tc.valueHolder)

	// when
	err := repo.Save(context.Background(), entity)

	// then
	assert.Nil(t, err, "no error expected")
}

//////////////////////
// Mock Helpers
//////////////////////

func (this WriterRepositoryTest) NewRepoWithMock(t *testing.T) (repository.WriterRepository, sqlmock.Sqlmock) {
	mock := mocking.NewSqlMock(t)
	repo := this.repositoryFactory.New()
	return repo, mock
}
