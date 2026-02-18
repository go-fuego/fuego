package fuego

// A ParamOption configures OpenAPI properties of [OpenAPIParam]
// i.e path/query parameters, cookies, and headers
type ParamOption = func(param *OpenAPIParam)

func ParamRequired() ParamOption {
	return func(param *OpenAPIParam) {
		param.Required = true
	}
}

func ParamNullable() ParamOption {
	return func(param *OpenAPIParam) {
		param.Nullable = true
	}
}

func ParamArray() ParamOption {
	return func(param *OpenAPIParam) {
		param.GoType = "array"
	}
}

func ParamString() ParamOption {
	return func(param *OpenAPIParam) {
		param.GoType = "string"
	}
}

func ParamInteger() ParamOption {
	return func(param *OpenAPIParam) {
		param.GoType = "integer"
	}
}

func ParamBool() ParamOption {
	return func(param *OpenAPIParam) {
		param.GoType = "boolean"
	}
}

func ParamDescription(description string) ParamOption {
	return func(param *OpenAPIParam) {
		param.Description = description
	}
}

func ParamDefault(value any) ParamOption {
	return func(param *OpenAPIParam) {
		param.Default = value
	}
}

// ParamExample adds an example to the parameter. As per the OpenAPI 3.0 standard, the example must be given a name.
func ParamExample(exampleName string, value any) ParamOption {
	return func(param *OpenAPIParam) {
		if param.Examples == nil {
			param.Examples = make(map[string]any)
		}
		param.Examples[exampleName] = value
	}
}

// ParamStatusCodes sets the status codes for which this parameter is required.
// Only used for response parameters.
// If empty, it is required for 200 status codes.
func ParamStatusCodes(codes ...int) ParamOption {
	return func(param *OpenAPIParam) {
		param.StatusCodes = codes
	}
}
