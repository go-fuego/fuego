package op

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type fieldValidationError struct {
	field string      // Name of the field e.g. "Name"
	tag   string      // Name of the validation tag, e.g. "required"
	param string      // Parameter of the validation tag, e.g. "3" in "min=3"
	value interface{} // Actual value of the field, e.g. "" (empty string so that's why the validation failed)
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
		switch err.tag {
		case "required":
			humanReadableError = fmt.Sprintf("%s is required", err.field)
		case "email":
			humanReadableError = fmt.Sprintf("%s should be a valid email", err.field)
		case "uuid":
			humanReadableError = fmt.Sprintf("%s should be a valid UUID", err.field)
		case "e164":
			humanReadableError = fmt.Sprintf("%s should be a valid international phone number (e.g. +33 6 06 06 06 06)", err.field)
		default:
			humanReadableError = fmt.Sprintf("%s should be %s", err.field, err.tag)
			if err.param != "" {
				humanReadableError += "=" + err.param
			}
		}
		messages = append(messages, humanReadableError)
	}

	return strings.Join(messages, ", ")
}

var v = validator.New()

func validate(a any) error {
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
				field: err.StructField(),
				tag:   err.Tag(),
				param: err.Param(),
				value: err.Value(),
			})
		}

		return validationError
	}
	return nil
}
