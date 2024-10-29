package fuego

func ParamRequired() func(param *OpenAPIParam) {
	return func(param *OpenAPIParam) {
		param.Required = true
	}
}

func ParamNullable() func(param *OpenAPIParam) {
	return func(param *OpenAPIParam) {
		param.Nullable = true
	}
}

func ParamString() func(param *OpenAPIParam) {
	return func(param *OpenAPIParam) {
		param.GoType = "string"
	}
}

func ParamInteger() func(param *OpenAPIParam) {
	return func(param *OpenAPIParam) {
		param.GoType = "integer"
	}
}

func ParamBool() func(param *OpenAPIParam) {
	return func(param *OpenAPIParam) {
		param.GoType = "boolean"
	}
}

func ParamDescription(description string) func(param *OpenAPIParam) {
	return func(param *OpenAPIParam) {
		param.Description = description
	}
}

func ParamDefault(value any) func(param *OpenAPIParam) {
	return func(param *OpenAPIParam) {
		param.Default = value
	}
}

// Example adds an example to the parameter. As per the OpenAPI 3.0 standard, the example must be given a name.
func ParamExample(exampleName string, value any) func(param *OpenAPIParam) {
	return func(param *OpenAPIParam) {
		if param.Examples == nil {
			param.Examples = make(map[string]any)
		}
		param.Examples[exampleName] = value
	}
}
