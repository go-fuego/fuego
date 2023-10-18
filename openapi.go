package op

import (
	"context"
	"encoding/json"
	"os"

	"log/slog"

	"github.com/getkin/kin-openapi/openapi3"
	httpSwagger "github.com/swaggo/http-swagger"
)

func ptr[T any](s T) *T {
	return &s
}

func NewOpenAPI() openapi3.T {
	info := &openapi3.Info{
		Title:       "OpenAPI",
		Description: "OpenAPI",
		Version:     "0.0.1",
	}
	paths := openapi3.Paths{
		"/": &openapi3.PathItem{
			Get: &openapi3.Operation{
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Value: &openapi3.Response{
							Description: ptr("OK"),
						},
					},
				},
			},
		},
	}
	spec := openapi3.T{
		OpenAPI: "3.0.0",
		Info:    info,
		Paths:   paths,
	}
	return spec
}

func (s *Server) GenerateOpenAPI() {

	// Validate
	err := s.spec.Validate(context.Background())
	if err != nil {
		slog.Error("Error validating spec", "error", err)
	}

	// Marshal spec to JSON
	dataJSON, err := json.Marshal(s.spec)
	if err != nil {
		slog.Error("Error marshalling spec to JSON", "error", err)
	}

	// Write spec to docs/openapi.json
	os.MkdirAll("docs", 0755)
	f, err := os.Create("docs/openapi.json")
	if err != nil {
		slog.Error("Error creating docs/openapi.json", "error", err)
	}
	defer f.Close()

	_, err = f.Write(dataJSON)
	if err != nil {
		slog.Error("Error marshalling spec to JSON", "error", err)
	}

	// Register docs/openapi.json
	Get(s, "/swagger/doc.json", func(ctx Ctx[any]) (any, error) {
		return s.spec, nil
	})

	GetStd(s, "/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
	))

}

func RegisterOpenAPIOperation(s *Server, method, path string) {
	operation := openapi3.NewOperation()
	requestBody := openapi3.NewRequestBody()
	schema := openapi3.NewObjectSchema().WithProperty("id", openapi3.NewStringSchema())
	requestBody.WithJSONSchema(schema)
	operation.RequestBody = &openapi3.RequestBodyRef{
		Value: requestBody,
	}

	operation.AddResponse(200, openapi3.NewResponse().WithDescription("OK"))
	s.spec.AddOperation(path, method, operation)
}
