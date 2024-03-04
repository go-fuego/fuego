// This package is used to test that Fuego can generate an OpenAPI 3.1 spec from the routes registered in a server correctly.
// Dependency graph is as such:
// github.com/go-fuego/fuego/openapi3 -> github.com/go-fuego/fuego -> github.com/go-fuego/fuego/openapi3-testing
// Arrows here are used to denote the import relationship, it is only in one direction.
package openapi3testing

import (
	"encoding/json"
	"fmt"

	"github.com/go-fuego/fuego"
	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
)

func validateSpec(s *fuego.Server) error {
	// 1. Load an OpenAPI Spec into bytes
	specBytes, err := json.Marshal(s.OpenApiSpec)
	if err != nil {
		return fmt.Errorf("cannot marshal OpenAPI spec: %w", err)
	}

	// 2. Create a new OpenAPI document using libopenapi
	document, docErrs := libopenapi.NewDocument(specBytes)
	if docErrs != nil {
		return fmt.Errorf("cannot create OpenAPI document: %w", docErrs)
	}

	// 3. Create a new validator
	docValidator, validatorErrs := validator.NewValidator(document)
	if validatorErrs != nil {
		return fmt.Errorf("cannot create OpenAPI validator: %v", validatorErrs)
	}

	// 4. Validate!
	valid, validationErrs := docValidator.ValidateDocument()
	if !valid {
		for _, e := range validationErrs {
			// 5. Handle the error
			fmt.Printf("Type: %s, Failure: %s\n", e.ValidationType, e.Message)
			fmt.Printf("Fix: %s\n\n", e.HowToFix)
		}
	}

	return nil
}
