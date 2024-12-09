package fuegogin

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-fuego/fuego"
)

type ContextWithBody[B any] struct{}

// Body implements fuego.Ctx.
func (c *ContextWithBody[B]) Body() (B, error) {
	panic("unimplemented")
}

// Context implements fuego.Ctx.
func (c *ContextWithBody[B]) Context() context.Context {
	panic("unimplemented")
}

// Cookie implements fuego.Ctx.
func (c *ContextWithBody[B]) Cookie(name string) (*http.Cookie, error) {
	panic("unimplemented")
}

// Header implements fuego.Ctx.
func (c *ContextWithBody[B]) Header(key string) string {
	panic("unimplemented")
}

// MainLang implements fuego.Ctx.
func (c *ContextWithBody[B]) MainLang() string {
	panic("unimplemented")
}

// MainLocale implements fuego.Ctx.
func (c *ContextWithBody[B]) MainLocale() string {
	panic("unimplemented")
}

// MustBody implements fuego.Ctx.
func (c *ContextWithBody[B]) MustBody() B {
	panic("unimplemented")
}

// PathParam implements fuego.Ctx.
func (c *ContextWithBody[B]) PathParam(name string) string {
	panic("unimplemented")
}

// QueryParam implements fuego.Ctx.
func (c *ContextWithBody[B]) QueryParam(name string) string {
	panic("unimplemented")
}

// QueryParamArr implements fuego.Ctx.
func (c *ContextWithBody[B]) QueryParamArr(name string) []string {
	panic("unimplemented")
}

// QueryParamBool implements fuego.Ctx.
func (c *ContextWithBody[B]) QueryParamBool(name string) bool {
	panic("unimplemented")
}

// QueryParamBoolErr implements fuego.Ctx.
func (c *ContextWithBody[B]) QueryParamBoolErr(name string) (bool, error) {
	panic("unimplemented")
}

// QueryParamInt implements fuego.Ctx.
func (c *ContextWithBody[B]) QueryParamInt(name string) int {
	panic("unimplemented")
}

// QueryParamIntErr implements fuego.Ctx.
func (c *ContextWithBody[B]) QueryParamIntErr(name string) (int, error) {
	panic("unimplemented")
}

// QueryParams implements fuego.Ctx.
func (c *ContextWithBody[B]) QueryParams() url.Values {
	panic("unimplemented")
}

// Redirect implements fuego.Ctx.
func (c *ContextWithBody[B]) Redirect(code int, url string) (any, error) {
	panic("unimplemented")
}

// Render implements fuego.Ctx.
func (c *ContextWithBody[B]) Render(templateToExecute string, data any, templateGlobsToOverride ...string) (fuego.CtxRenderer, error) {
	panic("unimplemented")
}

// Request implements fuego.Ctx.
func (c *ContextWithBody[B]) Request() *http.Request {
	panic("unimplemented")
}

// Response implements fuego.Ctx.
func (c *ContextWithBody[B]) Response() http.ResponseWriter {
	panic("unimplemented")
}

// SetCookie implements fuego.Ctx.
func (c *ContextWithBody[B]) SetCookie(cookie http.Cookie) {
	panic("unimplemented")
}

// SetHeader implements fuego.Ctx.
func (c *ContextWithBody[B]) SetHeader(key, value string) {
	panic("unimplemented")
}

// SetStatus implements fuego.Ctx.
func (c *ContextWithBody[B]) SetStatus(code int) {
	panic("unimplemented")
}
