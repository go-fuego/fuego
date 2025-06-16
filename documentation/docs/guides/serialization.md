# Serialization

Serialization is the process of converting Go data structures into formats like JSON or XML for transmission, while deserialization converts them back. Fuego handles this automatically using standard Go packages.

## Content Negotiation with Accept Header

Fuego implements [HTTP content negotiation](https://developer.mozilla.org/en-US/docs/Web/HTTP/Content_negotiation) out of the box. Your API automatically responds with different formats based on the client's `Accept` header without any additional code.

When a client makes a request to your Fuego API, it can specify the desired response format using the `Accept` header. Fuego will automatically detect this header and serialize your response data accordingly.

For example:

- A browser might send `Accept: text/html` to get an HTML page
- A mobile app might send `Accept: application/json` to get JSON data
- An XML-based client might send `Accept: application/xml` to get XML data

If no `Accept` header is provided, Fuego defaults to JSON (`application/json`).

## Supported Formats

To serialize data, just return the data you want to serialize from your controller. It will be automatically serialized into one of the following formats, depending on the `Accept` header in the request:

- JSON: `Accept: application/json` (default)
- XML: `Accept: application/xml`
- YAML: `Accept: application/yaml`
- HTML: `Accept: text/html`
- Plain text: `Accept: text/plain`

```go
type MyReturnType struct {
	Message string `json:"message"`
}

func helloWorld(c fuego.ContextNoBody) (MyReturnType, error) {
	return MyReturnType{Message: "Hello, World!"}, nil
}

// curl request: curl -X GET http://localhost:8080/ -H "Accept: application/json"
// response: {"message":"Hello, World!"}

// curl request: curl -X GET http://localhost:8080/ -H "Accept: application/xml"
// response: <MyReturnType><Message>Hello, World!</Message></MyReturnType>
```

This means you can build a single API endpoint that serves both your web frontend (HTML) and your API clients (JSON/XML) without duplicating code.

## Custom response - Bypass return type

If you want to bypass the automatic serialization, you can directly write to the response writer.

```go
func helloWorld(c fuego.ContextNoBody) (any, error) {
	w := c.Response()
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello, World!"))                                      // Write directly to the response writer.
	_ = json.NewEncoder(w).Encode(MyReturnType{Message: "Hello, World!"}) // You can also use json.NewEncoder(w).Encode to serialize data directly into JSON
	fuego.SendJSON(w, MyReturnType{Message: "Hello, World!"})             // Or use fuego.SendJSON to serialize data directly into JSON

	return nil, nil // If you return nil, nil fuego will not try to serialize a response
}
```

## Custom serialization

But you can also use the `Serialize` and `Deserialize` functions to manually serialize and deserialize data.

See [the documentation](https://pkg.go.dev/github.com/go-fuego/fuego#Server) for the `fuego.Server` struct for more information.

```go
package main

import (
	"net/http"

	"github.com/go-fuego/fuego"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	s := fuego.NewServer()

	s.Serialize = func(w http.ResponseWriter, ans any) {
		w.Header().Set("Content-Type", "text/plain")
		json.NewEncoder(w).Encode(ans)
	}

	fuego.Get(s, "/", helloWorld)

	s.Run()
}

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}
```

### Custom Content Negotiation Serialization and Deserialization

You can define custom serializers and deserializers for specific content types using the `WithContentTypeSerDes` option. This allows you to handle custom data formats beyond the built-in JSON, XML, YAML, and other supported formats.

Just implement the `SerDes` interface:

```go
package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

// Custom key-value serializer that handles data like "key1=value1;key2=value2"
// This contrived example merely demonstrates how to perform custom de/serialiation
type kvSerDes struct {
	delimiter string
}

func (s kvSerDes) Serialize(v any) ([]byte, error) {
	data, ok := v.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("expected map[string]string, got %T", v)
	}
	
	var parts []string
	for key, value := range data {
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	return []byte(strings.Join(parts, s.delimiter)), nil
}

func (s kvSerDes) Deserialize(ctx context.Context, input io.Reader) (any, error) {
	body, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]string)
	pairs := strings.Split(string(body), s.delimiter)
	for _, pair := range pairs {
		if kv := strings.Split(pair, "="); len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result, nil
}

func myController(c fuego.ContextWithBody[map[string]string]) (map[string]string, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}
	body["processed"] = "true"
	return body, nil
}

func main() {
	s := fuego.NewServer()

	fuego.Post(s, "/supports-keyvalue", myController,
		option.WithContentTypeSerDes("application/vnd.keyvalue", kvSerDes{delimiter: ";"}),
	)

	s.Run()
}
```

With this setup:
- Requests with `Content-Type: application/vnd.keyvalue` will be deserialized using your custom deserializer
- Responses with `Accept: application/vnd.keyvalue` will be serialized using your custom serializer
- The route still supports standard Fuego content negotiation for JSON, YAML, XML, etc.

Custom SerDes take precedence over built-in SerDes for the specified content type, meaning you can override Fuego's default de/serialization logic for JSON, XML, YAML, etc. if you choose.

Note that your `SerDes`'s `Deserialize` implementation must return the type expected by your controller's body (`fuego.ContextWithBody[FooType]`) or an internal server error (HTTP status 500) will be thrown.

You can also use this option at the group or server level to easily apply custom content negotiation to multiple routes by using `WithRouteOptions`.

## Combining Data and HTML with DataOrHTML

For routes that need to serve both API clients and web browsers, Fuego provides a convenient `DataOrHTML` helper that returns different content based on the `Accept` header:

```go
package main

import (
	"github.com/go-fuego/fuego"
)

type UserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/user/profile", func(c fuego.ContextNoBody) (interface{}, error) {
		userData := UserData{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		return fuego.DataOrHTML(
			userData,                      // When Accept: application/json, return this data
			renderUserProfile(userData),   // When Accept: text/html, render this HTML
		), nil
	})

	s.Run()
}

func renderUserProfile(user UserData) string {
	return "<h1>User Profile</h1><p>Name: " + user.Name + "</p><p>Email: " + user.Email + "</p>"
}
```

This approach allows you to build APIs and web interfaces with the same codebase, reducing duplication and ensuring consistency between your API and web UI.
