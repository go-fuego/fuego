package op

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewOpenAPI() openapi3.T {
	info := &openapi3.Info{
		Title:       "OpenAPI",
		Description: "OpenAPI",
		Version:     "0.0.1",
	}
	spec := openapi3.T{
		OpenAPI: "3.0.3",
		Info:    info,
		Paths:   openapi3.Paths{},
		Components: &openapi3.Components{
			Schemas:       make(map[string]*openapi3.SchemaRef),
			RequestBodies: make(map[string]*openapi3.RequestBodyRef),
			Responses:     make(map[string]*openapi3.ResponseRef),
		},
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
	err = os.MkdirAll("docs", 0o750)
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
		slog.Error("Error writing file", "error", err)
	}
	slog.Info("Updated docs/openapi.json")

	time.Sleep(2 * time.Second)

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

func RegisterOpenAPIOperation[T any, B any](s *Server, method, path string) (*openapi3.Operation, error) {
	generator := openapi3gen.NewGenerator(
		openapi3gen.UseAllExportedFields(),
	)

	operation := openapi3.NewOperation()

	// Tags
	tag := tagFromType(*new(T))
	if tag != "unknown-interface" {
		operation.Tags = append(operation.Tags, tag)
	}

	// Request body
	bodyTag := tagFromType(*new(B))
	if (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) && bodyTag != "unknown-interface" && bodyTag != "string" {
		bodySchema, err := generator.NewSchemaRefForValue(new(B), s.spec.Components.Schemas)
		if err != nil {
			return operation, err
		}
		s.spec.Components.Schemas[bodyTag] = bodySchema

		content := openapi3.NewContentWithSchema(bodySchema.Value, []string{"application/json"})
		content["application/json"].Schema.Ref = fmt.Sprintf("#/components/schemas/%s", bodyTag)

		requestBody := openapi3.NewRequestBody().
			WithRequired(true).
			WithDescription(fmt.Sprintf("Request body for %s", reflect.TypeOf(*new(B)).String())).
			WithContent(content)

		s.spec.Components.RequestBodies[bodyTag] = &openapi3.RequestBodyRef{
			Value: requestBody,
		}

		// add request body to operation
		operation.RequestBody = &openapi3.RequestBodyRef{
			Ref:   fmt.Sprintf("#/components/requestBodies/%s", bodyTag),
			Value: requestBody,
		}
	}

	// Response body
	responseSchema, err := generator.NewSchemaRefForValue(new(T), s.spec.Components.Schemas)
	if err != nil {
		return operation, err
	}
	s.spec.Components.Schemas[tag] = responseSchema

	content := openapi3.NewContentWithSchema(responseSchema.Value, []string{"application/json"})
	content["application/json"].Schema.Ref = fmt.Sprintf("#/components/schemas/%s", bodyTag)
	operation.AddResponse(200, openapi3.NewResponse().
		WithDescription("OK").
		WithContent(content),
	)

	// Path parameters
	for _, pathParam := range parsePathParams(path) {
		operation.AddParameter(&openapi3.Parameter{
			In:          "path",
			Name:        pathParam,
			Description: "",
			Required:    true,
			Schema:      openapi3.NewStringSchema().NewRef(),
		})
	}

	s.spec.AddOperation(path, method, operation)

	return operation, nil
}

func tagFromType(v any) string {
	if v == nil {
		return "unknown-interface"
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
