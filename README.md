<p align="center">
  <img src="./data/fuego.svg" height="200" alt="Fuego Logo" />
</p>

# Fuego 🔥

[![Go Reference](https://pkg.go.dev/badge/github.com/go-fuego/fuego.svg)](https://pkg.go.dev/github.com/go-fuego/fuego)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-fuego/fuego)](https://goreportcard.com/report/github.com/go-fuego/fuego)
[![Coverage Status](https://coveralls.io/repos/github/go-fuego/fuego/badge.svg?branch=main)](https://coveralls.io/github/go-fuego/fuego?branch=main)

> The framework for busy Go developers

The only Go framework generating OpenAPI documentation from code. Inspired by Nest, built for Go developers.

Also empowers `html/template`, `a-h/templ` and `maragudk/gomponents`: see [the example](./examples/simple-crud) - actually running [in prod](https://gourmet.quimerch.com)!

## Why Fuego?

Chi, Gin, Fiber and Echo are great frameworks. But since they were designed a long time ago, they do not enjoy the possibilities that modern Go provides. Fuego offers a lot of features that make it easy to develop APIs and web applications.

## Features

- **OpenAPI**: Fuego automatically generates OpenAPI documentation from **code** - _not from comments nor yaml files!_
- **`net/http` compatible**: Fuego is built on top of `net/http`, so you can use any `http.Handler` middleware or handler! Fuego also supports `log/slog`, `context` and `html/template`
- **Routing**: Fuego provides a simple API for the Go 1.22 `net/http`
- **Serialization/Deserialization**: Fuego automatically serializes and deserializes JSON, XML and HTML Forms based on user-provided structs (or not, if you want to do it yourself)
- **Validation**: Fuego provides a simple and fast validator based on `go-playground/validator`
- **Transformation**: easily transform your data by implementing the `fuego.InTransform` and `fuego.OutTransform` interfaces - also useful for **validation**
- **Middlewares**: easily add a custom `net/http` middleware or use the provided middlewares.
- **Error handling**: Fuego provides centralized error handling with the standard [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) incoming.
- **Rendering**: Fuego provides a simple and fast rendering system based on `html/template` - you can still also use your own template system like `templ` or `gomponents`

## Examples

```go
package main

import "github.com/go-fuego/fuego"

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	s.Run()
}
```

<details>
<summary>All features</summary>

```go
package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-fuego/fuego"
	"github.com/rs/cors"
)

type Received struct {
	Name string `json:"name" validate:"required"`
}

type MyResponse struct {
	Message       string `json:"message"`
	BestFramework string `json:"best"`
}

func main() {
	s := fuego.NewServer(
		fuego.WithPort(":8088"),
	)

	fuego.Use(s, cors.Default().Handler)
	fuego.Use(s, chiMiddleware.Compress(5, "text/html", "text/css"))

	// Fuego 🔥 handler with automatic OpenAPI generation, validation, (de)serialization and error handling
	fuego.Post(s, "/", func(c *fuego.ContextWithBody[Received]) (MyResponse, error) {
		data, err := c.Body()
		if err != nil {
			return MyResponse{}, err
		}

		c.Response().Header().Set("X-Hello", "World")

		return MyResponse{
			Message:       "Hello, " + data.Name,
			BestFramework: "Fuego!",
		}, nil
	})

	// Standard net/http handler with automatic OpenAPI route declaration
	fuego.GetStd(s, "/std", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	s.Run()
}

// InTransform will be called when using c.Body().
// It can be used to transform the entity and raise custom errors
func (r *Received) InTransform(context.Context) error {
	r.Name = strings.ToLower(r.Name)
	if r.Name == "fuego" {
		return errors.New("fuego is not a name")
	}
	return nil
}

// OutTransform will be called before sending data
func (r *MyResponse) OutTransform(context.Context) error {
	r.Message = strings.ToUpper(r.Message)
	return nil
}
```

```bash
curl  http://localhost:8088/std
# Hello, World!
curl http://localhost:8088 -X POST -d '{"name": "Your Name"}' -H 'Content-Type: application/json'
# {"message":"HELLO, YOUR NAME","best":"Fuego!"}
curl http://localhost:8088 -X POST -d '{"name": "Fuego"}' -H 'Content-Type: application/json'
# {"error":"cannot transform request body: cannot transform request body: fuego is not a name"}
```

</details>

## Contributing

See the [contributing guide](CONTRIBUTING.md)

## Disclaimer for experienced gophers

I know you might prefer to use `net/http` directly, but if having a frame can convince my company to use Go instead of Node, I'm happy to use it.

## License

[GPL](./LICENSE.txt)
