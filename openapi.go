package op

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"

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

func (s *Server) GenerateOpenAPI() openapi3.T {
	// Validate
	err := s.spec.Validate(context.Background())
	if err != nil {
		slog.Error("Error validating spec", "error", err)
	}

	// Marshal spec to JSON
	jsonSpec, err := json.Marshal(s.spec)
	if err != nil {
		slog.Error("Error marshalling spec to JSON", "error", err)
	}

	if !s.OpenapiConfig.DisableSwagger {
		generateSwagger(s, jsonSpec)
	}

	if !s.OpenapiConfig.DisableLocalSave {
		localSave(s.OpenapiConfig.JsonSpecLocalPath, jsonSpec)
	}

	return s.spec
}

func localSave(jsonSpecLocalPath string, jsonSpec []byte) {
	jsonFolder := filepath.Dir(jsonSpecLocalPath)

	err := os.MkdirAll(jsonFolder, 0o750)
	if err != nil {
		slog.Error("Error creating docs directory", "error", err)
	}

	f, err := os.Create(jsonSpecLocalPath)
	if err != nil {
		slog.Error("Error creating "+jsonSpecLocalPath, "error", err)
	}
	defer f.Close()
	_, err = f.Write(jsonSpec)
	if err != nil {
		slog.Error("Error writing file", "error", err)
	}

	slog.Info("Updated " + jsonSpecLocalPath)
}

// Registers the routes to serve the OpenAPI spec and Swagger UI.
func generateSwagger(s *Server, jsonSpec []byte) {
	GetStd(s, s.OpenapiConfig.JsonSpecUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(jsonSpec)
	})

	GetStd(s, s.OpenapiConfig.SwaggerUrl+"/", httpSwagger.Handler(
		httpSwagger.Layout(httpSwagger.BaseLayout),
		httpSwagger.PersistAuthorization(true),
		httpSwagger.URL(s.OpenapiConfig.JsonSpecUrl), // The url pointing to API definition
	))

	slog.Info(fmt.Sprintf("Raw json spec available at http://localhost%s%s", s.Addr, s.OpenapiConfig.JsonSpecUrl))
	slog.Info(fmt.Sprintf("OpenAPI generated at http://localhost%s%s/index.html", s.Addr, s.OpenapiConfig.SwaggerUrl))
}

func validateJsonSpecLocalPath(jsonSpecLocalPath string) bool {
	jsonSpecLocalPathRegexp := regexp.MustCompile(`^[^\/][\/a-zA-Z0-9\-\_]+(.json)$`)
	return jsonSpecLocalPathRegexp.MatchString(jsonSpecLocalPath)
}

func validateJsonSpecUrl(jsonSpecUrl string) bool {
	jsonSpecUrlRegexp := regexp.MustCompile(`^\/[\/a-zA-Z0-9\-\_]+(.json)$`)
	return jsonSpecUrlRegexp.MatchString(jsonSpecUrl)
}

func validateSwaggerUrl(swaggerUrl string) bool {
	swaggerUrlRegexp := regexp.MustCompile(`^\/[\/a-zA-Z0-9\-\_]+[a-zA-Z0-9\-\_]$`)
	return swaggerUrlRegexp.MatchString(swaggerUrl)
}

var generator = openapi3gen.NewGenerator(
	openapi3gen.UseAllExportedFields(),
)

func RegisterOpenAPIOperation[T any, B any](s *Server, method, path string) (*openapi3.Operation, error) {
	operation := openapi3.NewOperation()

	// Tags
	tag := tagFromType(*new(T))
	if tag != "unknown-interface" {
		operation.Tags = append(operation.Tags, tag)
	}

	// Request body
	bodyTag := tagFromType(*new(B))
	if (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) && bodyTag != "unknown-interface" && bodyTag != "string" {

		bodySchema, ok := s.spec.Components.Schemas[bodyTag]
		if !ok {
			var err error
			bodySchema, err = generator.NewSchemaRefForValue(new(B), s.spec.Components.Schemas)
			if err != nil {
				return operation, err
			}
			s.spec.Components.Schemas[bodyTag] = bodySchema
		}

		requestBody := openapi3.NewRequestBody().
			WithRequired(true).
			WithDescription("Request body for " + reflect.TypeOf(*new(B)).String())

		if bodySchema != nil {
			content := openapi3.NewContentWithSchema(bodySchema.Value, []string{"application/json"})
			content["application/json"].Schema.Ref = "#/components/schemas/" + bodyTag
			requestBody.WithContent(content)
		}

		s.spec.Components.RequestBodies[bodyTag] = &openapi3.RequestBodyRef{
			Value: requestBody,
		}

		// add request body to operation
		operation.RequestBody = &openapi3.RequestBodyRef{
			Ref:   "#/components/requestBodies/" + bodyTag,
			Value: requestBody,
		}
	}

	// Response body
	responseSchema, ok := s.spec.Components.Schemas[tag]
	if !ok {
		var err error
		responseSchema, err = generator.NewSchemaRefForValue(new(T), s.spec.Components.Schemas)
		if err != nil {
			return operation, err
		}
		s.spec.Components.Schemas[tag] = responseSchema
	}

	response := openapi3.NewResponse().WithDescription("OK")
	if responseSchema != nil {
		content := openapi3.NewContentWithSchema(responseSchema.Value, []string{"application/json"})
		content["application/json"].Schema.Ref = "#/components/schemas/" + tag
		response.WithContent(content)
	}
	operation.AddResponse(200, response)

	// Path parameters
	for _, pathParam := range parsePathParams(path) {
		parameter := openapi3.NewPathParameter(pathParam)
		parameter.Schema = openapi3.NewStringSchema().NewRef()
		operation.AddParameter(parameter)
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
