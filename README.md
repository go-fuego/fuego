<p align="center">
  <img src="./data/fuego.svg" height="200" alt="Fuego Logo" />
</p>

# Fuego 🔥

[![Go Reference](https://pkg.go.dev/badge/github.com/go-op/op.svg)](https://pkg.go.dev/github.com/go-op/op)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-op/op)](https://goreportcard.com/report/github.com/go-op/op)

> The Go framework for busy API developers

The only Go framework generating OpenAPI documentation from code. Inspired by Nest, built for Go developers.

## Why Fuego?

Chi, Gin, Fiber and Echo are great frameworks. But since they were designed a long time ago, they do not enjoy the possibilities that modern Go provides. Fuego offers a lot of helper functions and features that make it easy to develop APIs.

## Features

- **OpenAPI**: Fuego automatically generates OpenAPI documentation from code
- **`net/http` compatible**: Fuego is built on top of `net/http`, so you can use any `net/http` middleware or handler!
- **Routing**: Fuego provides a simple and fast router based on Go 1.22 `net/http`
- **Serialization/Deserialization**: Fuego automatically serializes and deserializes JSON and XML based on user-provided structs (or not, if you want to do it yourself)
- **Validation**: Fuego provides a simple and fast validator based on go-playground/validator
- **Transformation**: easily transform your data after deserialization
- **Middlewares**: easily add a custom `net/http` middleware or use the built-in middlewares.

## Examples

```go
package main

import (
	"net/http"

	"github.com/op-go/op"
	"github.com/rs/cors"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type Received struct {
	Name string `json:"name" validate:"required"`
}

type MyResponse struct {
	Message       string `json:"message"`
	BestFramework string `json:"best"`
}

func main() {
	s := op.New()

	op.UseStd(s, cors.Default().Handler)
	op.UseStd(s, chiMiddleware.Compress(5, "text/html", "text/css"))

	op.Post(s, "/", func(c op.Ctx[Received]) (MyResponse, error) {
		data, err := c.Body()
		if err != nil {
			return MyResponse{}, err
		}

		return MyResponse{
			Message:       "Hello, " + data.Name,
			BestFramework: "Fuego!",
		}, nil
	})

	op.GetStd(s, "/std", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	s.Run()
}
```

## Contributing

See the [contributing guide](CONTRIBUTING.md)

## Disclaimer for experienced gophers

I know you might prefer to use `net/http` directly, but if having a frame can convince my company to use Go instead of Node, I'm happy to use it.

## License

[GPL](./LICENSE.txt)
