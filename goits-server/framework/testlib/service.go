package testlib

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cpekyaman/goits/framework/orm"
	"github.com/stretchr/testify/assert"

	"github.com/cpekyaman/goits/framework/services"
)

// ServiceTest is the service test wrapper for running tests on services.
type ServiceTest struct {
	name     string
	metaData DBMetaData
	svc      ServiceFactory
}

// WithName sets the name of the target entity / service on this ServiceTest instance.
func (this *ServiceTest) WithName(n string) *ServiceTest {
	this.name = n
	return this
}

// WithDbMetaData sets the db related meta data of the target entity on this ServiceTest instance.
func (this *ServiceTest) WithDbMetaData(md DBMetaData) *ServiceTest {
	this.metaData = md
	return this
}

// WithFactory sets the function that will create new service instances on this ServiceTest instance.
func (this *ServiceTest) WithFactory(sf func() interface{}) *ServiceTest {
	this.svc = ServiceFactoryFunc(sf)
	return this
}

// ServiceFactory is the generic factory to instantiate a new service inside tests.
type ServiceFactory interface {
	New() interface{}
}

// ServiceFactoryFunc allows using any func with proper signature as ServiceFactory
type ServiceFactoryFunc func() interface{}

func (sif ServiceFactoryFunc) New() interface{} {
	return sif()
}

// NewServiceTest creates a new service test wrapper for given name.
func NewServiceTest() *ServiceTest {
	return &ServiceTest{}
}

func (this ServiceTest) NewMockDB(t *testing.T) sqlmock.Sqlmock {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err, "could not create mock db")
	orm.WithDB(mockDB, "sqlmock")

	t.Cleanup(func() {
		mockDB.Close()
	})

	return mock
}

// Tests and verifies GetAll method of service for success path.
func (this ServiceTest) GetAll_Success(t *testing.T, tc *TestContext) {
	// given
	mock := this.NewMockDB(t)
	si, ok := this.svc.New().(services.GetAllService)
	assert.True(t, ok, "service is not a GetAllService")

	rows := sqlmock.NewRows(this.metaData.Columns)
	tc.rowMocker.Mock(rows)

	ed := orm.EntityDefByName(this.name)
	assert.NotNil(t, ed, "no entity def to mock query")
	mock.ExpectQuery("from " + ed.TableName() + " order by " + ed.DefaultSort()).WillReturnRows(rows)

	// when
	rawResult, err := si.GetAll(context.Background())

	// then
	assert.Nil(t, err, "should not return error")
	assert.NotNil(t, rawResult, "should get non null result")

	result := TestResult{rawResult, err}
	tc.asserter.Assert(t, result)
}

// Tests and verifies GetAll method of service for error path.
func (this ServiceTest) GetAll_Error(t *testing.T) {
	// given
	mock := this.NewMockDB(t)
	si, ok := this.svc.New().(services.GetAllService)
	assert.True(t, ok, "service is not a GetAllService")

	merr := fmt.Errorf("sql error")
	mock.ExpectQuery("from .* order by ").WillReturnError(merr)

	// when
	rawResult, err := si.GetAll(context.Background())

	// then
	assert.Equal(t, merr, err, "should return error")
	assert.Nil(t, rawResult, "should not get a result")
}

func (this ServiceTest) Create_Binding_Error(t *testing.T) {
	// given
	mock := this.NewMockDB(t)
	mock.ExpectExec("insert into .*").WillReturnError(fmt.Errorf("sql: error"))

	si, ok := this.svc.New().(services.WriterService)
	assert.True(t, ok, "service is not a WriterService")

	expectedErr := fmt.Errorf("binding: error")

	// when
	err := si.Create(context.Background(), services.ObjectBinderFunc(func(target interface{}) error {
		return expectedErr
	}))

	// then
	assert.Equal(t, expectedErr, err, "should have returned binding error")
}
