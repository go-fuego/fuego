package op

import (
	"encoding/json"
	"net/http"

	"log/slog"
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
	return nil
}

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

	var b B
	err := json.NewDecoder(c.request.Body).Decode(&c.body)
	if err != nil {
		slog.Info("Error decoding body", "err", err)
		return b, err
	}
	slog.Info("Decoded body", "body", b)

	return *c.body, nil
}

type RteCtx[T any] struct {
	ReturnType T
	BodyType   any
	ErrorType  error
}
