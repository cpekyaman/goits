package validation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type baseEntity struct {
	Id   uint64
	Name string
}

type validatedEntity struct {
	baseEntity
	Amount int32
}

func TestGetContext_NotOk(t *testing.T) {
	// when
	_, ok := GetContext("test.OtherTest")

	// then
	assert.False(t, ok)
}

func TestGetContext_Ok(t *testing.T) {
	// given
	sv := Struct("test.AnotherTest")

	// when
	vc, ok := GetContext("test.AnotherTest")

	// then
	assert.True(t, ok, "should get a validation context")
	assert.Equal(t, sv, vc.sv, "validation struct of context is different then registered")
	assert.True(t, vc.IsValid(), "isValid should be true initially")
}

func TestValidate_UnknownField_Ok(t *testing.T) {
	// given
	sv := structValidation("test.Test")
	vc := GetContextFor(sv)

	// when
	vc.Validate("unknown", "garbage")

	// then
	assert.True(t, vc.IsValid(), "should still be valid")
}

func TestValidate_InvalidField(t *testing.T) {
	// given
	sv := structValidation("test.Test")
	sv.Field("demo").With(NotBlank())

	vc := GetContextFor(sv)

	// when
	vc.Validate("demo", " ")

	// then
	assert.False(t, vc.IsValid(), "should be invalid")
	assert.Equal(t, 1, len(vc.Errors().Errors), "we should have one validation error")
}

func TestValidateStruct_UnknownStruct_NoError(t *testing.T) {
	// when
	err := ValidateStruct("test.UnkonwnTest", validatedEntity{})

	// then
	assert.Nil(t, err, "no error expected")
}

func TestValidateStruct_NoError(t *testing.T) {
	// given
	sv := structValidation("test.Test")
	sv.
		Field("Name").With(NotBlank(), Pattern(PatternAlpha)).
		Field("Amount").With(IntMin(20))

	vc := GetContextFor(sv)

	// when
	err := vc.ValidateStruct(validatedEntity{baseEntity{uint64(1), "Demo"}, 30})

	// then
	assert.Nil(t, err, "no error expected")
}

func TestValidateStruct_Error(t *testing.T) {
	// given
	sv := structValidation("test.Test")
	sv.
		Field("Name").With(NotBlank(), Pattern(PatternAlpha)).
		Field("Amount").With(IntMin(20))

	vc := GetContextFor(sv)

	// when
	err := vc.ValidateStruct(validatedEntity{baseEntity{uint64(1), "Demo Fail"}, 10})

	// then
	assert.NotNil(t, err, "validation error expected")

	oe, ok := err.(*ObjectError)
	assert.True(t, ok, "should be ObjectError")

	errStr := oe.Error()
	assert.True(t, strings.Contains(errStr, "Amount:min"), errStr+" should contain Amount field")
	assert.True(t, strings.Contains(errStr, "Name:pattern"), errStr+" should contain Name field")
}
