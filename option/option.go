package option

import (
	"github.com/go-fuego/fuego"
)

// Group allows to group routes under a common path.
// Useful to group often used middlewares or options and reuse them.
// Example:
//
//	optionsPagination := option.Group(
//		option.QueryInt("per_page", "Number of items per page", param.Required()),
//		option.QueryInt("page", "Page number", param.Default(1)),
//	)
var Group = fuego.GroupOptions

// Middleware adds one or more route-scoped middleware.
var Middleware = fuego.OptionMiddleware

// Declare a query parameter for the route.
// This will be added to the OpenAPI spec.
// Example:
//
//	Query("name", "Filter by name", param.Example("cat name", "felix"), param.Nullable())
//
// The list of options is in the param package.
var Query = fuego.OptionQuery

// Declare an integer query parameter for the route.
// This will be added to the OpenAPI spec.
// The query parameter is transmitted as a string in the URL, but it is parsed as an integer.
// Example:
//
//	QueryInt("age", "Filter by age (in years)", param.Example("3 years old", 3), param.Nullable())
//
// The list of options is in the param package.
var QueryInt = fuego.OptionQueryInt

// Declare a boolean query parameter for the route.
// This will be added to the OpenAPI spec.
// The query parameter is transmitted as a string in the URL, but it is parsed as a boolean.
// Example:
//
//	QueryBool("is_active", "Filter by active status", param.Example("true", true), param.Nullable())
//
// The list of options is in the param package.
var QueryBool = fuego.OptionQueryBool

// Declare a header parameter for the route.
// This will be added to the OpenAPI spec.
// Example:
//
//	Header("Authorization", "Bearer token", param.Required())
//
// The list of options is in the param package.
var Header = fuego.OptionHeader

// Declare a cookie parameter for the route.
// This will be added to the OpenAPI spec.
// Example:
//
//	Cookie("session_id", "Session ID", param.Required())
//
// The list of options is in the param package.
var Cookie = fuego.OptionCookie

// Declare a path parameter for the route.
// This will be added to the OpenAPI spec.
// It will be marked as required by default by Fuego.
// If not set explicitly, the parameter will still be declared on the spec.
// Example:
//
//	Path("id", "ID of the item", param.Required())
//
// The list of options is in the param package.
var Path = fuego.OptionPath

// Registers a parameter for the route. Prefer using the [Query], [QueryInt], [Header], [Cookie] shortcuts.
var Param = fuego.OptionParam

// Tags adds one or more tags to the route.
var Tags = fuego.OptionTags

// Summary adds a summary to the route.
var Summary = fuego.OptionSummary

// Description adds a description to the route.
var Description = fuego.OptionDescription

// Security configures security requirements to the route.
//
// Single Scheme (AND Logic):
//
//	Add a single security requirement with multiple schemes.
//	All schemes must be satisfied:
//	Security(openapi3.SecurityRequirement{
//	  "basic": [],        // Requires basic auth
//	  "oauth2": ["read"]  // AND requires oauth with read scope
//	})
//
// Multiple Schemes (OR Logic):
//
//	Add multiple security requirements.
//	At least one requirement must be satisfied:
//	Security(
//	  openapi3.SecurityRequirement{"basic": []},        // First option
//	  openapi3.SecurityRequirement{"oauth2": ["read"]}  // Alternative option
//	})
//
// Mixing Approaches:
//
//	Combine AND logic within requirements and OR logic between requirements:
//	Security(
//	  openapi3.SecurityRequirement{
//	    "basic": [],             // Requires basic auth
//	    "oauth2": ["read:user"]  // AND oauth with read:user scope
//	  },
//	  openapi3.SecurityRequirement{"apiKey": []}  // OR alternative with API key
//	})
var Security = fuego.OptionSecurity

// OperationID adds an operation ID to the route.
var OperationID = fuego.OptionOperationID

// Deprecated marks the route as deprecated.
var Deprecated = fuego.OptionDeprecated

// AddError adds an error to the route.
var AddError = fuego.OptionAddError

// RequestContentType sets the accepted content types for the route.
// By default, the accepted content types is */*.
// This will override any options set at the server level.
var RequestContentType = fuego.OptionRequestContentType

// Hide hides the route from the OpenAPI spec.
var Hide = fuego.OptionHide
