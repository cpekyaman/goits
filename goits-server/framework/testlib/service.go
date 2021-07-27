package testlib

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cpekyaman/goits/framework/orm/db"
	"github.com/cpekyaman/goits/framework/orm/metadata"
	"github.com/cpekyaman/goits/framework/testlib/mocking"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/cpekyaman/goits/framework/caching"
	"github.com/cpekyaman/goits/framework/services"
	"github.com/cpekyaman/goits/framework/validation"
)

// ServiceTest is the service test wrapper for running tests on services.
type ServiceTest struct {
	name      string
	entityDef metadata.EntityDef
	metaData  DBMetaData
	qm        mocking.QueryMocker
	svc       ServiceFactory
}

// WithDbMetaData sets the db related meta data of the target entity on this ServiceTest instance.
func (this *ServiceTest) WithDbMetaData(md DBMetaData) *ServiceTest {
	this.metaData = md
	return this
}

// WithFactory sets the function that will create new service instances on this ServiceTest instance.
func (this *ServiceTest) WithFactory(sf func(caching.Cache, validation.ValidationProvider) interface{}) *ServiceTest {
	this.svc = ServiceFactoryFunc(sf)
	return this
}

// ServiceFactory is the generic factory to instantiate a new service inside tests.
type ServiceFactory interface {
	New(c caching.Cache, vp validation.ValidationProvider) interface{}
}

// ServiceFactoryFunc allows using any func with proper signature as ServiceFactory
type ServiceFactoryFunc func(c caching.Cache, vp validation.ValidationProvider) interface{}

func (sif ServiceFactoryFunc) New(c caching.Cache, vp validation.ValidationProvider) interface{} {
	return sif(c, vp)
}

// NewServiceTest creates a new service test wrapper for given entity definition.
func NewServiceTest(ed metadata.EntityDef) *ServiceTest {
	return &ServiceTest{
		name:      ed.Name(),
		entityDef: ed,
		qm:        mocking.NewQueryMocker(ed),
	}
}

// SvcMockContext contains a set of mocks to be used customizing behavior of the service under test when necessary.
type SvcMockContext struct {
	mock sqlmock.Sqlmock
	vp   *mocking.MockValidationProvider
	c    *mocking.MockCache
}

type ReaderMockContext struct {
	SvcMockContext
	svc services.ReaderService
}

type WriterMockContext struct {
	SvcMockContext
	svc services.WriterService
}

//////////////////////
// Tests For GetAll
//////////////////////

// Tests and verifies GetAll method of service for success path.
func (this ServiceTest) GetAll_Success(t *testing.T, tc *TestContext) {
	// given
	mc := this.NewReaderTestContext(t, gomock.NewController(t))

	_, rows := this.MockFindAllWithRows(mc.mock)
	tc.rowMocker.Mock(rows)

	// when
	rawResult, err := mc.svc.GetAll(context.Background())

	// then
	assert.Nil(t, err, "should not return error")
	assert.NotNil(t, rawResult, "should get non null result")

	result := TestResult{rawResult, err}
	tc.asserter.Assert(t, result)
}

// Tests and verifies GetAll method of service for error path.
func (this ServiceTest) GetAll_Error(t *testing.T) {
	// given
	mc := this.NewReaderTestContext(t, gomock.NewController(t))

	merr := fmt.Errorf("sql: error")
	this.qm.ExpectFindAll(mc.mock).WillReturnError(merr)

	// when
	rawResult, err := mc.svc.GetAll(context.Background())

	// then
	assert.Equal(t, merr, err, "should return error")
	assert.Nil(t, rawResult, "should not get a result")
}

//////////////////////
// Tests For GetById
//////////////////////

func (this ServiceTest) GetById_Cached_Hit(t *testing.T, tc *TestContext) {
	// given
	mc := this.NewReaderTestContext(t, gomock.NewController(t))

	findOneId := uint64(4)
	mc.c.EXPECT().
		GetOrCompute(gomock.Eq(caching.IdToKey(findOneId)), gomock.Any()).
		Times(1).
		Return(tc.valueHolder, nil)

	// when
	v, err := mc.svc.GetById(context.Background(), findOneId)

	// then
	assert.Nil(t, err, err)
	tc.asserter.Assert(t, TestResult{v, err})
}

func (this ServiceTest) GetById_Cached_Error(t *testing.T) {
	// given
	mc := this.NewReaderTestContext(t, gomock.NewController(t))

	findOneId := uint64(4)
	expectedErr := errors.New("sql: failure")
	mc.c.EXPECT().
		GetOrCompute(gomock.Eq(caching.IdToKey(findOneId)), gomock.Any()).
		Times(1).
		Return(nil, expectedErr)

	// when
	_, err := mc.svc.GetById(context.Background(), findOneId)

	// then
	assert.NotNil(t, err, "should get back an error")
	assert.Equal(t, expectedErr, err, "unexpected error")
}

func (this ServiceTest) GetById_NoCache_Success(t *testing.T, tc *TestContext) {
	// given
	mc := this.NewReaderTestContext(t, gomock.NewController(t))

	findOneId := uint64(4)
	_, rows := this.MockFindOneWithRows(findOneId, mc.mock)
	tc.rowMocker.Mock(rows)

	// when
	v, err := mc.svc.GetById(context.Background(), findOneId)

	// then
	assert.Nil(t, err, err)
	tc.asserter.Assert(t, TestResult{v, err})
}

func (this ServiceTest) GetById_NoCache_Error(t *testing.T) {
	// given
	mc := this.NewReaderTestContext(t, gomock.NewController(t))

	findOneId := uint64(4)
	this.qm.ExpectFindOne(mc.mock).WillReturnError(mocking.SqlError)

	// when
	_, err := mc.svc.GetById(context.Background(), findOneId)

	// then
	assert.NotNil(t, err, "should get back an error")
	assert.Equal(t, mocking.SqlError, err, "unexpected error")
}

//////////////////////
// Tests For Create
//////////////////////

func (this ServiceTest) Create_Binding_Error(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	mc := this.NewWriterTestContext(t, ctrl)

	expectedErr := errors.New("binding: error")

	// when
	err := mc.svc.Create(context.Background(), this.errObjectBinder(expectedErr))

	// then
	assert.Equal(t, expectedErr, err, "should have returned binding error")
}

func (this ServiceTest) Create_Validation_Error(t *testing.T, tc *TestContext) {
	// given
	ctrl := gomock.NewController(t)
	mc := this.NewWriterTestContext(t, ctrl)

	expectedErr := errors.New("validation: error")
	mc.vp.EXPECT().ValidateStruct(gomock.Eq(this.name), gomock.Any()).Return(expectedErr)

	// when
	err := mc.svc.Create(context.Background(), services.ObjectBinderFunc(tc.valueBinder))

	// then
	assert.Equal(t, expectedErr, err, "should have returned validation error")
}

func (this ServiceTest) Create_Db_Error(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	mc := this.NewWriterTestContext(t, ctrl)

	mc.vp.EXPECT().ValidateStruct(gomock.Eq(this.name), gomock.Any()).Return(nil)

	this.qm.ExpectInsertError(mc.mock)

	// when
	err := mc.svc.Create(context.Background(), this.noopObjectBinder())

	// then
	assert.Equal(t, mocking.SqlError, err, "should have returned db error")
}

func (this ServiceTest) Create_Success(t *testing.T, tc *TestContext) {
	// given
	ctrl := gomock.NewController(t)
	mc := this.NewWriterTestContext(t, ctrl)

	exec := this.qm.ExpectInsert(mc.mock, 1)
	tc.execMocker.Mock(exec)

	mc.vp.EXPECT().ValidateStruct(gomock.Eq(this.name), gomock.Any()).Return(nil)

	// when
	err := mc.svc.Create(context.Background(), services.ObjectBinderFunc(tc.valueBinder))

	// then
	assert.Nil(t, err, "create should be successfull")
}

//////////////////////
// Tests For Update
//////////////////////

func (this ServiceTest) Update_Find_Error(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	mc := this.NewWriterTestContext(t, ctrl)

	id := uint64(1)
	this.qm.ExpectFindOne(mc.mock).WillReturnError(mocking.SqlError)

	// when
	err := mc.svc.Update(context.Background(), id, this.noopObjectBinder())

	// then
	assert.Equal(t, mocking.SqlError, err, "should have returned sql error")
}

func (this ServiceTest) Update_NoDataFound_Error(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	mc := this.NewWriterTestContext(t, ctrl)

	id := uint64(1)
	this.MockFindOneWithRows(id, mc.mock)

	// when
	err := mc.svc.Update(context.Background(), id, this.noopObjectBinder())

	// then
	assert.NotNil(t, err, "should get back an error")
	assert.True(t, strings.Contains(err.Error(), "no rows"), fmt.Sprintf("unexpected error: %v", err))
}

func (this ServiceTest) Update_Binding_Error(t *testing.T, tc *TestContext) {
	// given
	ctrl := gomock.NewController(t)
	mc := this.NewWriterTestContext(t, ctrl)

	id := uint64(1)
	_, rows := this.MockFindOneWithRows(id, mc.mock)
	tc.rowMocker.Mock(rows)

	expectedErr := errors.New("binding: error")

	// when
	err := mc.svc.Update(context.Background(), id, this.errObjectBinder(expectedErr))

	// then
	assert.Equal(t, expectedErr, err, "should have returned binding error")
}

func (this ServiceTest) Update_Validation_Error(t *testing.T, tc *TestContext) {
	// given
	ctrl := gomock.NewController(t)
	mc := this.NewWriterTestContext(t, ctrl)

	id := uint64(1)
	_, rows := this.MockFindOneWithRows(id, mc.mock)
	tc.rowMocker.Mock(rows)

	expectedErr := errors.New("validation: error")
	mc.vp.EXPECT().ValidateStruct(gomock.Eq(this.name), gomock.Any()).Return(expectedErr)

	// when
	err := mc.svc.Update(context.Background(), id, services.ObjectBinderFunc(tc.valueBinder))

	// then
	assert.Equal(t, expectedErr, err, "should have returned validation error")
}

func (this ServiceTest) Update_Db_Error(t *testing.T, tc *TestContext) {
	// given
	ctrl := gomock.NewController(t)
	mc := this.NewWriterTestContext(t, ctrl)

	mc.vp.EXPECT().ValidateStruct(gomock.Eq(this.name), gomock.Any()).Return(nil)

	id := uint64(1)
	this.qm.ExpectUpdateError(mc.mock, tc.valueHolder)

	// when
	err := mc.svc.Update(context.Background(), id, this.noopObjectBinder())

	// then
	assert.Equal(t, mocking.SqlError, err, "should have returned db error")
}

func (this ServiceTest) Update_Success(t *testing.T, tc *TestContext) {
	// given
	ctrl := gomock.NewController(t)
	mc := this.NewWriterTestContext(t, ctrl)

	mc.vp.EXPECT().ValidateStruct(gomock.Eq(this.name), gomock.Any()).Return(nil)

	id := uint64(1)
	_, rows := this.MockFindOneWithRows(id, mc.mock)
	tc.rowMocker.Mock(rows)

	exec := this.qm.ExpectUpdate(mc.mock, tc.valueHolder)
	tc.execMocker.Mock(exec)

	mc.c.EXPECT().Invalidate(gomock.Eq(caching.IdToKey(id)))

	// when
	err := mc.svc.Update(context.Background(), id, services.ObjectBinderFunc(tc.valueBinder))

	// then
	assert.Nil(t, err, "update should be successfull")
}

//////////////////////
// Mock Helpers
//////////////////////

func (this ServiceTest) NewReaderTestContext(t *testing.T, ctrl *gomock.Controller) ReaderMockContext {
	mock := this.NewMockDB(t)
	c := mocking.NewMockCache(ctrl)
	vp := mocking.NewMockValidationProvider(ctrl)

	si, ok := this.svc.New(c, vp).(services.ReaderService)
	assert.True(t, ok, "service is not a ReaderService")

	return ReaderMockContext{SvcMockContext{mock, vp, c}, si}
}

func (this ServiceTest) NewWriterTestContext(t *testing.T, ctrl *gomock.Controller) WriterMockContext {
	mock := this.NewMockDB(t)

	c := mocking.NewMockCache(ctrl)
	vp := mocking.NewMockValidationProvider(ctrl)

	si, ok := this.svc.New(c, vp).(services.WriterService)
	assert.True(t, ok, "service is not a WriterService")

	return WriterMockContext{SvcMockContext{mock, vp, c}, si}
}

func (this ServiceTest) MockFindOneWithRows(id uint64, mock sqlmock.Sqlmock) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	eq, rows := this.qm.ExpectQueryWithRows(this.qm.ExpectFindOne(mock), this.metaData.Columns)
	eq.WithArgs(id)
	return eq, rows
}

func (this ServiceTest) MockFindAllWithRows(mock sqlmock.Sqlmock) (*sqlmock.ExpectedQuery, *sqlmock.Rows) {
	eq, rows := this.qm.ExpectQueryWithRows(this.qm.ExpectFindAll(mock), this.metaData.Columns)
	return eq, rows
}

func (this ServiceTest) NewMockDB(t *testing.T) sqlmock.Sqlmock {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err, "could not create mock db")
	db.WithDB(mockDB, "sqlmock")

	t.Cleanup(func() {
		mockDB.Close()
	})

	return mock
}

func (this ServiceTest) noopObjectBinder() services.ObjectBinderFunc {
	return services.ObjectBinderFunc(func(target interface{}) error {
		return nil
	})
}

func (this ServiceTest) errObjectBinder(err error) services.ObjectBinderFunc {
	return services.ObjectBinderFunc(func(target interface{}) error {
		return err
	})
}
