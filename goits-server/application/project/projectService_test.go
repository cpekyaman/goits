package project

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cpekyaman/goits/framework/caching"
	"github.com/cpekyaman/goits/framework/testlib"
	"github.com/cpekyaman/goits/framework/validation"
)

var st *testlib.ServiceTest

func init() {
	st = testlib.NewServiceTest(projectED).
		WithDbMetaData(testlib.DBMetaData{Columns: []string{"id", "name", "description"}}).
		WithFactory(func(c caching.Cache, vp validation.ValidationProvider) interface{} {
			return newProjectService(newProjectRepository(), c, vp)
		})
}

func TestSVC_Project_GetAll_Success(t *testing.T) {
	context := testlib.NewTestContext().
		WithRowMock(func(r *sqlmock.Rows) {
			r.AddRow(1, "Demo 1", "Demo Project One")
			r.AddRow(2, "Demo 2", "Demo Project Two")
		}).
		WithAsserter(func(t *testing.T, result testlib.TestResult) {
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

func TestSVC_Project_GetAll_Error(t *testing.T) {
	st.GetAll_Error(t)
}

func TestSVC_Project_GetById_Success(t *testing.T) {
	prj := &Project{
		Name:        "Single",
		Description: "Single Project",
	}
	context := testlib.NewTestContext().
		WithValue(prj).
		WithAsserter(func(t *testing.T, result testlib.TestResult) {
			p, ok := result.RawResult.(*Project)

			assert.True(t, ok, "not a project entity")

			assert.Equal(t, "Single", p.Name, "object does not have correct name")
			assert.Equal(t, "Single Project", p.Description, "object does not have correct description")
		})

	st.GetById_Cached_Hit(t, context)
}

func TestSVC_Project_GetById_Error(t *testing.T) {
	st.GetById_Cached_Error(t)
}

func TestSVC_Project_Create_Binding_Error(t *testing.T) {
	st.Create_Binding_Error(t)
}
func TestSVC_Project_Create_Db_Error(t *testing.T) {
	st.Create_Db_Error(t)
}

func TestSVC_Project_Create_Validation_Error(t *testing.T) {
	tc := testlib.NewTestContext().
		WithBinder(defaultBinder)

	st.Create_Validation_Error(t, tc)
}

func TestSVC_Project_Create_Success(t *testing.T) {
	tc := testlib.NewTestContext().
		WithBinder(defaultBinder).
		WithExecMock(func(exec *sqlmock.ExpectedExec) {
			args := []driver.Value{"demo project", "demo", 1, 2}
			exec.WithArgs(args...)
		})

	st.Create_Success(t, tc)
}

func TestSVC_Project_Update_Find_Error(t *testing.T) {
	st.Update_Find_Error(t)
}

func TestSVC_Project_Update_NoDataFound_Error(t *testing.T) {
	st.Update_NoDataFound_Error(t)
}
func TestSVC_Project_Update_Binding_Error(t *testing.T) {
	id := uint64(1)

	tc := testlib.NewTestContext().
		WithValue(&Project{}).
		WithRowMock(func(r *sqlmock.Rows) {
			r.AddRow(id, "test", "test project")
		})

	st.Update_Binding_Error(t, tc)
}

func TestSVC_Project_Update_Validation_Error(t *testing.T) {
	id := uint64(1)

	tc := testlib.NewTestContext().
		WithValue(&Project{}).
		WithBinder(defaultBinder).
		WithRowMock(func(r *sqlmock.Rows) {
			r.AddRow(id, "test", "test project")
		})

	st.Update_Validation_Error(t, tc)
}

func TestSVC_Project_Update_Db_Error(t *testing.T) {
	id := uint64(1)

	tc := testlib.NewTestContext().
		WithValue(&Project{}).
		WithRowMock(func(r *sqlmock.Rows) {
			r.AddRow(id, "test", "test project")
		})

	st.Update_Db_Error(t, tc)
}

func TestSVC_Project_Update_Success(t *testing.T) {
	id := uint64(1)

	tc := testlib.NewTestContext().
		WithValue(&Project{}).
		WithBinder(defaultBinder).
		WithRowMock(func(r *sqlmock.Rows) {
			r.AddRow(id, "test", "test project")
		}).
		WithExecMock(func(exec *sqlmock.ExpectedExec) {
			args := []driver.Value{"demo project", "demo", 1, 2, 1, 2}
			exec.WithArgs(args...)
		})

	st.Update_Success(t, tc)
}

func defaultBinder(target interface{}) error {
	prj, ok := target.(*Project)
	if !ok {
		return nil
	}
	prj.Name = "demo"
	prj.Description = "demo project"
	prj.Status = 1
	prj.Type = 2
	prj.Version = 2

	return nil
}
