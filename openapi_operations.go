package fuego

import (
	"slices"

	"github.com/getkin/kin-openapi/openapi3"
)

type ParamType string

const (
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
	Example  string
	Type     ParamType
}

// Overrides the description for the route.
func (r Route[ResponseBody, RequestBody]) Description(description string) Route[ResponseBody, RequestBody] {
	r.Operation.Description = description
	return r
}

// Overrides the summary for the route.
func (r Route[ResponseBody, RequestBody]) Summary(summary string) Route[ResponseBody, RequestBody] {
	r.Operation.Summary = summary
	return r
}

// Overrides the operationID for the route.
func (r Route[ResponseBody, RequestBody]) OperationID(operationID string) Route[ResponseBody, RequestBody] {
	r.Operation.OperationID = operationID
	return r
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

// Header registers a header parameter for the route.
func (r Route[ResponseBody, RequestBody]) Header(name, description string, params ...OpenAPIParamOption) Route[ResponseBody, RequestBody] {
	r.Param(HeaderParamType, name, description, params...)
	return r
}

// Cookie registers a cookie parameter for the route.
func (r Route[ResponseBody, RequestBody]) Cookie(name, description string, params ...OpenAPIParamOption) Route[ResponseBody, RequestBody] {
	r.Param(CookieParamType, name, description, params...)
	return r
}

// QueryParam registers a query parameter for the route.
func (r Route[ResponseBody, RequestBody]) QueryParam(name, description string, params ...OpenAPIParamOption) Route[ResponseBody, RequestBody] {
	r.Param(QueryParamType, name, description, params...)
	return r
}

// Replace the tags for the route.
// By default, the tag is the type of the response body.
func (r Route[ResponseBody, RequestBody]) Tags(tags ...string) Route[ResponseBody, RequestBody] {
	r.Operation.Tags = tags
	return r
}

// AddTags adds tags to the route.
func (r Route[ResponseBody, RequestBody]) AddTags(tags ...string) Route[ResponseBody, RequestBody] {
	r.Operation.Tags = append(r.Operation.Tags, tags...)
	return r
}

// AddError adds an error to the route.
func (r Route[ResponseBody, RequestBody]) AddError(code int, description string, errorType ...any) Route[ResponseBody, RequestBody] {
	addResponse(r.mainRouter, r.Operation, code, description, errorType...)
	return r
}

func addResponse(s *Server, operation *openapi3.Operation, code int, description string, errorType ...any) {
	var responseSchema schemaTag

	if len(errorType) > 0 {
		responseSchema = schemaTagFromType(s, errorType[0])
	} else {
		responseSchema = schemaTagFromType(s, HTTPError{})
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

// RemoveTags removes tags from the route.
func (r Route[ResponseBody, RequestBody]) RemoveTags(tags ...string) Route[ResponseBody, RequestBody] {
	for _, tag := range tags {
		for i, t := range r.Operation.Tags {
			if t == tag {
				r.Operation.Tags = slices.Delete(r.Operation.Tags, i, i+1)
				break
			}
		}
	}
	return r
}

func (r Route[ResponseBody, RequestBody]) Deprecated() Route[ResponseBody, RequestBody] {
	r.Operation.Deprecated = true
	return r
}
