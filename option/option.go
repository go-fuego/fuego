package option

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/param"
)

// Group allows to group routes under a common path.
// Useful to group often used middlewares or options and reuse them.
// Example:
//
//	optionsPagination := option.Group(
//		option.QueryInt("per_page", "Number of items per page", param.Required()),
//		option.QueryInt("page", "Page number", param.Default(1)),
//	)
func Group(options ...func(*fuego.BaseRoute)) func(*fuego.BaseRoute) {
	return func(r *fuego.BaseRoute) {
		for _, option := range options {
			option(r)
		}
	}
}

// Middleware adds one or more route-scoped middleware.
func Middleware(middleware ...func(http.Handler) http.Handler) func(*fuego.BaseRoute) {
	return func(r *fuego.BaseRoute) {
		r.Middlewares = append(r.Middlewares, middleware...)
	}
}

// Declare a query parameter for the route.
// This will be added to the OpenAPI spec.
// Example:
//
//	Query("name", "Filter by name", param.Example("cat name", "felix"), param.Nullable())
//
// The list of options is in the param package.
func Query(name, description string, options ...func(*fuego.OpenAPIParam)) func(*fuego.BaseRoute) {
	options = append(options, param.Description(description), paramType(fuego.QueryParamType), param.String())
	return func(r *fuego.BaseRoute) {
		Param(name, options...)(r)
	}
}

// Declare an integer query parameter for the route.
// This will be added to the OpenAPI spec.
// The query parameter is transmitted as a string in the URL, but it is parsed as an integer.
// Example:
//
//	QueryInt("age", "Filter by age (in years)", param.Example("3 years old", 3), param.Nullable())
//
// The list of options is in the param package.
func QueryInt(name, description string, options ...func(*fuego.OpenAPIParam)) func(*fuego.BaseRoute) {
	options = append(options, param.Description(description), paramType(fuego.QueryParamType), param.Integer())
	return func(r *fuego.BaseRoute) {
		Param(name, options...)(r)
	}
}

// Declare a boolean query parameter for the route.
// This will be added to the OpenAPI spec.
// The query parameter is transmitted as a string in the URL, but it is parsed as a boolean.
// Example:
//
//	QueryBool("is_active", "Filter by active status", param.Example("true", true), param.Nullable())
//
// The list of options is in the param package.
func QueryBool(name, description string, options ...func(*fuego.OpenAPIParam)) func(*fuego.BaseRoute) {
	options = append(options, param.Description(description), paramType(fuego.QueryParamType), param.Bool())
	return func(r *fuego.BaseRoute) {
		Param(name, options...)(r)
	}
}

// Declare a header parameter for the route.
// This will be added to the OpenAPI spec.
// Example:
//
//	Header("Authorization", "Bearer token", param.Required())
//
// The list of options is in the param package.
func Header(name, description string, options ...func(*fuego.OpenAPIParam)) func(*fuego.BaseRoute) {
	options = append(options, param.Description(description), paramType(fuego.HeaderParamType))
	return func(r *fuego.BaseRoute) {
		Param(name, options...)(r)
	}
}

// Declare a cookie parameter for the route.
// This will be added to the OpenAPI spec.
// Example:
//
//	Cookie("session_id", "Session ID", param.Required())
//
// The list of options is in the param package.
func Cookie(name, description string, options ...func(*fuego.OpenAPIParam)) func(*fuego.BaseRoute) {
	options = append(options, param.Description(description), paramType(fuego.CookieParamType))
	return func(r *fuego.BaseRoute) {
		Param(name, options...)(r)
	}
}

func paramType(paramType fuego.ParamType) func(*fuego.OpenAPIParam) {
	return func(param *fuego.OpenAPIParam) {
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

// Registers a parameter for the route. Prefer using the [Query], [QueryInt], [Header], [Cookie] shortcuts.
func Param(name string, options ...func(*fuego.OpenAPIParam)) func(*fuego.BaseRoute) {
	param := fuego.OpenAPIParam{
		Name: name,
	}
	// Applies options to fuego.OpenAPIParam
	for _, option := range options {
		option(&param)
	}

	// Applies fuego.OpenAPIParam to openapi3.Parameter
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

	return func(r *fuego.BaseRoute) {
		r.Operation.AddParameter(openapiParam)
		if r.Params == nil {
			r.Params = make(map[string]fuego.OpenAPIParam)
		}
		r.Params[name] = param
	}
}

// Tags adds one or more tags to the route.
func Tags(tags ...string) func(*fuego.BaseRoute) {
	return func(r *fuego.BaseRoute) {
		r.Operation.Tags = append(r.Operation.Tags, tags...)
	}
}

// Summary adds a summary to the route.
func Summary(summary string) func(*fuego.BaseRoute) {
	return func(r *fuego.BaseRoute) {
		r.Operation.Summary = summary
	}
}

// Description adds a description to the route.
func Description(description string) func(*fuego.BaseRoute) {
	return func(r *fuego.BaseRoute) {
		r.Operation.Description = description
	}
}

// OperationID adds an operation ID to the route.
func OperationID(operationID string) func(*fuego.BaseRoute) {
	return func(r *fuego.BaseRoute) {
		r.Operation.OperationID = operationID
	}
}

// Deprecated marks the route as deprecated.
func Deprecated() func(*fuego.BaseRoute) {
	return func(r *fuego.BaseRoute) {
		r.Operation.Deprecated = true
	}
}

// AddError adds an error to the route.
func AddError(code int, description string, errorType ...any) func(*fuego.BaseRoute) {
	var responseSchema fuego.SchemaTag
	return func(r *fuego.BaseRoute) {
		if len(errorType) > 0 {
			responseSchema = fuego.SchemaTagFromType(r.MainRouter, errorType[0])
		} else {
			responseSchema = fuego.SchemaTagFromType(r.MainRouter, fuego.HTTPError{})
		}
		content := openapi3.NewContentWithSchemaRef(&responseSchema.SchemaRef, []string{"application/json"})

		response := openapi3.NewResponse().
			WithDescription(description).
			WithContent(content)
		r.Operation.AddResponse(code, response)
	}
}

// RequestContentType sets the accepted content types for the route.
// By default, the accepted content types is */*.
func RequestContentType(consumes ...string) func(*fuego.BaseRoute) {
	return func(r *fuego.BaseRoute) {
		r.AcceptedContentTypes = consumes
	}
}
