package op

import (
	"encoding/json"
	"fmt"
	"net/http"

	"log/slog"
)

const (
	maxBodySize = 1048576
)

// Context for the request. BodyType is the type of the request body. Please do not use a pointer type as parameter.
type Ctx[BodyType any] struct {
	body    *BodyType
	request *http.Request
}

// PathParams returns the path parameters of the request.
func (c Ctx[B]) PathParams() map[string]string {
	return nil
}

// QueryParams returns the query parameters of the request.
func (c Ctx[B]) QueryParams() map[string]string {
	queryParams := c.request.URL.Query()
	params := make(map[string]string)
	for k, v := range queryParams {
		params[k] = v[0]
	}
	return params
}

// QueryParam returns the query parameter with the given name.
func (c Ctx[B]) QueryParam(name string) string {
	return c.request.URL.Query().Get(name)
}

// Request returns the http request.
func (c Ctx[B]) Request() *http.Request {
	return c.request
}

// MustBody works like Body, but panics if there is an error.
func (c *Ctx[B]) MustBody() B {
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
func (c *Ctx[B]) Body() (B, error) {
	if c.body != nil {
		return *c.body, nil
	}

	// Limit the size of the request body.
	c.request.Body = http.MaxBytesReader(nil, c.request.Body, maxBodySize)

	// Deserialize the request body.
	dec := json.NewDecoder(c.request.Body)
	if config.DisallowUnknownFields {
		dec.DisallowUnknownFields()
	}
	err := dec.Decode(&c.body)
	if err != nil {
		return *c.body, fmt.Errorf("cannot decode request body: %w", err)
	}
	slog.Info("Decoded body", "body", *c.body)

	// Validation
	// err = validation.Validate(rs.Validate, t)
	// if err != nil {
	// 	errWrapped := fmt.Errorf("cannot validate request body: %w", err)
	// 	common.SendError(w, exceptions.BadRequest{Err: err, Message: err.Error()})
	// 	return t, errWrapped
	// }

	// Normalize input if possible.
	if normalizableBody, ok := any(c.body).(Normalizable); ok {
		err := normalizableBody.Normalize()
		if err != nil {
			return *c.body, fmt.Errorf("error normalizing request body: %w", err)
		}
		c.body, ok = any(normalizableBody).(*B)
		if !ok {
			return *c.body, fmt.Errorf("error retyping request body: %w",
				fmt.Errorf("normalized body is not of type %T but should be", *new(B)))
		}

		slog.Info("Normalized body", "body", *c.body)
	}

	return *c.body, nil
}

type RteCtx[T any] struct {
	ReturnType T
	BodyType   any
	ErrorType  error
}
