package project

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cpekyaman/goits/framework/orm"
	"github.com/cpekyaman/goits/framework/testlib"
	"github.com/stretchr/testify/assert"
)

const (
	findOneId uint64 = 5
)

var readerTest *testlib.ReaderRepositoryTest
var writerTest *testlib.WriterRepositoryTest

func init() {
	readerTest = testlib.NewReaderRepositoryTest().
		WithDbMetaData(testlib.DBMetaData{Columns: []string{"id", "name"}}).
		WithEntity(projectED).
		WithInstanceFactory(func() orm.ReaderRepository { return newProjectRepository() })

	writerTest = testlib.NewWriterRepositoryTest().
		WithEntity(projectED).
		WithInstanceFactory(func() orm.WriterRepository { return newProjectRepository() })
}

func TestFindById_DataFound(t *testing.T) {
	id := uint64(1)

	context := testlib.NewTestContext().
		WithValue(&Project{}).
		DbMockWith(func(r *sqlmock.Rows) {
			r.AddRow(id, "test project")
		}).
		AssertWith(func(t *testing.T, result testlib.TestResult) {
			prj, ok := result.RawResult.(*Project)
			assert.True(t, ok, "not a project")

			assert.Equal(t, id, prj.Id, "id is not correct")
			assert.Equal(t, "test project", prj.Name, "name is not correct")
		})

	readerTest.FindById_DataFound(t, context)
}

func TestFindById_NoDataFound(t *testing.T) {
	context := testlib.NewTestContext().
		WithValue(&Project{}).
		AssertWith(func(t *testing.T, result testlib.TestResult) {
			prj, ok := result.RawResult.(*Project)
			assert.True(t, ok, "not a project")

			assert.Equal(t, uint64(0), prj.Id, "should not set id")
			assert.Equal(t, "", prj.Name, "should not set name")
		})

	readerTest.FindById_NoDataFound(t, context)
}

func TestFindById_Error(t *testing.T) {
	id := uint64(1)

	context := testlib.NewTestContext().
		WithValue(&Project{}).
		DbMockWith(func(rows *sqlmock.Rows) {
			rows.AddRow(id, 8)
		})

	readerTest.FindById_Error(t, context)
}

func TestFindAll_DataFound(t *testing.T) {
	context := testlib.NewTestContext().
		WithValue(&[]Project{}).
		DbMockWith(func(rows *sqlmock.Rows) {
			rows.AddRow(1, "first project")
			rows.AddRow(2, "second project")
			rows.AddRow(3, "third project")
		}).
		AssertWith(func(t *testing.T, result testlib.TestResult) {
			pp, ok := result.RawResult.(*[]Project)
			assert.True(t, ok, "not a project slice")
			projects := *pp

			assert.Equal(t, 3, len(projects), "not all rows are returned")
			for _, v := range projects {
				assert.NotEmpty(t, v.Id, "has empty struct")
			}
		})

	readerTest.FindAll_DataFound(t, context)
}

func TestFindAll_NoDataFound(t *testing.T) {
	readerTest.FindAll_NoDataFound(t, &[]Project{})
}

func TestFindAll_Error(t *testing.T) {
	context := testlib.NewTestContext().
		WithValue(&[]Project{}).
		DbMockWith(func(rows *sqlmock.Rows) {
			rows.AddRow(1, "first project")
			rows.AddRow(2, "second project")
			rows.AddRow(3, "third project")
		})

	readerTest.FindAll_Error(t, context)
}

func TestCreate_Error(t *testing.T) {
	context := testlib.NewTestContext().
		WithValue(newDummyProject())

	writerTest.Create_Error(t, context)
}

func TestCreate_Success(t *testing.T) {
	context := testlib.NewTestContext().
		WithValue(newDummyProject())

	writerTest.Create_Success(t, context)
}

func TestUpdate_Error(t *testing.T) {
	prj := newDummyProject()
	prj.Id = uint64(100)

	context := testlib.NewTestContext().WithValue(prj)

	writerTest.Update_Error(t, context)
}

func TestUpdate_Success(t *testing.T) {
	prj := newDummyProject()
	prj.Id = uint64(100)

	context := testlib.NewTestContext().WithValue(prj)

	writerTest.Update_Success(t, context)
}

func newDummyProject() *Project {
	prj := NewProject()
	prj.Name = "Dummy"
	prj.Description = "New Dummy Project"
	return prj
}
