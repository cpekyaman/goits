package validation

import (
	"reflect"
)

var validatorRegistry map[string]*StructValidation

func init() {
	validatorRegistry = make(map[string]*StructValidation)
}

// Struct creates a new StructValidation for typeName and returns it for further customization.
func Struct(typeName string) *StructValidation {
	vl := structValidation(typeName)
	validatorRegistry[typeName] = vl
	return vl
}

// structValidation is used internally to create a StructValidation without registering it.
func structValidation(typeName string) *StructValidation {
	return &StructValidation{
		name:       typeName,
		validators: make(map[string]*FieldValidation),
	}
}

// Field creates a new FieldValidation for fieldName and returns it for further customization.
func (this *StructValidation) Field(fieldName string) *FieldValidation {
	fv := &FieldValidation{
		sv:         this,
		name:       fieldName,
		validators: make([]Validator, 0),
	}
	this.validators[fieldName] = fv
	return fv
}

// Get gets the FieldValidation instance registered for fieldName for validation purposes.
func (this *StructValidation) Get(fieldName string) (*FieldValidation, bool) {
	fv, ok := this.validators[fieldName]
	return fv, ok
}

// StructValidation is the container object that stores validations registered for a type.
type StructValidation struct {
	name       string
	validators map[string]*FieldValidation
}

// With adds the Validator v to list of validators to be applied to field.
func (this *FieldValidation) With(varr ...Validator) *StructValidation {
	this.validators = append(this.validators, varr...)
	return this.sv
}

// FieldValidation is the container object that stores validations to be applied for a struct field.
type FieldValidation struct {
	sv         *StructValidation
	name       string
	validators []Validator
}

// Validator is the contract that all validator implementations should satisfy.
type Validator interface {
	Type() string

	Params() map[string]interface{}

	Validate(value interface{}) bool
}

// validatorImpl is an instance of Validator.
type validatorImpl struct {
	name   string
	params map[string]interface{}
	vFunc  func(interface{}) bool
}

// Type returns the type of the validator for logging and rendering localized error messages.
func (this validatorImpl) Type() string {
	return this.name
}

// Params returns the parameters that the validator used.
func (this validatorImpl) Params() map[string]interface{} {
	return this.params
}

// Validate is the place where validation takes place.
func (this validatorImpl) Validate(value interface{}) bool {
	return this.vFunc(value)
}

// fieldValue gets the actual typed value if the the input is a reflect.Value.
func fieldValue(value interface{}) interface{} {
	in := value
	rv, ok := in.(reflect.Value)
	if ok {
		in = rv.Interface()
	}
	return in
}
