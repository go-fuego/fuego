package op

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"log/slog"
)

const (
	maxBodySize = 1048576
)

type Ctx[B any] interface {
	Body() (B, error)
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

	config Config
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

// Normalizable is an interface for entities that can be normalized.
// Useful for example for trimming strings, add custom fields, etc.
// Can also raise an error if the entity is not valid.
type Normalizable interface {
	Normalize() error // Normalizes the entity.
}

// Body returns the body of the request.
// If (*B) implements [Normalizable], it will be normalized.
// It caches the result, so it can be called multiple times.
func (c *Context[B]) Body() (B, error) {
	if c.body != nil {
		return *c.body, nil
	}

	// Limit the size of the request body.
	c.request.Body = http.MaxBytesReader(nil, c.request.Body, maxBodySize)

	switch any(c.body).(type) {
	case *string:
		// Read the request body.
		body, err := io.ReadAll(c.request.Body)
		if err != nil {
			return *new(B), fmt.Errorf("cannot read request body: %w", err)
		}
		// c.body = (*B)(unsafe.Pointer(&body))
		s := string(body)
		c.body = any(&s).(*B)
		slog.Info("Read body", "body", *c.body)
	default:
		// Deserialize the request body.
		dec := json.NewDecoder(c.request.Body)
		if c.config.DisallowUnknownFields {
			dec.DisallowUnknownFields()
		}
		err := dec.Decode(&c.body)
		if err != nil {
			return *new(B), fmt.Errorf("cannot decode request body: %w", err)
		}
		slog.Info("Decoded body", "body", *c.body)

		// Validation
		err = validate(*c.body)
		if err != nil {
			return *c.body, fmt.Errorf("cannot validate request body: %w", err)
		}
	}

	// Normalize input if possible.
	if normalizableBody, ok := any(c.body).(Normalizable); ok {
		err := normalizableBody.Normalize()
		if err != nil {
			return *c.body, fmt.Errorf("cannot normalize request body: %w", err)
		}
		c.body, ok = any(normalizableBody).(*B)
		if !ok {
			return *c.body, fmt.Errorf("cannot retype request body: %w",
				fmt.Errorf("normalized body is not of type %T but should be", *new(B)))
		}

		slog.Info("Normalized body", "body", *c.body)
	}

	return *c.body, nil
}
