package internal

import (
	"context"
	"time"
)

// Base context shared by all adaptors (net/http, gin, echo, etc...)
type CommonContext[B any] struct {
	CommonCtx context.Context
}

func (c CommonContext[B]) Context() context.Context {
	return c.CommonCtx
}

// ContextNoBody implements the context interface via [net/http.Request.Context]
func (c CommonContext[B]) Deadline() (deadline time.Time, ok bool) {
	return c.Context().Deadline()
}

// ContextNoBody implements the context interface via [net/http.Request.Context]
func (c CommonContext[B]) Done() <-chan struct{} {
	return c.Context().Done()
}

// ContextNoBody implements the context interface via [net/http.Request.Context]
func (c CommonContext[B]) Err() error {
	return c.Context().Err()
}

// ContextNoBody implements the context interface via [net/http.Request.Context]
func (c CommonContext[B]) Value(key any) any {
	return c.Context().Value(key)
}
