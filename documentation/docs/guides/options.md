# Options

You can customize the server with the following function options.

All the options start with `With` and are located in the `fuego` package.

```go
package main

import "github.com/go-fuego/fuego"

func main() {
	s := fuego.NewServer(
		fuego.WithPort(":8080"),
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
	fuego.WithPort(":8080"),
)
```

### CORS

CORS middleware is not registered as a usual middleware, because it applies on routes that aren't registered. For example, OPTIONS /foo is not a registered route, but it's a request that needs to be handled by the CORS middleware.

```go
s := fuego.NewServer(
	fuego.WithCorsMiddleware(rs{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}),
)
```
