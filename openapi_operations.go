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
	if response.Type == nil {
		panic("Type in Response cannot be nil")
	}

	responseSchema := SchemaTagFromType(openAPI, response.Type)
	if len(response.ContentTypes) == 0 {
		response.ContentTypes = []string{"application/json", "application/xml"}
	}

	content := openapi3.NewContentWithSchemaRef(&responseSchema.SchemaRef, response.ContentTypes)
	return openapi3.NewResponse().
		WithDescription(description).
		WithContent(content)
}

func (openAPI *OpenAPI) buildOpenapi3RequestBody(requestBody RequestBody) *openapi3.RequestBody {
	if requestBody.Type == nil {
		panic("Type in RequestBody cannot be nil")
	}

	bodySchema := SchemaTagFromType(openAPI, requestBody.Type)
	if len(requestBody.ContentTypes) == 0 {
		requestBody.ContentTypes = []string{"application/json", "application/xml"}
	}

	content := openapi3.NewContentWithSchemaRef(&bodySchema.SchemaRef, requestBody.ContentTypes)
	return openapi3.NewRequestBody().
		WithRequired(true).
		WithDescription("Request body for " + reflect.TypeOf(requestBody.Type).String()).
		WithContent(content)
}

// openAPIResponse describes a response error in the OpenAPI spec.
type openAPIResponse struct {
	Description string
	Response
	Code int
}
