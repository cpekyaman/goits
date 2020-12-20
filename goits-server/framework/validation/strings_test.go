package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotEmpty_Valid(t *testing.T) {
	assert.False(t, NotEmpty().Validate(""), "empty string should be invalid")
	assert.True(t, NotEmpty().Validate("  "), "blank string should be valid")
	assert.True(t, NotEmpty().Validate("test"), "non-empty string should be valid")
}

func TestNotBlank_Valid(t *testing.T) {
	assert.False(t, NotBlank().Validate(""), "empty string should be invalid")
	assert.False(t, NotBlank().Validate("  "), "blank string should be valid")
	assert.True(t, NotBlank().Validate(" test "), "non-empty string with spaces should be valid")
	assert.True(t, NotBlank().Validate("test"), "non-empty string should be valid")
}

func TestStrLen_Valid(t *testing.T) {
	// with
	values := []struct {
		min   int
		max   int
		str   string
		valid bool
	}{
		{0, 5, "test", true},
		{0, 6, "notatest", false},
		{5, 9, "hellooo", true},
		{4, 8, "go", false},
		{3, 7, "hello", true},
	}

	for _, tv := range values {
		// given
		validator := StrLen(tv.min, tv.max)

		t.Run(fmt.Sprintf("%s in (%d,%d)", tv.str, tv.min, tv.max), func(t *testing.T) {
			// when
			result := validator.Validate(tv.str)

			// then
			if result != tv.valid {
				t.Errorf("expected to be %v but was %v", tv.valid, result)
			}
		})
	}
}

func TestPattern_Valid(t *testing.T) {
	// with
	values := []struct {
		pattern PatternType
		str     string
		valid   bool
	}{
		{PatternAlpha, "test", true},
		{PatternAlpha, "test1", false},
		{PatternAlpha, "test project", false},
		{PatternWordAlpha, "test project", true},
		{PatternAlNum, "test", true},
		{PatternAlNum, "test1", true},
		{PatternAlNum, "test project 2", false},
		{PatternAlNum, "test-project", false},
		{PatternWordAlnum, "test project 2", true},
	}

	for _, tv := range values {
		// given
		validator := Pattern(tv.pattern)

		t.Run(fmt.Sprintf("(%s) matches (%s)", tv.str, tv.pattern), func(t *testing.T) {
			// when
			result := validator.Validate(tv.str)

			// then
			if result != tv.valid {
				t.Errorf("expected to be %v but was %v", tv.valid, result)
			}
		})
	}
}
