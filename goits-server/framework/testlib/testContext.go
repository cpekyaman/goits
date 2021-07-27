// utilities and helpers for writing tests.
package testlib

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// NewTestContext creates a new empty test context.
func NewTestContext() *TestContext {
	return &TestContext{}
}

// TestContext is used to provide test helpers with necessary hooks.
type TestContext struct {
	rowMocker   RowMocker
	execMocker  ExecMocker
	asserter    Asserter
	valueHolder interface{}
	valueBinder func(target interface{}) error
}

// WithRowMock provides the function that will set the rows to be returned as mock results.
// This mock is used to set read queries with the expected result rows to be returned.
func (tc *TestContext) WithRowMock(mockFunc func(rows *sqlmock.Rows)) *TestContext {
	tc.rowMocker = RowMockerFunc(mockFunc)
	return tc
}

// WithExecMock provides the function that will set the rows to be returned as mock results.
// This mock is used to set modifying queries with correct expected arguments .
func (tc *TestContext) WithExecMock(mockFunc func(exec *sqlmock.ExpectedExec)) *TestContext {
	tc.execMocker = ExecMockerFunc(mockFunc)
	return tc
}

// WithAsserter provides the function that will verify the result obtained by test runner helper.
// The asserter is essentiall used for read query involving tests to verify returned result is expected.
func (tc *TestContext) WithAsserter(assertFunc func(t *testing.T, result TestResult)) *TestContext {
	tc.asserter = AsserterFunc(assertFunc)
	return tc
}

// WithValue provides the pointer to actual target struct to scan query results into.
// This is used to provide pointer to a struct which the query results will be scanned into.
func (tc *TestContext) WithValue(target interface{}) *TestContext {
	tc.valueHolder = target
	return tc
}

// WithBinder provides an ObjectBinder function to use when it is relevant.
// The binder is used for create / update related tests to fill in the object with exptected data.
func (tc *TestContext) WithBinder(f func(target interface{}) error) *TestContext {
	tc.valueBinder = f
	return tc
}

// RowMocker is the interface to receive mock rows implementation to fill in the mock results.
type RowMocker interface {
	// Mock is to simply set the rows to be returned as mock results.
	Mock(rows *sqlmock.Rows)
}

// RowMockerFunc is a wrapper type to use compatible functions as RowMocker.
type RowMockerFunc func(rows *sqlmock.Rows)

// Mock wraps the provided function in order to use it as RowMocker.
func (f RowMockerFunc) Mock(rows *sqlmock.Rows) {
	f(rows)
}

// ExecMocker is used to customize expected db exec.
type ExecMocker interface {
	Mock(exec *sqlmock.ExpectedExec)
}

// ExecMockerFunc is a wrapper type to use compatible functions as ExecMocker.
type ExecMockerFunc func(exec *sqlmock.ExpectedExec)

// Mock wraps the provided function in order to use it as ExecMocker.
func (f ExecMockerFunc) Mock(exec *sqlmock.ExpectedExec) {
	f(exec)
}

// DBMetaData provides the data template to be mocked.
// It is used to tell sqlmock what the structure of expected data is (such as when mocking row results).
type DBMetaData struct {
	Columns []string
}

// TestResult keeps the result of when part in a test execution to be asserted by Asserter.
type TestResult struct {
	RawResult interface{}
	Error     error
}

// Asserter is a generic asserter interface for tests to consume results for their case specific assertions.
type Asserter interface {
	Assert(t *testing.T, result TestResult)
}

// AsserterFunc allows using any function with correct signature as Asserter
type AsserterFunc func(t *testing.T, result TestResult)

func (f AsserterFunc) Assert(t *testing.T, result TestResult) {
	f(t, result)
}

// NoOpAsserter is an asserter that does nothing.
// It is provided so that tests don't need to do a nil check for asserter.
func NoOpAsserter() Asserter {
	return AsserterFunc(func(t *testing.T, r TestResult) {})
}
