package op

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
)

const (
	maxBodySize = 1048576
)

// Ctx is the context of the request.
// It contains the request body, the path parameters, the query parameters, and the http request.
type Ctx[B any] interface {
	// Body returns the body of the request.
	// If (*B) implements [InTransformer], it will be transformed after deserialization.
	// It caches the result, so it can be called multiple times.
	Body() (B, error)

	// MustBody works like Body, but panics if there is an error.
	MustBody() B

	// PathParam returns the path parameter with the given name.
	// If it does not exist, it returns an empty string.
	// Example:
	//   op.Get(s, "/recipes/{recipe_id}", func(c op.Ctx[any]) (any, error) {
	//	 	id := c.PathParam("recipe_id")
	//   	...
	//   })
	PathParam(name string) string
	PathParams() map[string]string
	QueryParam(name string) string
	QueryParams() map[string]string

	Render(data any, templates ...string) (HTML, error)

	Request() *http.Request        // Request returns the underlying http request.
	Response() http.ResponseWriter // Response returns the underlying http response writer.

	// Context returns the context of the request.
	// Same as c.Request().Context().
	// This is the context related to the request, not the context of the server.
	Context() context.Context
}

func NewContext[B any](w http.ResponseWriter, r *http.Request, options readOptions) *Context[B] {
	c := &Context[B]{
		response: w,
		request:  r,
		readOptions: readOptions{
			DisallowUnknownFields: options.DisallowUnknownFields,
			MaxBodySize:           options.MaxBodySize,
		},
		fsInitMessage: "(filesystem given '.' from current working directory)",
		fs:            os.DirFS("."),
	}

	return c
}

// Context for the request. BodyType is the type of the request body. Please do not use a pointer type as parameter.
type Context[BodyType any] struct {
	body       *BodyType
	request    *http.Request
	response   http.ResponseWriter
	pathParams map[string]string

	fsInitMessage string // The input given to NewContext. Used for error messages.
	fs            fs.FS

	readOptions readOptions
}

// readOptions are options for reading the request body.
type readOptions struct {
	DisallowUnknownFields bool
	MaxBodySize           int64
	LogBody               bool
}

var _ Ctx[any] = &Context[any]{} // Check that Context implements Ctx.

// Context returns the context of the request.
// Same as c.Request().Context().
func (c Context[B]) Context() context.Context {
	return c.request.Context()
}

// Render renders the given templates with the given data.
// It returns just an empty string, because the response is written directly to the http.ResponseWriter.
func (c Context[B]) Render(data any, templates ...string) (HTML, error) {
	tmpl, err := template.ParseFS(c.fs, templates...)
	if err != nil {

		var pathError *fs.PathError
		if errors.As(err, &pathError) {
			wd, _ := os.Getwd()
			return "", fmt.Errorf("template '%s' does not exist in directory '%s': %w", pathError.Path, wd, err)
		}

		return "", fmt.Errorf("%w %s", err, c.fsInitMessage)
	}

	c.response.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(c.response, data)

	return "", err
}

// PathParams returns the path parameters of the request.
func (c Context[B]) PathParam(name string) string {
	param := c.pathParams[name]
	if param == "" {
		slog.Error("Path parameter might be invalid", "name", name, "valid parameters", c.pathParams)
	}
	return param // TODO go1.22: get (*http.Request) PathValue(name)
}

// PathParams returns the path parameters of the request.
func (c Context[B]) PathParams() map[string]string {
	return nil
}

// QueryParams returns the query parameters of the request.
func (c Context[B]) QueryParams() map[string]string {
	queryParams := c.request.URL.Query()
	params := make(map[string]string)
	for k, v := range queryParams {
		params[k] = v[0]
	}
	return params
}

// QueryParam returns the query parameter with the given name.
func (c Context[B]) QueryParam(name string) string {
	return c.request.URL.Query().Get(name)
}

// Request returns the http request.
func (c Context[B]) Request() *http.Request {
	return c.request
}

// Response returns the http response writer.
func (c Context[B]) Response() http.ResponseWriter {
	return c.response
}

// MustBody works like Body, but panics if there is an error.
func (c *Context[B]) MustBody() B {
	b, err := c.Body()
	if err != nil {
		panic(err)
	}
	return b
}

// Body returns the body of the request.
// If (*B) implements [InTransformer], it will be transformed after deserialization.
// It caches the result, so it can be called multiple times.
// The reason why the body is cached is because it is not possible to read an http request body multiple times, not because of performance.
// For decoding, it uses the Content-Type header. If it is not set, defaults to application/json.
func (c *Context[B]) Body() (B, error) {
	if c.body != nil {
		return *c.body, nil
	}

	// Limit the size of the request body.
	if c.readOptions.MaxBodySize != 0 {
		c.request.Body = http.MaxBytesReader(nil, c.request.Body, c.readOptions.MaxBodySize)
	}

	var body B
	var err error
	switch c.request.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded", "text/plain":
		s, errReadingString := readString[string](c.request.Body, c.readOptions)
		body = any(s).(B)
		err = errReadingString
	case "application/json":
		fallthrough
	default:
		body, err = readJSON[B](c.request.Body, c.readOptions)
	}

	c.body = &body

	return body, err
}
