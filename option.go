package fuego

import (
	"fmt"
	"maps"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// GroupOptions allows to group routes under a common path.
// Useful to group often used middlewares or options and reuse them.
// Example:
//
//	optionsPagination := GroupOptions(
//		OptionQueryInt("per_page", "Number of items per page", ParamRequired()),
//		OptionQueryInt("page", "Page number", ParamDefault(1)),
//	)
func GroupOptions(options ...func(*BaseRoute)) func(*BaseRoute) {
	return func(r *BaseRoute) {
		for _, option := range options {
			option(r)
		}
	}
}

// OptionMiddleware adds one or more route-scoped middleware.
func OptionMiddleware(middleware ...func(http.Handler) http.Handler) func(*BaseRoute) {
	return func(r *BaseRoute) {
		r.Middlewares = append(r.Middlewares, middleware...)
	}
}

// OptionQuery declares a query parameter for the route.
// This will be added to the OpenAPI spec.
// Example:
//
//	OptionQuery("name", "Filter by name", ParamExample("cat name", "felix"), ParamNullable())
//
// The list of options is in the param package.
func OptionQuery(name, description string, options ...ParamOption) func(*BaseRoute) {
	options = append(options, ParamDescription(description), paramType(QueryParamType), ParamString())
	return func(r *BaseRoute) {
		OptionParam(name, options...)(r)
	}
}

// OptionQueryInt declares an integer query parameter for the route.
// This will be added to the OpenAPI spec.
// The query parameter is transmitted as a string in the URL, but it is parsed as an integer.
// Example:
//
//	OptionQueryInt("age", "Filter by age (in years)", ParamExample("3 years old", 3), ParamNullable())
//
// The list of options is in the param package.
func OptionQueryInt(name, description string, options ...ParamOption) func(*BaseRoute) {
	options = append(options, ParamDescription(description), paramType(QueryParamType), ParamInteger())
	return func(r *BaseRoute) {
		OptionParam(name, options...)(r)
	}
}

// OptionQueryBool declares a boolean query parameter for the route.
// This will be added to the OpenAPI spec.
// The query parameter is transmitted as a string in the URL, but it is parsed as a boolean.
// Example:
//
//	OptionQueryBool("is_active", "Filter by active status", ParamExample("true", true), ParamNullable())
//
// The list of options is in the param package.
func OptionQueryBool(name, description string, options ...ParamOption) func(*BaseRoute) {
	options = append(options, ParamDescription(description), paramType(QueryParamType), ParamBool())
	return func(r *BaseRoute) {
		OptionParam(name, options...)(r)
	}
}

// OptionQueryArray declares a query parameter for the route with an array type.
// This will be added to the OpenAPI spec.
// Example:
//
//	OptionQueryArray("tags", "Filter by tags", reflect.Int, ParamExample("tag list", "1,2,3"))
//
// The list of options is in the param package.
func OptionQueryArray(name, description string, elemKind reflect.Kind, options ...ParamOption) func(*BaseRoute) {
	return func(r *BaseRoute) {
		param, openapiParam := buildParam(name, append(options, ParamDescription(description), paramType(QueryParamType))...)

		// Create an array schema
		arraySchema := openapi3.NewSchema()
		arraySchema.Type = &openapi3.Types{openapi3.TypeArray}

		// Set the items type based on the element kind
		itemsSchema := openapi3.NewSchema()
		switch elemKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			itemsSchema.Type = &openapi3.Types{openapi3.TypeInteger}
		case reflect.Float32, reflect.Float64:
			itemsSchema.Type = &openapi3.Types{openapi3.TypeNumber}
		case reflect.Bool:
			itemsSchema.Type = &openapi3.Types{openapi3.TypeBoolean}
		default:
			itemsSchema.Type = &openapi3.Types{openapi3.TypeString}
		}

		arraySchema.Items = &openapi3.SchemaRef{Value: itemsSchema}
		openapiParam.Schema = &openapi3.SchemaRef{Value: arraySchema}

		r.Operation.AddParameter(openapiParam)
		if r.Params == nil {
			r.Params = make(map[string]OpenAPIParam)
		}
		r.Params[name] = param
	}
}

// OptionHeader declares a header parameter for the route.
// This will be added to the OpenAPI spec.
// Example:
//
//	OptionHeader("Authorization", "Bearer token", ParamRequired())
//
// The list of options is in the param package.
func OptionHeader(name, description string, options ...ParamOption) func(*BaseRoute) {
	options = append(options, ParamDescription(description), paramType(HeaderParamType))
	return func(r *BaseRoute) {
		OptionParam(name, options...)(r)
	}
}

// OptionCookie declares a cookie parameter for the route.
// This will be added to the OpenAPI spec.
// Example:
//
//	OptionCookie("session_id", "Session ID", ParamRequired())
//
// The list of options is in the param package.
func OptionCookie(name, description string, options ...ParamOption) func(*BaseRoute) {
	options = append(options, ParamDescription(description), paramType(CookieParamType))
	return func(r *BaseRoute) {
		OptionParam(name, options...)(r)
	}
}

// OptionPath declares a path parameter for the route.
// This will be added to the OpenAPI spec.
// It will be marked as required by default by Fuego.
// Example:
//
//	OptionPath("id", "ID of the item")
//
// The list of options is in the param package.
func OptionPath(name, description string, options ...ParamOption) func(*BaseRoute) {
	options = append(options, ParamDescription(description), paramType(PathParamType), ParamRequired())
	return func(r *BaseRoute) {
		OptionParam(name, options...)(r)
	}
}

func paramType(paramType ParamType) func(*OpenAPIParam) {
	return func(param *OpenAPIParam) {
		param.Type = paramType
	}
}

func panicsIfNotCorrectType(openapiParam *openapi3.Parameter, exampleValue any) any {
	if exampleValue == nil {
		return nil
	}
	if openapiParam.Schema.Value.Type.Is("integer") {
		_, ok := exampleValue.(int)
		if !ok {
			panic("example value must be an integer")
		}
	}
	if openapiParam.Schema.Value.Type.Is("boolean") {
		_, ok := exampleValue.(bool)
		if !ok {
			panic("example value must be a boolean")
		}
	}
	if openapiParam.Schema.Value.Type.Is("string") {
		_, ok := exampleValue.(string)
		if !ok {
			panic("example value must be a string")
		}
	}
	return exampleValue
}

// OptionResponseHeader declares a response header for the route.
// This will be added to the OpenAPI spec, under the given default status code response.
// Example:
//
//	OptionResponseHeader("Content-Range", "Pagination range", ParamExample("42 pets", "unit 0-9/42"), ParamDescription("https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Range"))
//	OptionResponseHeader("Set-Cookie", "Session cookie", ParamExample("session abc123", "session=abc123; Expires=Wed, 09 Jun 2021 10:18:14 GMT"))
//
// The list of options is in the param package.
func OptionResponseHeader(name, description string, options ...ParamOption) func(*BaseRoute) {
	apiParam, openapiParam := buildParam(name, options...)

	openapiParam.Name = ""
	openapiParam.In = ""

	if len(apiParam.StatusCodes) == 0 {
		apiParam.StatusCodes = []int{200}
	}

	return func(r *BaseRoute) {
		for _, code := range apiParam.StatusCodes {
			codeString := strconv.Itoa(code)
			responseForCurrentCode := r.Operation.Responses.Value(codeString)
			if responseForCurrentCode == nil {
				response := openapi3.NewResponse().WithDescription("OK")
				r.Operation.AddResponse(code, response)
				responseForCurrentCode = r.Operation.Responses.Value(codeString)
			}

			responseForCurrentCodeHeaders := responseForCurrentCode.Value.Headers
			if responseForCurrentCodeHeaders == nil {
				responseForCurrentCode.Value.Headers = make(map[string]*openapi3.HeaderRef)
			}

			responseForCurrentCode.Value.Headers[name] = &openapi3.HeaderRef{
				Value: &openapi3.Header{
					Parameter: *openapiParam,
				},
			}
		}
	}
}

func buildParam(name string, options ...ParamOption) (OpenAPIParam, *openapi3.Parameter) {
	param := OpenAPIParam{
		Name: name,
	}
	// Applies options to OpenAPIParam
	for _, option := range options {
		option(&param)
	}

	// Applies [OpenAPIParam] to [openapi3.Parameter]
	// Why not use openapi3.NewHeaderParameter(name) directly?
	// Because we might change the openapi3 library in the future,
	// and we want to keep the flexibility to change the implementation without changing the API.
	openapiParam := &openapi3.Parameter{
		Name:        name,
		In:          string(param.Type),
		Description: param.Description,
		Schema:      openapi3.NewStringSchema().NewRef(),
	}
	if param.GoType != "" {
		openapiParam.Schema.Value.Type = &openapi3.Types{param.GoType}
	}
	openapiParam.Schema.Value.Nullable = param.Nullable
	openapiParam.Schema.Value.Default = panicsIfNotCorrectType(openapiParam, param.Default)

	if param.Required {
		openapiParam.Required = param.Required
	}
	for name, exampleValue := range param.Examples {
		if openapiParam.Examples == nil {
			openapiParam.Examples = make(openapi3.Examples)
		}
		exampleOpenAPI := openapi3.NewExample(name)
		exampleOpenAPI.Value = panicsIfNotCorrectType(openapiParam, exampleValue)
		openapiParam.Examples[name] = &openapi3.ExampleRef{Value: exampleOpenAPI}
	}

	return param, openapiParam
}

// OptionParam registers a parameter for the route. Prefer using the [OptionQuery], [OptionQueryInt], [OptionHeader], [OptionCookie] shortcuts.
func OptionParam(name string, options ...ParamOption) func(*BaseRoute) {
	param, openapiParam := buildParam(name, options...)

	return func(r *BaseRoute) {
		r.Operation.AddParameter(openapiParam)
		if r.Params == nil {
			r.Params = make(map[string]OpenAPIParam)
		}
		r.Params[name] = param
	}
}

// OptionTags adds one or more tags to the route.
func OptionTags(tags ...string) func(*BaseRoute) {
	return func(r *BaseRoute) {
		for _, tag := range tags {
			if slices.Contains(r.Operation.Tags, tag) {
				return
			}
			r.Operation.Tags = append(r.Operation.Tags, tag)
		}
	}
}

// OptionTagInfo adds a tag with name and description to the route and registers it in the global OpenAPI tags.
// This ensures the tag appears in the OpenAPI spec with its description.
func OptionTagInfo(name, description string) func(*BaseRoute) {
	return func(r *BaseRoute) {
		if !slices.Contains(r.Operation.Tags, name) {
			r.Operation.Tags = append(r.Operation.Tags, name)
		}

		if r.OpenAPI.Description().Tags == nil {
			r.OpenAPI.Description().Tags = openapi3.Tags{}
		}

		existingTag := r.OpenAPI.Description().Tags.Get(name)
		if existingTag != nil {
			if existingTag.Description == "" {
				existingTag.Description = description
			}
			return
		}

		r.OpenAPI.Description().Tags = append(r.OpenAPI.Description().Tags, &openapi3.Tag{
			Name:        name,
			Description: description,
		})
	}
}

// OptionSummary adds a summary to the route.
func OptionSummary(summary string) func(*BaseRoute) {
	return func(r *BaseRoute) {
		r.Operation.Summary = summary
	}
}

// OptionDescription overrides the default description set by Fuego.
// By default, the description is set by Fuego with some info,
// like the controller function name and the package name.
func OptionDescription(description string) func(*BaseRoute) {
	return func(r *BaseRoute) {
		r.Operation.Description = description
	}
}

// OptionAddDescription appends a description to the route.
// By default, the description is set by Fuego with some info,
// like the controller function name and the package name.
func OptionAddDescription(description string) func(*BaseRoute) {
	return func(r *BaseRoute) {
		r.Operation.Description += description
	}
}

// OptionOverrideDescription overrides the default description set by Fuego.
// By default, the description is set by Fuego with some info,
// like the controller function name and the package name.
func OptionOverrideDescription(description string) func(*BaseRoute) {
	return func(r *BaseRoute) {
		r.overrideDescription = true
		r.Operation.Description = description
	}
}

// OptionOperationID adds an operation ID to the route.
func OptionOperationID(operationID string) func(*BaseRoute) {
	return func(r *BaseRoute) {
		r.Operation.OperationID = operationID
	}
}

// OptionDeprecated marks the route as deprecated.
func OptionDeprecated() func(*BaseRoute) {
	return func(r *BaseRoute) {
		r.Operation.Deprecated = true
	}
}

// OptionAddError adds an error to the route.
// It replaces any existing error previously set with the same code.
// Required: should only supply one type to `errorType`
// Deprecated: Use [OptionAddResponse] instead
func OptionAddError(code int, description string, errorType ...any) func(*BaseRoute) {
	var responseSchema SchemaTag
	return func(r *BaseRoute) {
		if len(errorType) > 1 {
			panic("errorType should not be more than one")
		}

		if len(errorType) > 0 {
			responseSchema = SchemaTagFromType(r.OpenAPI, errorType[0])
		} else {
			responseSchema = SchemaTagFromType(r.OpenAPI, HTTPError{})
		}
		content := openapi3.NewContentWithSchemaRef(&responseSchema.SchemaRef, []string{"application/json"})

		response := openapi3.NewResponse().
			WithDescription(description).
			WithContent(content)

		if r.Operation.Responses == nil {
			r.Operation.Responses = openapi3.NewResponses()
		}
		r.Operation.Responses.Set(strconv.Itoa(code), &openapi3.ResponseRef{Value: response})
	}
}

// Response represents a fuego.Response that can be used
// when setting custom response types on routes
type Response struct {
	// user provided type
	Type any
	// content-type of the response i.e application/json
	ContentTypes []string
}

// OptionAddResponse adds a response to a route by status code
// It replaces any existing response set by any status code, this will override 200.
// Required: Response.Type must be set
// Optional: Response.ContentTypes will default to `application/json` and `application/xml` if not set
func OptionAddResponse(code int, description string, response Response) func(*BaseRoute) {
	return func(r *BaseRoute) {
		if r.Operation.Responses == nil {
			r.Operation.Responses = openapi3.NewResponses()
		}
		r.Operation.Responses.Set(
			strconv.Itoa(code), &openapi3.ResponseRef{
				Value: r.OpenAPI.buildOpenapi3Response(description, response),
			},
		)
	}
}

// RequestBody represents a fuego.RequestBody that can be used
// when setting custom request type on routes
type RequestBody struct {
	// user provided type
	Type any
	// content-type of the request i.e application/json
	ContentTypes []string
}

// OptionRequestBody sets a request body to a route
// Can be called multiple times to set additional ContentType values.
// If two of the same ContentTypes are provided the last provided will be
// what is set in the spec.
// Required: RequestBody.Type must be set
// Optional: RequestBody.ContentTypes will default to `application/json` and `application/xml` if not set
func OptionRequestBody(requestBody RequestBody) func(*BaseRoute) {
	return func(r *BaseRoute) {
		if r.Operation.RequestBody == nil {
			r.Operation.RequestBody = &openapi3.RequestBodyRef{
				Value: r.OpenAPI.buildOpenapi3RequestBody(requestBody),
			}
			return
		}
		content := r.OpenAPI.buildContent(requestBody.Type, requestBody.ContentTypes...)
		maps.Copy(r.Operation.RequestBody.Value.Content, content)
	}
}

// OptionDefaultResponse adds a default response to a route
// Required: Response.Type must be set
// Optional: Response.ContentTypes will default to `application/json` and `application/xml` if not set
func OptionDefaultResponse(description string, response Response) func(*BaseRoute) {
	return func(r *BaseRoute) {
		if r.Operation.Responses == nil {
			r.Operation.Responses = openapi3.NewResponses()
		}
		r.Operation.Responses.Set(
			"default", &openapi3.ResponseRef{
				Value: r.OpenAPI.buildOpenapi3Response(description, response),
			},
		)
	}
}

// OptionRequestContentType sets the accepted content types for the route.
// By default, the accepted content types is */*.
// This will override any options set at the server level.
func OptionRequestContentType(consumes ...string) func(*BaseRoute) {
	return func(r *BaseRoute) {
		r.RequestContentTypes = consumes
	}
}

// OptionHide hides the route from the OpenAPI spec.
func OptionHide() func(*BaseRoute) {
	return func(r *BaseRoute) {
		r.Hidden = true
	}
}

// OptionShow shows the route from the OpenAPI spec.
func OptionShow() func(*BaseRoute) {
	return func(r *BaseRoute) {
		r.Hidden = false
	}
}

// OptionDefaultStatusCode sets the default status code for the route.
func OptionDefaultStatusCode(defaultStatusCode int) func(*BaseRoute) {
	return func(r *BaseRoute) {
		r.DefaultStatusCode = defaultStatusCode
	}
}

// OptionSecurity configures security requirements to the route.
//
// Single Scheme (AND Logic):
//
//	Add a single security requirement with multiple schemes.
//	All schemes must be satisfied:
//	OptionSecurity(openapi3.SecurityRequirement{
//	  "basic": [],        // Requires basic auth
//	  "oauth2": ["read"]  // AND requires oauth with read scope
//	})
//
// Multiple Schemes (OR Logic):
//
//	Add multiple security requirements.
//	At least one requirement must be satisfied:
//	OptionSecurity(
//	  openapi3.SecurityRequirement{"basic": []},        // First option
//	  openapi3.SecurityRequirement{"oauth2": ["read"]}  // Alternative option
//	})
//
// Mixing Approaches:
//
//	Combine AND logic within requirements and OR logic between requirements:
//	OptionSecurity(
//	  openapi3.SecurityRequirement{
//	    "basic": [],             // Requires basic auth
//	    "oauth2": ["read:user"]  // AND oauth with read:user scope
//	  },
//	  openapi3.SecurityRequirement{"apiKey": []}  // OR alternative with API key
//	})
func OptionSecurity(securityRequirements ...openapi3.SecurityRequirement) func(*BaseRoute) {
	return func(r *BaseRoute) {
		if r.OpenAPI.Description().Components == nil {
			panic("zero security schemes have been registered with the server")
		}

		// Validate the security scheme exists in components
		for _, req := range securityRequirements {
			for schemeName := range req {
				if _, exists := r.OpenAPI.Description().Components.SecuritySchemes[schemeName]; !exists {
					panic(fmt.Sprintf("security scheme '%s' not defined in components", schemeName))
				}
			}
		}

		if r.Operation.Security == nil {
			r.Operation.Security = &openapi3.SecurityRequirements{}
		}

		// Append all provided security requirements
		*r.Operation.Security = append(*r.Operation.Security, securityRequirements...)
	}
}

// OptionStripTrailingSlash ensure that the route declaration
// will have its ending trailing slash stripped.
func OptionStripTrailingSlash() func(*BaseRoute) {
	return func(r *BaseRoute) {
		if len(r.Path) > 1 {
			r.Path = strings.TrimRight(r.Path, "/")
		}
	}
}
