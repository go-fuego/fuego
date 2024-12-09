package param

import "github.com/go-fuego/fuego"

// Required sets the parameter as required.
// If the parameter is not present, the request will fail.
var Required = fuego.ParamRequired

// Nullable sets the parameter as nullable.
var Nullable = fuego.ParamNullable

// Integer sets the parameter type to integer.
// The query parameter is transmitted as a string in the URL, but it is parsed as an integer.
// Please prefer QueryInt for clarity.
var Integer = fuego.ParamInteger

// Bool sets the parameter type to boolean.
// The query parameter is transmitted as a string in the URL, but it is parsed as a boolean.
// Please prefer QueryBool for clarity.
var Bool = fuego.ParamBool

// Description sets the description for the parameter.
var Description = fuego.ParamDescription

// Default sets the default value for the parameter.
// Type is checked at start-time.
var Default = fuego.ParamDefault

// Example adds an example to the parameter. As per the OpenAPI 3.0 standard, the example must be given a name.
var Example = fuego.ParamExample

// StatusCodes sets the status codes for which this parameter is required.
// Only used for response parameters.
// If empty, it is required for 200 status codes.
var StatusCodes = fuego.ParamStatusCodes
