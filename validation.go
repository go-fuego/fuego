package fuego

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// explainError translates a validator error into a human readable string.
func explainError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s should be a valid email", err.Field())
	case "uuid":
		return fmt.Sprintf("%s should be a valid UUID", err.Field())
	case "e164":
		return fmt.Sprintf("%s should be a valid international phone number (e.g. +33 6 06 06 06 06)", err.Field())
	default:
		resp := fmt.Sprintf("%s should be %s", err.Field(), err.Tag())
		if err.Param() != "" {
			resp += "=" + err.Param()
		}
		return resp
	}
}

var v = validator.New()

func validate(a any) error {
	_, ok := a.(map[string]any)
	if ok {
		return nil
	}

	err := v.Struct(a)
	if err == nil {
		return nil
	}

	// this check is only needed when your code could produce an
	// invalid value for validation such as interface with nil value
	if _, exists := err.(*validator.InvalidValidationError); exists {
		return fmt.Errorf("validation error: %w", err)
	}

	validationError := HTTPError{
		Err:    err,
		Status: http.StatusBadRequest,
		Title:  "Validation Error",
	}
	var errorsSummary []string
	for _, err := range err.(validator.ValidationErrors) {
		errorsSummary = append(errorsSummary, explainError(err))
		validationError.Errors = append(validationError.Errors, ErrorItem{
			Name:   err.StructNamespace(),
			Reason: err.Error(),
			More: map[string]any{
				"nsField": err.StructNamespace(),
				"field":   err.StructField(),
				"tag":     err.Tag(),
				"param":   err.Param(),
				"value":   err.Value(),
			},
		})
	}

	validationError.Detail = strings.Join(errorsSummary, ", ")

	return validationError
}
