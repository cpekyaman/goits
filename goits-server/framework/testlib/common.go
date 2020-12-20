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
	asserter    Asserter
	valueHolder interface{}
}

// DbMockWith provides the function that will set the rows to be returned as mock results.
func (tc *TestContext) DbMockWith(mockFunc func(rows *sqlmock.Rows)) *TestContext {
	tc.rowMocker = RowMockerFunc(mockFunc)
	return tc
}

// AssertWith provides the function that will verify the result obtained by test runner helper.
func (tc *TestContext) AssertWith(assertFunc func(t *testing.T, result TestResult)) *TestContext {
	tc.asserter = AsserterFunc(assertFunc)
	return tc
}

// WithValue provides the pointer to actual target struct to scan query results into.
func (tc *TestContext) WithValue(target interface{}) *TestContext {
	tc.valueHolder = target
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

// Provides the list of columns of the target entity to be mocked.
type DBMetaData struct {
	Columns []string
}

// keeps the result of when part in a test execution
type TestResult struct {
	RawResult interface{}
	Error     error
}

// a generic asserter interface for tests to consume results for their own assertions
type Asserter interface {
	Assert(t *testing.T, result TestResult)
}

// allows using any function with correct signature as Asserter
type AsserterFunc func(t *testing.T, result TestResult)

func (f AsserterFunc) Assert(t *testing.T, result TestResult) {
	f(t, result)
}

// an asserter that does nothing
func NoOpAsserter() Asserter {
	return AsserterFunc(func(t *testing.T, r TestResult) {})
}
