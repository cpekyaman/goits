//go:generate mockgen -source=validator.go -destination=validator_mock.go -package=mocking
package validation

import (
	"errors"
	"reflect"
)

var notAStructErr = errors.New("validation: not a struct")
var notInitErr = errors.New("validation: validator not initialized")

var vp ValidationProvider

func InitValidation() {
	vp = contextValidationProvider{}
}

// Provider returns the default initialized ValidationProvider for the application.
func Provider() ValidationProvider {
	return vp
}

// ValidationProvider provides the public interface for invoking registered validations.
type ValidationProvider interface {
	ValidateStruct(typeName string, entity interface{}) error
}

// ValidateStruct performs the validation of given struct instance registered under the given name.
func ValidateStruct(typeName string, entity interface{}) error {
	if vp != nil {
		return vp.ValidateStruct(typeName, entity)
	}
	return notInitErr
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
