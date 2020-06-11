package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Form struct, which holds form errors and and form values.
type Form struct {
	url.Values
	Errors errors
}

// New initializes a new form.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required marks all provided fields as required by checking
// if it's empty.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MaxLength checks given fields value is less than the max length.
// If the fields value is bigger, an error will be added to the form
// errors.
func (f *Form) MaxLength(field string, maxLength int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > maxLength {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters", maxLength))
	}
}

// PermittedValues checks if the value of the given field has one of the given options.
// If the value is not within the provided options, an error will be added to the form.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}

	for _, opt := range opts {
		if value == opt {
			return
		}
	}

	f.Errors.Add(field, "This field is invalid")
}

// Valid returns true if the form has no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
