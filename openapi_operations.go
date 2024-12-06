package fuego

import (
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
	Default  any // Default value for the parameter
	Example  string
	Examples map[string]any
	Type     ParamType
	GoType   string // integer, string, bool
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

// AddError adds an error to the route.
//
// Deprecated: Use `option.AddError` from github.com/go-fuego/fuego/option instead.
func (r Route[ResponseBody, RequestBody]) AddError(code int, description string, errorType ...any) Route[ResponseBody, RequestBody] {
	addResponse(r.mainRouter, r.Operation, code, description, errorType...)
	return r
}

func addResponse(s *Server, operation *openapi3.Operation, code int, description string, errorType ...any) {
	var responseSchema SchemaTag

	if len(errorType) > 0 {
		responseSchema = SchemaTagFromType(s, errorType[0])
	} else {
		responseSchema = SchemaTagFromType(s, HTTPError{})
	}
	content := openapi3.NewContentWithSchemaRef(&responseSchema.SchemaRef, []string{"application/json"})

	response := openapi3.NewResponse().
		WithDescription(description).
		WithContent(content)

	operation.AddResponse(code, response)
}

// openAPIError describes a response error in the OpenAPI spec.
type openAPIError struct {
	Code        int
	Description string
	ErrorType   any
}
