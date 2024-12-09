package fuego

import (
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
)

type ParamType string // Query, Header, Cookie

const (
	PathParamType   ParamType = "path"
	QueryParamType  ParamType = "query"
	HeaderParamType ParamType = "header"
	CookieParamType ParamType = "cookie"
)

type OpenAPIParam struct {
	Name        string
	Description string
	OpenAPIParamOption
}

type OpenAPIParamOption struct {
	Required bool
	Nullable bool
	// Default value for the parameter.
	// Type is checked at start-time.
	Default  any
	Example  string
	Examples map[string]any
	Type     ParamType
	// integer, string, bool
	GoType string
	// Status codes for which this parameter is required.
	// Only used for response parameters.
	// If empty, it is required for 200 status codes.
	StatusCodes []int
}

// Param registers a parameter for the route.
// The paramType can be "query", "header" or "cookie" as defined in [ParamType].
// [Cookie], [Header], [QueryParam] are shortcuts for Param.
func (r Route[ResponseBody, RequestBody]) Param(paramType ParamType, name, description string, params ...OpenAPIParamOption) Route[ResponseBody, RequestBody] {
	openapiParam := openapi3.NewHeaderParameter(name)
	openapiParam.Description = description
	openapiParam.Schema = openapi3.NewStringSchema().NewRef()
	openapiParam.In = string(paramType)

	for _, param := range params {
		if param.Required {
			openapiParam.Required = param.Required
		}
		if param.Example != "" {
			openapiParam.Example = param.Example
		}
	}

	r.Operation.AddParameter(openapiParam)

	return r
}

// Registers a response for the route, only if error for this code is not already set.
func addResponseIfNotSet(s *Server, operation *openapi3.Operation, code int, description string, errorType ...any) {
	var responseSchema SchemaTag

	if len(errorType) > 0 {
		responseSchema = SchemaTagFromType(s, errorType[0])
	} else {
		responseSchema = SchemaTagFromType(s, HTTPError{})
	}
	content := openapi3.NewContentWithSchemaRef(&responseSchema.SchemaRef, []string{"application/json"})

	if operation.Responses.Value(strconv.Itoa(code)) == nil {
		response := openapi3.NewResponse().
			WithDescription(description).
			WithContent(content)

		operation.AddResponse(code, response)
	}
}

// openAPIError describes a response error in the OpenAPI spec.
type openAPIError struct {
	Code        int
	Description string
	ErrorType   any
}
