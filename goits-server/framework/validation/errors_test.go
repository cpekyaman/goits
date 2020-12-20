package validation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldError_Error(t *testing.T) {
	// given
	fe := FieldError{
		Name:   "myField",
		Type:   "min",
		Params: map[string]interface{}{"min": 5},
	}

	// when
	err := fe.Error()

	// then
	assert.Equal(t, "validation: myField:min", err, "error message is not correct")
}

func TestObjectError_Error(t *testing.T) {
	// given
	first := FieldError{"first", "max", map[string]interface{}{"max": 10}}
	second := FieldError{"second", "pattern", map[string]interface{}{"pattern": PatternAlpha}}

	oe := ObjectError{
		Name:   "Test",
		Errors: []FieldError{first, second},
	}

	// when
	err := oe.Error()

	// then
	assert.True(t, strings.Contains(err, first.Error()), "object error does not contain first field")
	assert.True(t, strings.Contains(err, second.Error()), "object error does not contain second field")
}
