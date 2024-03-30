# Options

You can customize the server with the following function options.

All the options start with `With` and are located in the `fuego` package.

```go
package main

import "github.com/go-fuego/fuego"

func main() {
	s := fuego.NewServer(
		fuego.WithPort(8080),
		fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
			DisableSwagger   : true,
		}),
	)

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	s.Run()
}
```

## Some options

### Port

You can change the port of the server with the `WithPort` option.

```go
s := fuego.NewServer(
	fuego.WithPort(8080),
)
```

### CORS

CORS middleware is not registered as a usual middleware, because it applies on routes that aren't registered. For example, `OPTIONS /foo` is not a registered route (only `GET /foo` is registered for example), but it's a request that needs to be handled by the CORS middleware.

```go
import "github.com/rs/cors"

s := fuego.NewServer(
	fuego.WithCorsMiddleware(cors.New(cors.Options{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler),
)
```

### Custom OpenAPI generation

You can customize the used OpenAPI generation with the `WithOpenAPIGenerator` option.
Please refer to the [openapi3gen](https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3gen) package for more information.

```go
package main

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"

	"github.com/go-fuego/fuego"
)

type Input struct {
	Name string `json:"name" description:"some name" readOnly:"true"`
}

type Output struct {
	Input Input `json:"input"`
}

func main() {
	s := fuego.NewServer(
		fuego.WithOpenAPIGenerator(
			openapi3gen.NewGenerator(
				openapi3gen.UseAllExportedFields(),
				openapi3gen.SchemaCustomizer(func(name string, t reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {
					if v := tag.Get("readOnly"); v == "true" {
						schema.ReadOnly = true
					}

					if v := tag.Get("description"); v != "" {
						schema.Description = v
					}

					return nil
				}),
			),
		),
	)

	// add routes

	s.Run()
}

```
