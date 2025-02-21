# Controllers & Deserialization

Controllers are the main way to interact with the application. They are responsible for handling the requests and responses.

## Controller types

The standard controller type is a function that receives a `fuego.ContextWithBody` and returns a response and an error.

```go
// Standard fuego controller
func MyController(c fuego.ContextWithBody[Body]) (MyResponse, error)
```

Fuego handles [Content Negotiation](https://developer.mozilla.org/en-US/docs/Web/HTTP/Content_negotiation), so you can serve different content types based on the `Accept` header of the request. One controller, multiple responses formats.

> Just care about data types! Fuego will handle the rest.

Please note that Fuego relies on the underlying libraries (`net/http`, `encoding/json`, `gin` if using `fuegogin`) to handle deserialization. We avoid at all cost the use of reflection to keep the performance high.

### Registering a controller

See [Routing](./routing.md) for more information.

```go
fuego.Get(s, "/", MyController)
```

## Request body

Use the `fuego.ContextWithBody` interface. Useful for `POST`, `PUT`, `PATCH` requests for example.

```go
type MyInput struct {
    Name string `json:"name"`
}

func MyController(c fuego.ContextWithBody[MyInput]) (*MyResponse, error) {
    body, err := c.Body()
    if err != nil {
        return nil, err
    }

    return &MyResponse{
        Name: body.Name,
    }, nil
}
```

Fuego will automatically parse the request body (`JSON`, `XML`, `YAML`, `application/x-www-form-urlencoded` and `multipart/form-data`) [according to the `Content-Type` header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type) of the request.

```curl
curl -X POST http://localhost:9999/ -d '{"name": "My name"}' -H "Content-Type: application/json"
# Response: {"name": "My name"}

curl -X POST http://localhost:9999/ -d '<MyInput><Name>My name</Name></MyInput>' -H "Content-Type: application/xml"
# Response: {"name": "My name"}
```

It will then validate it using the input struct, see [Validation](./validation.md).

### I don't need request body

Use the `fuego.ContextNoBody` interface. Useful for `GET`, `DELETE`, `HEAD`, `OPTIONS` requests for example.

```go
func MyController(c fuego.ContextNoBody) (MyResponse, error) {
    return MyResponse{Name: "My name"}, nil
}
```

### Binary body

If you just want to read the body of the request as a byte slice, you can use the `[]byte` receiver type.

Don't forget to set the request `Content-Type` header to `application/octet-stream`.

```go
fuego.Put(s, "/blob", func(c fuego.ContextWithBody[[]byte]) (any, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}

	return body, nil
})
```

## Query parameters (dynamic)

They are declared (for OpenAPI and validation) at the route registration level. It is not type-safe (it relies on the same string on the route registration and the controller) BUT it raises warning if you make a typo and use a non-declared query parameter.

See [Route Options](./options.md) for more information.

```go
package main

import (
    "github.com/go-fuego/fuego"
)

type MyInput struct {
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

func main() {
    s := fuego.NewServer()

    fuego.Get(s, "/", myController,
        option.Query("name", "Name of the user", param.Required(), param.Example("example 1", "Napoleon")),
        myReusableOption,
    )

    s.Run()
}
```

```curl
curl -X GET http://localhost:9999/?name=MyName
# Response: {"name": "MyName"}
```

## Query parameters (type-safe)

> 🚧 Below contains incoming syntax, not available currently

This syntax allows users to have strong static typing between the input and the query parameters.

```go
type Params struct {
    Limit int    `query:"limit"`
    Group string `header:"X-User-Group"`
}
func myController(c fuego.Context[MyInput, Params]) (*MyResponse, error) {
    params, err := c.Params()
    if err != nil {
        return nil, err
    }

    return &MyResponse{
        Name: params.Group,
    }, nil
}
```

## Headers

You can always go further in the request and response by using the underlying net/http request and response, by using `c.Request` and `c.Response`.

### Get request header

```go
func MyController(c fuego.ContextNoBody) (MyResponse, error) {
	value := c.Header("X-My-Header")
	return MyResponse{}, nil
}
```

### Set response header

```go
func MyController(c fuego.ContextNoBody) (MyResponse, error) {
	c.SetHeader("X-My-Header", "value")
	return MyResponse{}, nil
}
```

## Cookies

### Get request cookie

```go
func MyController(c fuego.ContextNoBody) (MyResponse, error) {
	value := c.Cookie("my-cookie")
	return MyResponse{}, nil
}
```

### Set response cookie

```go
func MyController(c fuego.ContextNoBody) (MyResponse, error) {
	c.SetCookie("my-cookie", "value")
	return MyResponse{}, nil
}
```
