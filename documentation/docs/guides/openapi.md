# OpenAPI Specification

Fuego automatically generates an OpenAPI specification for your API.

```go
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", helloWorld)

	s.Run()
}

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}
```

Result for this simple example:

![Swagger UI](../../static/img/hello-world-openapi.png)

The core idea of Fuego is to generate the OpenAPI specification automatically, so you don't have to worry about it. However, you can customize it if you want.

## Customize the operations caracteristics

Each route can be customized to add more information to the OpenAPI specification.

Just add methods after the route declaration.

```go
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", helloWorld).
		Summary("A simple hello world").
		Description("This is a simple hello world").
		Deprecated()

	s.Run()
}

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}

```

## Output

Fuego automatically provides an OpenAPI specification for your API in several ways:

- **JSON file**
- **Swagger UI**
- **JSON endpoint**

Fuego will indicate in a log the paths where the OpenAPI specifications and Swagger UI are available.

## Customize the OpenAPI specification output

You can customize the path and to activate or not the feature, with the option WithOpenAPIConfig.

```go
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer(fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
		DisableSwagger   : false, // If true, the server will not serve the swagger ui nor the openapi json spec
		DisableLocalSave : false, // If true, the server will not save the openapi json spec locally
		SwaggerUrl       : "/xxx", // URL to serve the swagger ui
		JsonUrl          : "/xxx/swagger.json", // URL to serve the openapi json spec
		JsonFilePath     : "./foo/bar.json", // Local path to save the openapi json spec
	}))

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	s.Run()
}
```
