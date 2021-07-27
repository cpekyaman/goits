package validation

import (
	"reflect"
)

// contextValidationProvider is the implementation ValidationProvider which uses ValidationContext object to perform validations.
type contextValidationProvider struct{}

// ValidateStruct tries to create a new ValidationContext for the given type and performs validation checks.
func (this contextValidationProvider) ValidateStruct(typeName string, entity interface{}) error {
	vc, ok := GetContext(typeName)
	if !ok {
		return nil
	}

	return vc.ValidateStruct(entity)
}

// GetContext creates a new ValidationContext for typeName if the type as validations.
func GetContext(typeName string) (*ValidationContext, bool) {
	sv, ok := validatorRegistry[typeName]
	if !ok {
		return nil, false
	}
	return GetContextFor(sv), true
}

// GetContextFor creates a new ValidationContext for the given existing StructValidation.
func GetContextFor(sv *StructValidation) *ValidationContext {
	return &ValidationContext{
		sv: sv,
		oe: &ObjectError{
			Name: sv.name,
		},
		valid: true,
	}
}

// ValidationContext represents execution of registered validators on a type.
type ValidationContext struct {
	sv    *StructValidation
	oe    *ObjectError
	valid bool
}

// ValidateStruct ıterates over the fields of given struct via reflection and validates them.
func (this *ValidationContext) ValidateStruct(entity interface{}) error {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return notAStructErr
	}

	this.innerValidate(v)
	return this.Errors()
}

// innerValidate is the internal method that recursively iterates over fields and invokes validations.
func (this *ValidationContext) innerValidate(v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)

		switch f.Type.Kind() {
		case reflect.Struct:
			if f.Anonymous {
				this.innerValidate(v.Field(i))
			} else {
				this.Validate(f.Name, v.Field(i))
			}
		case reflect.Ptr:
			this.innerValidate(v.Elem())
		default:
			this.Validate(f.Name, v.Field(i))
		}
	}
}

// Validate applies registered field validators to given value and returns this for method chaining.
func (this *ValidationContext) Validate(field string, value interface{}) *ValidationContext {
	fv, ok := this.sv.Get(field)
	if !ok {
		return this
	}

	for _, v := range fv.validators {
		if v.Validate(value) == false {
			this.valid = false
			this.oe.Errors = append(this.oe.Errors, FieldError{
				Name:   field,
				Type:   v.Type(),
				Params: v.Params(),
			})
		}
	}

	return this
}

// IsValid returns true if there are no validation errors.
func (this *ValidationContext) IsValid() bool {
	return this.valid
}

// Errors returns the ObjectError instance if there are any errors.
func (this *ValidationContext) Errors() *ObjectError {
	if this.valid {
		return nil
	} else {
		return this.oe
	}
}
