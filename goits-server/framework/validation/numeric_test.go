package validation

import (
	"fmt"
	"testing"
)

func TestIntMin_Valid(t *testing.T) {
	// with
	values := []struct {
		min   int32
		value int32
		valid bool
	}{
		{0, 5, true},
		{5, 2, false},
		{0, -5, false},
		{-5, -10, false},
		{-5, -3, true},
		{2, 2, true},
	}

	for _, tv := range values {
		// given
		validator := IntMin(tv.min)

		t.Run(fmt.Sprintf("%d gte %d", tv.value, tv.min), func(t *testing.T) {
			// when
			result := validator.Validate(tv.value)

			// then
			if result != tv.valid {
				t.Errorf("expected to be %v but was %v", tv.valid, result)
			}
		})
	}
}

func TestIntMax_Valid(t *testing.T) {
	// with
	values := []struct {
		max   int32
		value int32
		valid bool
	}{
		{0, 5, false},
		{5, 2, true},
		{0, -5, true},
		{-5, -10, true},
		{-5, -3, false},
		{4, 4, true},
	}

	for _, tv := range values {
		// given
		validator := IntMax(tv.max)

		t.Run(fmt.Sprintf("%d lte %d", tv.value, tv.max), func(t *testing.T) {
			// when
			result := validator.Validate(tv.value)

			// then
			if result != tv.valid {
				t.Errorf("expected to be %v but was %v", tv.valid, result)
			}
		})
	}
}

func TestIntRange_Valid(t *testing.T) {
	// with
	values := []struct {
		min   int32
		max   int32
		value int32
		valid bool
	}{
		{0, 10, 5, true},
		{0, 5, -1, false},
		{10, 15, 20, false},
		{-10, -5, -8, true},
		{-5, 0, -6, false},
	}

	for _, tv := range values {
		// given
		validator := IntRange(tv.min, tv.max)

		t.Run(fmt.Sprintf("%d in (%d,%d)", tv.value, tv.min, tv.max), func(t *testing.T) {
			// when
			result := validator.Validate(tv.value)

			// then
			if result != tv.valid {
				t.Errorf("expected to be %v but was %v", tv.valid, result)
			}
		})
	}
}
