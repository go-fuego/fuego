# Options

## Server options

You can customize the server with the following function options.

All the options start with `With` and are located in the `fuego` package, see [the full list](https://pkg.go.dev/github.com/go-fuego/fuego#WithAddr).

```go
package main

import "github.com/go-fuego/fuego"

func main() {
	s := fuego.NewServer(
		fuego.WithAddr("localhost:8080"),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				DisableSwaggerUI: true,
			}),
		),
	)

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	s.Run()
}
```

### Address

You can change the address of the server with the `WithAddr` option.

```go
import "github.com/go-fuego/fuego"

func main() {
	s := fuego.NewServer(
		fuego.WithAddr("localhost:8080"),
	)
}
```

## Engine options

They are options at the Engine level, reusable for all routers (`net/http`, `gin`, `echo`).

```go
s := fuego.NewServer(
	fuego.WithEngineOptions(
		fuego.WithErrorHandler(func(err error) error {
			return fmt.Errorf("my wrapper: %w", err)
		}),
		fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
			UIHandler: func(specURL string) http.Handler {
				return dummyMiddleware(fuego.DefaultOpenAPIHandler(specURL))
			},
		}),
	),
)
```

## Route options

They are options at the route registration level. They allow you to declare query parameters, middlewares, description and more.

```go
package main

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

type MyInput struct {
	Name string `json:"name"`
}

type MyResponse struct {
	Name string `json:"name"`
}

func myController(c fuego.ContextWithBody[MyInput]) (*MyResponse, error) {
	name := c.QueryParam("name")
	return &MyResponse{
		Name: name,
	}, nil
}

var myReusableOption = option.Group(
	option.QueryInt("per_page", "Number of items per page", param.Default(100), param.Example("100 per page", 100)),
	option.QueryInt("page", "Page number", param.Default(1), param.Example("page 9", 9)),
)

func myCustomOption(r *fuego.BaseRoute) {
	r.XXX = "YYY" // Direct access to the route struct to inject custom behavior
}

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", myController,
		option.Query("name", "Name of the user", param.Required(), param.Example("example 1", "Napoleon")),

		option.Summary("Name getting route"),
		option.Description("This is the longdescription of the route"),
		option.Tags("Name", "Getting"),
		myCustomOption,
		myReusableOption,
	)

	s.Run()
}
```

### Set route options at Group level

You can pass route options on a routes Group and they will be inherited by all routes in the group.

I personally recommend **using the `option.Group` instead** of this to adopt a more composable approach instead of using inheritance here.

```go
package main

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func main() {
	s := fuego.NewServer()

	g := fuego.Group(s, "/pets",
		option.Summary("Pets operations"),
		option.Description("Operations about pets"),
		option.Tags("pets"),
	)

	fuego.Get(g, "/", func(c fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})
}
```

### Set route options at Server level

You can pass route options on the server and they will be inherited by all routes.

I personally recommend **using the `option.Group` instead** of this to adopt a more composable approach instead of using inheritance here.

```go
package main

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func main() {
	s := fuego.NewServer(
		fuego.WithRouteOptions(
			option.Summary("Pets operations"),
			option.Description("Operations about pets"),
			option.Tags("pets"),
		),
	)

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})
}
```
