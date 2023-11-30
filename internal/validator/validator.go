package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// Check if all fields are valid
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// Add error in the FieldErrors map if any
func (v *Validator) AddFieldError(key string, msg string) {

	//Initialize FieldErrors map should incase of nil value
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exist := v.FieldErrors[key]; !exist {
		v.FieldErrors[key] = msg
	}
}

// Check the field and add an error message to FieldErrors if validation check is not ok
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank() returns true if a value is not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if value contains not morethan n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= 100
}

// PermittedInt() returns true if a value is in a list of permitted integers.
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
