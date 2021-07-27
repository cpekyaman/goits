package assertions

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ArrayEmpty checks if the given value is an array pointer and if the underlying array is empty.
func ArrayEmpty(t *testing.T, valueHolder interface{}) {
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
