package fuego

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-fuego/fuego/openapi3"
)

func NewOpenApiSpec() openapi3.Document {
	spec := openapi3.NewDocument()

	return spec
}

func (s *Server) generateOpenAPI() openapi3.Document {
	jsonSpec, err := json.Marshal(s.OpenApiSpec)
	if err != nil {
		slog.Error("Error marshalling OpenAPI spec", "error", err)
	}

	if !s.OpenapiConfig.DisableSwagger {
		generateSwagger(s, jsonSpec)
	}

	if !s.OpenapiConfig.DisableLocalSave {
		err := localSave(s.OpenapiConfig.JsonFilePath, jsonSpec)
		if err != nil {
			slog.Error("Error saving spec to local path", "error", err, "path", s.OpenapiConfig.JsonFilePath)
		}
	}

	return s.OpenApiSpec
}

func localSave(jsonSpecLocalPath string, jsonSpec []byte) error {
	jsonFolder := filepath.Dir(jsonSpecLocalPath)

	err := os.MkdirAll(jsonFolder, 0o750)
	if err != nil {
		return errors.New("error creating docs directory")
	}

	f, err := os.Create(jsonSpecLocalPath) // #nosec G304 (file path provided by developer, not by user)
	if err != nil {
		return errors.New("error creating file")
	}
	defer f.Close()

	_, err = f.Write(jsonSpec)
	if err != nil {
		return errors.New("error writing file ")
	}

	slog.Info("JSON file: " + jsonSpecLocalPath)
	return nil
}

// Registers the routes to serve the OpenAPI spec and Swagger UI.
func generateSwagger(s *Server, jsonSpec []byte) {
	GetStd(s, s.OpenapiConfig.JsonUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(jsonSpec)
	})

	Handle(s, s.OpenapiConfig.SwaggerUrl+"/", s.UIHandler(s.OpenapiConfig.JsonUrl))

	slog.Info(fmt.Sprintf("JSON spec: http://localhost%s%s", s.Server.Addr, s.OpenapiConfig.JsonUrl))
	slog.Info(fmt.Sprintf("OpenAPI UI: http://localhost%s%s/index.html", s.Server.Addr, s.OpenapiConfig.SwaggerUrl))
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

func RegisterOpenAPIOperation[T any, B any](s *Server, method, path string) (*openapi3.Operation, error) {
	operation := &openapi3.Operation{
		Summary:     "Summary",
		Description: "Description",
	}

	// Tags
	tag := tagFromType(*new(T))
	if tag != "unknown-interface" {
		operation.Tags = append(operation.Tags, tag)
	}

	// Request body
	bodyTag := tagFromType(*new(B))
	if (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) && bodyTag != "unknown-interface" && bodyTag != "string" {

		bodySchema, ok := s.OpenApiSpec.Components.Schemas[bodyTag]
		if !ok {
			bodySchema = openapi3.ToSchema(*new(B))
		}

		requestBody := &openapi3.RequestBody{
			Required: true,
			Content:  make(map[openapi3.MimeType]openapi3.SchemaObject),
		}

		if bodySchema != nil {
			requestBody.Content["application/json"] = openapi3.SchemaObject{
				Schema: bodySchema,
			}
		}

		s.OpenApiSpec.Components.RequestBodies[bodyTag] = requestBody

		// add request body to operation
		operation.RequestBody = requestBody
	}

	// Response body
	responseSchema, ok := s.OpenApiSpec.Components.Schemas[tag]
	if !ok {
		responseSchema = openapi3.ToSchema(*new(T))
	}

	operation.Responses = make(map[string]*openapi3.Response)
	operation.Responses["200"] = &openapi3.Response{
		Description: "OK",
		Content: map[openapi3.MimeType]openapi3.SchemaObject{
			"application/json": {
				Schema: responseSchema,
			},
		},
	}

	// Path parameters
	for _, pathParam := range parsePathParams(path) {
		operation.Parameters = append(operation.Parameters, &openapi3.Parameter{
			Name: pathParam,
			In:   "path",
			Schema: openapi3.Schema{
				Type: "string",
			},
		})
	}

	s.OpenApiSpec.Paths.AddPath(path, strings.ToLower(method), operation)

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
