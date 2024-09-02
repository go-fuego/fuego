package param

import "github.com/go-fuego/fuego"

func Required() func(param *fuego.OpenAPIParam) {
	return func(param *fuego.OpenAPIParam) {
		param.Required = true
	}
}

func Nullable() func(param *fuego.OpenAPIParam) {
	return func(param *fuego.OpenAPIParam) {
		param.Nullable = true
	}
}

func Integer() func(param *fuego.OpenAPIParam) {
	return func(param *fuego.OpenAPIParam) {
		param.GoType = "integer"
	}
}

func Bool() func(param *fuego.OpenAPIParam) {
	return func(param *fuego.OpenAPIParam) {
		param.GoType = "boolean"
	}
}

func Description(description string) func(param *fuego.OpenAPIParam) {
	return func(param *fuego.OpenAPIParam) {
		param.Description = description
	}
}

func Default(value any) func(param *fuego.OpenAPIParam) {
	return func(param *fuego.OpenAPIParam) {
		param.Default = value
	}
}

// Example adds an example to the parameter. As per the OpenAPI 3.0 standard, the example must be given a name.
func Example(exampleName string, value any) func(param *fuego.OpenAPIParam) {
	return func(param *fuego.OpenAPIParam) {
		if param.Examples == nil {
			param.Examples = make(map[string]any)
		}
		param.Examples[exampleName] = value
	}
}
