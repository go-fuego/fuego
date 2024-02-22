package fuego

import (
	"slices"

	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPIParam struct {
	Required bool
	Example  string
	Type     string // "query", "header", "cookie"
}

func (r Route[ResponseBody, RequestBody]) Description(description string) Route[ResponseBody, RequestBody] {
	r.operation.Description = description
	return r
}

func (r Route[ResponseBody, RequestBody]) Summary(summary string) Route[ResponseBody, RequestBody] {
	r.operation.Summary = summary
	return r
}

func (r Route[ResponseBody, RequestBody]) OperationID(operationID string) Route[ResponseBody, RequestBody] {
	r.operation.OperationID = operationID
	return r
}

// Param registers a parameter for the route.
// The paramType can be "query", "header" or "cookie".
// [Cookie], [Header], [QueryParam] are shortcuts for Param.
func (r Route[ResponseBody, RequestBody]) Param(paramType, name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	openapiParam := openapi3.NewHeaderParameter(name)
	openapiParam.Description = description
	openapiParam.Schema = openapi3.NewStringSchema().NewRef()
	openapiParam.In = paramType

	for _, param := range params {
		if param.Required {
			openapiParam.Required = param.Required
		}
		if param.Example != "" {
			openapiParam.Example = param.Example
		}
	}

	r.operation.AddParameter(openapiParam)

	return r
}

// Header registers a header parameter for the route.
func (r Route[ResponseBody, RequestBody]) Header(name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	r.Param("header", name, description, params...)
	return r
}

// Cookie registers a cookie parameter for the route.
func (r Route[ResponseBody, RequestBody]) Cookie(name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	r.Param("cookie", name, description, params...)
	return r
}

// QueryParam registers a query parameter for the route.
func (r Route[ResponseBody, RequestBody]) QueryParam(name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	r.Param("query", name, description, params...)
	return r
}

func (r Route[ResponseBody, RequestBody]) Tags(tags ...string) Route[ResponseBody, RequestBody] {
	r.operation.Tags = tags
	return r
}

func (r Route[ResponseBody, RequestBody]) AddTags(tags ...string) Route[ResponseBody, RequestBody] {
	r.operation.Tags = append(r.operation.Tags, tags...)
	return r
}

func (r Route[ResponseBody, RequestBody]) RemoveTags(tags ...string) Route[ResponseBody, RequestBody] {
	for _, tag := range tags {
		for i, t := range r.operation.Tags {
			if t == tag {
				r.operation.Tags = slices.Delete(r.operation.Tags, i, i+1)
				break
			}
		}
	}
	return r
}

func (r Route[ResponseBody, RequestBody]) Deprecated() Route[ResponseBody, RequestBody] {
	r.operation.Deprecated = true
	return r
}
