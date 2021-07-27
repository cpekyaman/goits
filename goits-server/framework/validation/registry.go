package validation

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

// WithMandatoryName register default validators for the name field of the struct.
func (this *StructValidation) WithMandatoryName() *StructValidation {
	this.Field("name").With(notBlank, Pattern(PatternAlNum))
	return this
}

// WithMandatoryDesc register default validators for the description field of the struct.
func (this *StructValidation) WithMandatoryDesc() *StructValidation {
	this.Field("description").With(notBlank, Pattern(PatternAlNum))
	return this
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
