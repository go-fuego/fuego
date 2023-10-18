package op

import (
	"net/http"
)

const (
	maxBodySize = 1048576
)

// Ctx is the context of the request.
// It contains the request body, the path parameters, the query parameters, and the http request.
type Ctx[B any] interface {
	// Body returns the body of the request.
	// If (*B) implements [Normalizable], it will be normalized.
	// It caches the result, so it can be called multiple times.
	Body() (B, error)

	// MustBody works like Body, but panics if there is an error.
	MustBody() B
	PathParams() map[string]string
	QueryParam(name string) string
	QueryParams() map[string]string
	Request() *http.Request
}

// Context for the request. BodyType is the type of the request body. Please do not use a pointer type as parameter.
type Context[BodyType any] struct {
	body    *BodyType
	request *http.Request

	readOptions readOptions
}

// readOptions are options for reading the request body.
type readOptions struct {
	DisallowUnknownFields bool
	MaxBodySize           int64
	LogBody               bool
}

var _ Ctx[any] = &Context[any]{} // Check that Context implements Ctx.

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

// MustBody works like Body, but panics if there is an error.
func (c *Context[B]) MustBody() B {
	b, err := c.Body()
	if err != nil {
		panic(err)
	}
	return b
}

// Body returns the body of the request.
// If (*B) implements [Normalizable], it will be normalized.
// It caches the result, so it can be called multiple times.
func (c *Context[B]) Body() (B, error) {
	if c.body != nil {
		return *c.body, nil
	}

	// Limit the size of the request body.
	if c.readOptions.MaxBodySize != 0 {
		c.request.Body = http.MaxBytesReader(nil, c.request.Body, c.readOptions.MaxBodySize)
	}

	switch any(new(B)).(type) {
	case *string:
		s, err := readString[string](c.request.Body, c.readOptions)
		body := any(s).(B)
		return body, err
	default:
		return readJSON[B](c.request.Body, c.readOptions)
	}
}
