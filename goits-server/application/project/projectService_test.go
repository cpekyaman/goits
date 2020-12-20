package project

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cpekyaman/goits/framework/testlib"
)

var st *testlib.ServiceTest

func init() {
	st = testlib.NewServiceTest().
		WithName("Project").
		WithDbMetaData(testlib.DBMetaData{Columns: []string{"id", "name", "description"}}).
		WithFactory(func() interface{} { return newProjectService() })
}

func TestGetAllSuccess(t *testing.T) {
	context := testlib.NewTestContext().
		DbMockWith(func(r *sqlmock.Rows) {
			r.AddRow(1, "Demo 1", "Demo Project One")
			r.AddRow(2, "Demo 2", "Demo Project Two")
		}).
		AssertWith(func(t *testing.T, result testlib.TestResult) {
			projects, ok := result.RawResult.([]Project)

			assert.True(t, ok, "not a project slice")
			assert.Equal(t, 2, len(projects), "number of elements is not correct")

			assert.Equal(t, uint64(1), projects[0].Id, "first object does not have correct id")
			assert.Equal(t, "Demo 1", projects[0].Name, "first object does not have correct name")

			assert.Equal(t, uint64(2), projects[1].Id, "second object does not have correct id")
			assert.Equal(t, "Demo 2", projects[1].Name, "second object does not have correct name")
		})

	st.GetAll_Success(t, context)
}

func TestGetAllError(t *testing.T) {
	st.GetAll_Error(t)
}

func TestCreate_BindingError(t *testing.T) {
	st.Create_Binding_Error(t)
}
