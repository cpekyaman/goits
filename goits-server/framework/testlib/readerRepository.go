package testlib

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cpekyaman/goits/framework/orm/metadata"
	"github.com/cpekyaman/goits/framework/orm/repository"
	"github.com/cpekyaman/goits/framework/testlib/assertions"
	"github.com/cpekyaman/goits/framework/testlib/mocking"
	"github.com/stretchr/testify/assert"
)

type ReaderRepositoryTest struct {
	entityDef         metadata.EntityDef
	metaData          DBMetaData
	repositoryFactory ReaderRepositoryFactory
	mocker            mocking.QueryMocker
}

func (this *ReaderRepositoryTest) WithDbMetaData(md DBMetaData) *ReaderRepositoryTest {
	this.metaData = md
	return this
}

func (this *ReaderRepositoryTest) WithInstanceFactory(rf func() repository.ReaderRepository) *ReaderRepositoryTest {
	this.repositoryFactory = ReaderRepositoryFactoryFunc(rf)
	return this
}

type ReaderRepositoryFactory interface {
	New() repository.ReaderRepository
}

type ReaderRepositoryFactoryFunc func() repository.ReaderRepository

func (this ReaderRepositoryFactoryFunc) New() repository.ReaderRepository {
	return this()
}

// creates a new service test wrapper for given name
func NewReaderRepositoryTest(ed metadata.EntityDef) *ReaderRepositoryTest {
	return &ReaderRepositoryTest{
		entityDef: ed,
		mocker:    mocking.NewQueryMocker(ed),
	}
}

//////////////////////
// Tests For FindOneById
//////////////////////

func (this ReaderRepositoryTest) FindById_DataFound(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)

	findOneId := uint64(4)
	_, rows := this.MockFindOneWithRows(findOneId, mock)
	tc.rowMocker.Mock(rows)

	// when
	err := repo.FindOneById(context.Background(), tc.valueHolder, findOneId)

	// then
	assert.Nil(t, err, err)
	tc.asserter.Assert(t, TestResult{tc.valueHolder, err})
}

func (this ReaderRepositoryTest) FindById_NoDataFound(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)

	findOneId := uint64(4)
	this.MockFindOneWithRows(findOneId, mock)

	// when
	err := repo.FindOneById(context.Background(), tc.valueHolder, findOneId)

	// then
	assert.NotNil(t, err, "should get back an error")
	assert.True(t, strings.Contains(err.Error(), "no rows"), fmt.Sprintf("unexpected error: %v", err))
	tc.asserter.Assert(t, TestResult{tc.valueHolder, err})
}

func (this ReaderRepositoryTest) FindById_Error(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)

	findOneId := uint64(4)
	this.mocker.ExpectFindOne(mock).WillReturnError(mocking.SqlError)

	// when
	err := repo.FindOneById(context.Background(), tc.valueHolder, findOneId)

	// then
	assert.NotNil(t, err, "should get back an error")
	assert.Equal(t, mocking.SqlError, err, "unexpected error")
}

//////////////////////
// Tests For FindOneByAttribute
//////////////////////

func (this ReaderRepositoryTest) FindByAttribute_DataFound(t *testing.T, tc *TestContext, attr string, attrVal interface{}) {
	// given
	repo, mock := this.NewRepoWithMock(t)

	_, rows := this.MockFindOneByAttributeWithRows(attr, attrVal, mock)
	tc.rowMocker.Mock(rows)

	// when
	err := repo.FindOneByAttribute(context.Background(), tc.valueHolder, attr, attrVal)

	// then
	assert.Nil(t, err, err)
	tc.asserter.Assert(t, TestResult{tc.valueHolder, err})
}

func (this ReaderRepositoryTest) FindByAttribute_NoDataFound(t *testing.T, tc *TestContext, attr string, attrVal interface{}) {
	// given
	repo, mock := this.NewRepoWithMock(t)

	this.MockFindOneByAttributeWithRows(attr, attrVal, mock)

	// when
	err := repo.FindOneByAttribute(context.Background(), tc.valueHolder, attr, attrVal)

	// then
	assert.NotNil(t, err, "should get back an error")
	assert.True(t, strings.Contains(err.Error(), "no rows"), fmt.Sprintf("unexpected error: %v", err))
	tc.asserter.Assert(t, TestResult{tc.valueHolder, err})
}

func (this ReaderRepositoryTest) FindByAttribute_Error(t *testing.T, tc *TestContext, attr string, attrVal interface{}) {
	// given
	repo, mock := this.NewRepoWithMock(t)

	this.mocker.ExpectFindOneByAttr(mock, attr).WillReturnError(mocking.SqlError)

	// when
	err := repo.FindOneByAttribute(context.Background(), tc.valueHolder, attr, attrVal)

	// then
	assert.NotNil(t, err, "should get back an error")
	assert.Equal(t, mocking.SqlError, err, "unexpected error")
}

//////////////////////
// Tests For FindAll
//////////////////////

func (this ReaderRepositoryTest) FindAll_DataFound(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)

	_, rows := this.MockFindAllWithRows(mock)
	tc.rowMocker.Mock(rows)

	// when
	err := repo.FindAll(context.Background(), tc.valueHolder)

	// then
	assert.Nil(t, err, "should not get error")
	tc.asserter.Assert(t, TestResult{tc.valueHolder, nil})
}

func (this ReaderRepositoryTest) FindAll_NoDataFound(t *testing.T, valueHolder interface{}) {
	// given
	repo, mock := this.NewRepoWithMock(t)

	this.MockFindAllWithRows(mock)

	// when
	err := repo.FindAll(context.Background(), valueHolder)

	// then
	assert.Nil(t, err, "should not get an error")
	assertions.ArrayEmpty(t, valueHolder)
}

func (this ReaderRepositoryTest) FindAll_Error(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)

	eq, rows := this.MockFindAllWithRows(mock)
	tc.rowMocker.Mock(rows)
	eq.WillReturnError(fmt.Errorf("for test"))

	// when
	err := repo.FindAll(context.Background(), tc.valueHolder)

	// then
	assert.NotNil(t, err, "should not back get error")
	assert.True(t, strings.Contains(err.Error(), "for test"), "unexpected error")
	assertions.ArrayEmpty(t, tc.valueHolder)
}

//////////////////////
// Tests For FindAllByAttributes
//////////////////////

func (this ReaderRepositoryTest) FindAllByAttributes_DataFound(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)

	_, rows := this.MockFindAllWithRows(mock)
	tc.rowMocker.Mock(rows)

	// when
	err := repo.FindAll(context.Background(), tc.valueHolder)

	// then
	assert.Nil(t, err, "should not get error")
	tc.asserter.Assert(t, TestResult{tc.valueHolder, nil})
}

//////////////////////
// Mock Helpers
//////////////////////

func (this ReaderRepositoryTest) MockFindOneWithRows(id uint64, mock sqlmock.Sqlmock) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	eq, rows := this.mocker.ExpectQueryWithRows(this.mocker.ExpectFindOne(mock), this.metaData.Columns)
	eq.WithArgs(id)
	return eq, rows
}

func (this ReaderRepositoryTest) MockFindAllWithRows(mock sqlmock.Sqlmock) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	eq, rows := this.mocker.ExpectQueryWithRows(this.mocker.ExpectFindAll(mock), this.metaData.Columns)
	return eq, rows
}

func (this ReaderRepositoryTest) MockFindOneByAttributeWithRows(attr string, attrVal interface{}, mock sqlmock.Sqlmock) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	eq, rows := this.mocker.ExpectQueryWithRows(this.mocker.ExpectFindOneByAttr(mock, attr), this.metaData.Columns)
	eq.WithArgs(attrVal)
	return eq, rows
}

func (this ReaderRepositoryTest) MockFindAllByAttributesWithRows(attrs map[string]interface{}, mock sqlmock.Sqlmock) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	eq, params := this.mocker.ExpectFindAllByAttributes(mock, attrs)
	eq, rows := this.mocker.ExpectQueryWithRows(eq, this.metaData.Columns)
	eq.WithArgs(params...)
	return eq, rows
}

func (this ReaderRepositoryTest) NewRepoWithMock(t *testing.T) (repository.ReaderRepository, sqlmock.Sqlmock) {
	mock := mocking.NewSqlMock(t)
	repo := this.repositoryFactory.New()
	return repo, mock
}
