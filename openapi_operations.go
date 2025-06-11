package fuego

import (
	"reflect"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/go-fuego/fuego/internal"
)

type ParamType = internal.ParamType // Query, Header, Cookie

const (
	PathParamType   ParamType = "path"
	QueryParamType  ParamType = "query"
	HeaderParamType ParamType = "header"
	CookieParamType ParamType = "cookie"
)

type OpenAPIParam = internal.OpenAPIParam

// Registers a response for the route, only if error for this code is not already set.
func addResponseIfNotSet(openapi *OpenAPI, operation *openapi3.Operation, code int, description string, response Response) {
	if operation.Responses.Value(strconv.Itoa(code)) != nil {
		return
	}
	operation.AddResponse(code, openapi.buildOpenapi3Response(description, response))
}

func (openAPI *OpenAPI) buildOpenapi3Response(description string, response Response) *openapi3.Response {
	return openapi3.NewResponse().
		WithDescription(description).
		WithContent(openAPI.buildContent(response.Type, response.ContentTypes...))
}

func (openAPI *OpenAPI) buildOpenapi3RequestBody(requestBody RequestBody) *openapi3.RequestBody {
	return openapi3.NewRequestBody().
		WithRequired(true).
		WithDescription("Request body for " + reflect.TypeOf(requestBody.Type).String()).
		WithContent(openAPI.buildContent(requestBody.Type, requestBody.ContentTypes...))
}

func (openAPI *OpenAPI) buildContent(t any, consumes ...string) openapi3.Content {
	if t == nil {
		panic("Type in RequestBody cannot be nil")
	}
	if len(consumes) == 0 {
		consumes = []string{"application/json", "application/xml"}
	}

	bodySchema := SchemaTagFromType(openAPI, t)
	return openapi3.NewContentWithSchemaRef(&bodySchema.SchemaRef, consumes)
}

// openAPIResponse describes a response error in the OpenAPI spec.
type openAPIResponse struct {
	Response
	Description string
	Code        int
}
