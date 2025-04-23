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

Result for this simple example at [http://localhost:9999/swagger/index.html](http://localhost:9999/swagger/index.html):

![Swagger UI](../../static/img/hello-world-openapi.jpeg)

The core idea of Fuego is to generate the OpenAPI specification automatically,
so you don't have to worry about it. However, you can customize it if you want.

## Route Options

Each route can be customized to add more information to the OpenAPI specification.

Just add options after the route declaration.

```go
package main

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", helloWorld,
		option.Summary("A simple hello world"),
		option.Description("This is a simple hello world example"),
		option.Query("name", "Name to greet", param.Required(), param.Default("World")),
		option.Tags("Hello"),
		option.Deprecated(),
	)

	s.Run()
}

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}
```

## Group Options, Options Groups & Custom Options

You can also customize the OpenAPI specification for a group of routes.

```go
package main

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

// Define a reusable group of options
var optionPagination = option.Group(
	option.QueryInt("page", "Page number", param.Default(1)),
	option.QueryInt("limit", "Items per page", param.Default(10)),
)

// Custom options for the group
var customOption = func(r *fuego.BaseRoute) {
	r.XXX = YYY // Direct access to the route struct to inject custom behavior
}

func main() {
	s := fuego.NewServer()

	api := fuego.Group(s, "/users",
		option.Summary("Users routes"),
		option.Description("Default description for all Users routes"),
		option.Tags("users"),
	)

	fuego.Get(api, "/", helloWorld,
		optionPagination,
		customOption,
		option.Summary("A simple hello world"), // Replace the default summary
	)

	s.Run()
}

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}
```

## Output

Fuego automatically provides an OpenAPI specification for your API in several ways:

- **JSON file** - Saved locally at the specified path
- **OpenAPI JSON endpoint** - Available at the specified URL
- **Swagger UI** - Interactive documentation UI available at the specified URL

Fuego will indicate in a log the paths where the OpenAPI specifications and
Swagger UI are available.

You can customize the paths and activate or deactivate these features with the `WithOpenAPIConfig` option.

```go
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer(
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				Disabled:          false,                   // If true, the server will not serve nor generate any OpenAPI resources
				DisableSwaggerUI:  false,                   // If true, the server will not serve the swagger ui
				DisableLocalSave:  false,                   // If true, the server will not save the openapi json spec locally
				DisableMessages:   false,                   // If true, the engine will not print messages
				PrettyFormatJSON:  true,                    // Pretty prints the OpenAPI spec with proper JSON indentation
				SwaggerURL:        "/swagger",              // URL to serve the swagger ui
				SpecURL:           "/swagger/openapi.json", // URL to serve the openapi json spec
				JSONFilePath:      "doc/openapi.json",      // Local path to save the openapi json spec
				UIHandler:         fuego.DefaultOpenAPIHandler, // Custom UI handler
			}),
		),
	)

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	s.Run()
}
```

## Custom UI

By default, Fuego uses [Stoplight Elements](https://stoplight.io/open-source/elements) for the OpenAPI UI.
You can customize this by providing your own `UIHandler` in the OpenAPIConfig.

Example with `http-swagger`:

```go
import (
	"net/http"

	"github.com/go-fuego/fuego"
	httpSwagger "github.com/swaggo/http-swagger"
)

func openApiHandler(specURL string) http.Handler {
	return httpSwagger.Handler(
		httpSwagger.Layout(httpSwagger.BaseLayout),
		httpSwagger.PersistAuthorization(true),
		httpSwagger.URL(specURL), // The url pointing to API definition
	)
}

func main() {
	s := fuego.NewServer(
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				UIHandler: openApiHandler,
			}),
		),
	)

	fuego.Get(s, "/", helloWorld)

	s.Run()
}
```

Comparison of UI options:

|               | StopLight Elements (Default) | Swagger UI     | Disabled |
| ------------- | ---------------------------- | -------------- | -------- |
| Works offline | No ❌                        | Yes ✅         | -        |
| Binary Size   | Smaller                      | Larger (+10Mb) | Smallest |

## Get OpenAPI Spec at build time

With Go, you cannot generate things at build time, but you can separate the
OpenAPI generation from the server, by using the
`(Server).OutputOpenAPISpec()` function.

```go title="main.go" showLineNumbers
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", helloWorld)

	if justGenerateOpenAPI { // A CLI flag, an env variable, etc.
		s.OutputOpenAPISpec()
		return
	}

	s.Run()
}
```

## Hide From OpenAPI Spec

Certain routes such as web routes you may not want to be part of the OpenAPI spec.

You can prevent them from being added to the OpenAPI documentation using route options.

```go
package main

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func main() {
	s := fuego.NewServer()

	// Create a group of routes to be hidden
	web := fuego.Group(s, "/",
		option.Hide(), // Hide all routes in this group from OpenAPI
	)

	fuego.Get(web, "/", func(c fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	// These routes will still be added to the spec
	api := fuego.Group(s, "/api")
	fuego.Get(api, "/hello", func(c fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	s.Run()
}
```
