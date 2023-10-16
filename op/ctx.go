package op

import (
	"encoding/json"
	"net/http"

	"log/slog"
)

const (
	maxBodySize = 1048576
)

// Context for the request. BodyType is the type of the request body.
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

// Body returns the body of the request. It caches the result, so it can be called multiple times.
func (c *Ctx[B]) Body() (B, error) {
	if c.body != nil {
		return *c.body, nil
	}

	c.request.Body = http.MaxBytesReader(nil, c.request.Body, maxBodySize)

	dec := json.NewDecoder(c.request.Body)
	if config.DisallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(&c.body)
	if err != nil {
		slog.Info("Error decoding body", "err", err)
		return *c.body, err
	}
	slog.Info("Decoded body", "body", *c.body)

	return *c.body, nil
}

type RteCtx[T any] struct {
	ReturnType T
	BodyType   any
	ErrorType  error
}
