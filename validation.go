package fuego

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type fieldValidationError struct {
	DevField string      `json:"devField,omitempty"` // Dev-friendly name of the field e.g. "user.Name"
	Field    string      `json:"field,omitempty"`    // User-friendly name of the field e.g. "Name"
	Tag      string      `json:"tag,omitempty"`      // Name of the validation tag, e.g. "required"
	Param    string      `json:"param,omitempty"`    // Parameter of the validation tag, e.g. "3" in "min=3"
	Value    interface{} `json:"value,omitempty"`    // Actual value of the field, e.g. "" (empty string so that's why the validation failed)
}

type structValidationError struct {
	Errors []fieldValidationError
}

func (e structValidationError) Status() int {
	return http.StatusBadRequest
}

func (e structValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "validation error"
	}

	messages := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		humanReadableError := ""
		switch err.Tag {
		case "required":
			humanReadableError = fmt.Sprintf("%s is required", err.Field)
		case "email":
			humanReadableError = fmt.Sprintf("%s should be a valid email", err.Field)
		case "uuid":
			humanReadableError = fmt.Sprintf("%s should be a valid UUID", err.Field)
		case "e164":
			humanReadableError = fmt.Sprintf("%s should be a valid international phone number (e.g. +33 6 06 06 06 06)", err.Field)
		default:
			humanReadableError = fmt.Sprintf("%s should be %s", err.Field, err.Tag)
			if err.Param != "" {
				humanReadableError += "=" + err.Param
			}
		}
		messages = append(messages, humanReadableError)
	}

	return strings.Join(messages, ", ")
}

func (e structValidationError) Info() map[string]any {
	return map[string]any{
		"validation": e.Errors,
	}
}

var v = validator.New()

func validate(a any) error {
	_, ok := a.(map[string]any)
	if ok {
		return nil
	}

	err := v.Struct(a)
	if err != nil {
		// this check is only needed when your code could produce an
		// invalid value for validation such as interface with nil value
		if _, exists := err.(*validator.InvalidValidationError); exists {
			return fmt.Errorf("validation error: %w", err)
		}

		validationError := structValidationError{}
		for _, err := range err.(validator.ValidationErrors) {
			validationError.Errors = append(validationError.Errors, fieldValidationError{
				DevField: err.StructNamespace(),
				Field:    err.StructField(),
				Tag:      err.Tag(),
				Param:    err.Param(),
				Value:    err.Value(),
			})
		}

		return validationError
	}
	return nil
}
