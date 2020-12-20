package testlib

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cpekyaman/goits/framework/orm"
	"github.com/stretchr/testify/assert"
)

type ReaderRepositoryTest struct {
	entityDef         orm.EntityDef
	metaData          DBMetaData
	repositoryFactory ReaderRepositoryFactory
}

func (this *ReaderRepositoryTest) WithEntity(ed orm.EntityDef) *ReaderRepositoryTest {
	this.entityDef = ed
	return this
}

func (this *ReaderRepositoryTest) WithDbMetaData(md DBMetaData) *ReaderRepositoryTest {
	this.metaData = md
	return this
}

func (this *ReaderRepositoryTest) WithInstanceFactory(rf func() orm.ReaderRepository) *ReaderRepositoryTest {
	this.repositoryFactory = ReaderRepositoryFactoryFunc(rf)
	return this
}

type ReaderRepositoryFactory interface {
	New() orm.ReaderRepository
}

type ReaderRepositoryFactoryFunc func() orm.ReaderRepository

func (this ReaderRepositoryFactoryFunc) New() orm.ReaderRepository {
	return this()
}

// creates a new service test wrapper for given name
func NewReaderRepositoryTest() *ReaderRepositoryTest {
	return &ReaderRepositoryTest{}
}

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
	repo, mock := this.NewRepoWithMock(t)

	findOneId := uint64(4)
	this.MockFindOneWithRows(findOneId, mock)

	// when
	err := repo.FindOneById(context.Background(), tc.valueHolder, findOneId)

	// then
	assert.NotNil(t, err, "should get back an error")
	fmt.Println(err.Error())
	assert.True(t, strings.Contains(err.Error(), "no rows"), fmt.Sprintf("unexpected error: %v", err))
	tc.asserter.Assert(t, TestResult{tc.valueHolder, err})
}

func (this ReaderRepositoryTest) FindById_Error(t *testing.T, tc *TestContext) {
	// given
	repo, mock := this.NewRepoWithMock(t)

	findOneId := uint64(4)
	eq, rows := this.MockFindOneWithRows(findOneId, mock)
	eq.WillReturnError(fmt.Errorf("for test"))
	tc.rowMocker.Mock(rows)

	// when
	err := repo.FindOneById(context.Background(), tc.valueHolder, findOneId)

	// then
	assert.NotNil(t, err, "should get back an error")
	assert.True(t, strings.Contains(err.Error(), "for test"), "unexpected error")
}

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
	this.verifyArrayPtrEmpty(t, valueHolder)
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
	this.verifyArrayPtrEmpty(t, tc.valueHolder)
}

func (this ReaderRepositoryTest) verifyArrayPtrEmpty(t *testing.T, valueHolder interface{}) {
	arrPtr := reflect.ValueOf(valueHolder)
	assert.Equal(t, reflect.Ptr, arrPtr.Kind(), "not a pointer")

	arrElem := arrPtr.Elem()
	switch reflect.TypeOf(arrElem.Interface()).Kind() {
	case reflect.Slice:
		assert.Equal(t, 0, arrElem.Len(), "no rows should be returned")
	default:
		assert.Fail(t, "not a slice")
	}
}

func (this ReaderRepositoryTest) MockFindOneWithRows(id uint64, mock sqlmock.Sqlmock) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	eq, rows := this.MockFinderWithRows("select \\* from "+this.entityDef.TableName()+" where "+this.entityDef.PKColumn()+" = \\$1", mock)
	eq.WithArgs(id)
	return eq, rows
}

func (this ReaderRepositoryTest) MockFindAllWithRows(mock sqlmock.Sqlmock) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	eq, rows := this.MockFinderWithRows("select \\* from "+this.entityDef.TableName()+" order by "+this.entityDef.DefaultSort(), mock)
	return eq, rows
}

func (this ReaderRepositoryTest) MockFinderWithRows(queryRegex string, mock sqlmock.Sqlmock) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	rows := sqlmock.NewRows(this.metaData.Columns)
	eq := mock.ExpectQuery(queryRegex).WillReturnRows(rows)
	return eq, rows
}

func (this ReaderRepositoryTest) NewRepoWithMock(t *testing.T) (orm.ReaderRepository, sqlmock.Sqlmock) {
	mock := NewSqlMock(t)
	repo := this.repositoryFactory.New()
	return repo, mock
}
