package validation

import (
	"fmt"
	"strings"
)

// FieldError represents a validation error for specific field and validator.
type FieldError struct {
	Name   string
	Type   string
	Params map[string]interface{}
}

// Error returns a simple string representation of the FieldError for logging.
func (this FieldError) Error() string {
	return fmt.Sprintf("validation: %s:%s", this.Name, this.Type)
}

// ObjectError is the container type that stores all field validation errors.
type ObjectError struct {
	Name   string
	Errors []FieldError
}

// Error returns a simple string representation of the ObjectError for logging.
func (this *ObjectError) Error() string {
	var sb strings.Builder
	for _, err := range this.Errors {
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}
