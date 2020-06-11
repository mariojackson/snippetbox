package forms

type errors map[string][]string

// Add adds the given error message to the given field.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message of the given field.
func (e errors) Get(field string) string {
	errorMessages := e[field]
	if len(errorMessages) == 0 {
		return ""
	}
	return errorMessages[0]
}
