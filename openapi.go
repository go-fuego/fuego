package op

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
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
	err = os.MkdirAll("docs", 0o755)
	if err != nil {
		slog.Error("Error creating docs directory", "error", err)
	}
	f, err := os.Create("docs/openapi.json")
	if err != nil {
		slog.Error("Error creating docs/openapi.json", "error", err)
	}
	defer f.Close()
	_, err = f.Write(dataJSON)
	if err != nil {
		slog.Error("Error marshalling spec to JSON", "error", err)
	}

	// Serve spec as JSON
	GetStd(s, "/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(dataJSON)
	})

	// Swagger UI
	GetStd(s, "/swagger/", httpSwagger.Handler(
		httpSwagger.Layout(httpSwagger.BaseLayout),
		httpSwagger.PersistAuthorization(true),
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
	))

	slog.Info(fmt.Sprintf("OpenAPI generated at http://localhost%s/swagger/index.html", s.Addr))
}

func RegisterOpenAPIOperation[T any, B any](s *Server, method, path string) {
	generator := openapi3gen.NewGenerator(
		openapi3gen.UseAllExportedFields(),
	)

	operation := openapi3.NewOperation()

	// Request body
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		requestBody := openapi3.NewRequestBody()
		bodySchema, err := generator.NewSchemaRefForValue(new(B), nil)
		if err != nil {
			panic(err)
		}
		requestBody.WithJSONSchema(bodySchema.Value)
		operation.RequestBody = &openapi3.RequestBodyRef{
			Value: requestBody,
		}
	}

	// Response body
	responseSchema, err := generator.NewSchemaRefForValue(new(T), nil)
	if err != nil {
		panic(err)
	}

	// Path parameters
	pathParams := parsePathParams(path)
	for _, pathParam := range pathParams {
		operation.AddParameter(&openapi3.Parameter{
			In:          "path",
			Name:        pathParam,
			Description: "",
			Required:    true,
			Schema:      openapi3.NewStringSchema().NewRef(),
		})
	}

	// Tags
	tag := tagFromType(*new(T))
	operation.Tags = []string{tag}

	operation.AddResponse(200, openapi3.NewResponse().
		WithDescription("OK").
		WithJSONSchema(responseSchema.Value),
	)

	s.spec.AddOperation(path, method, operation)
}

func tagFromType(v any) string {
	if v == nil {
		return "default"
	}

	return dive(reflect.TypeOf(v), 4)
}

// dive returns the name of the type of the given reflect.Type.
// If the type is a pointer, slice, array, map, channel, function, or unsafe pointer,
// it will dive into the type and return the name of the type it points to.
func dive(t reflect.Type, maxDepth int) string {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		if maxDepth == 0 {
			return "default"
		}
		return dive(t.Elem(), maxDepth-1)
	default:
		return t.Name()
	}
}
