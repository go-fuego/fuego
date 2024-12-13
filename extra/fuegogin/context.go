package fuegogin

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/go-fuego/fuego"
)

type ContextWithBody[B any] struct {
	ginCtx *gin.Context
}

type ContextNoBody = ContextWithBody[any]

// Body implements fuego.Ctx.
func (c *ContextWithBody[B]) Body() (B, error) {
	var body B
	err := c.ginCtx.Bind(&body)
	return body, err
}

// Context implements fuego.Ctx.
func (c *ContextWithBody[B]) Context() context.Context {
	return c.ginCtx
}

// Cookie implements fuego.Ctx.
func (c *ContextWithBody[B]) Cookie(name string) (*http.Cookie, error) {
	panic("unimplemented")
}

// Header implements fuego.Ctx.
func (c *ContextWithBody[B]) Header(key string) string {
	return c.ginCtx.GetHeader(key)
}

// MustBody implements fuego.Ctx.
func (c *ContextWithBody[B]) MustBody() B {
	body, err := c.Body()
	if err != nil {
		panic(err)
	}
	return body
}

// PathParam implements fuego.Ctx.
func (c *ContextWithBody[B]) PathParam(name string) string {
	return c.ginCtx.Param(name)
}

// QueryParam implements fuego.Ctx.
func (c *ContextWithBody[B]) QueryParam(name string) string {
	return c.ginCtx.Query(name)
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
	return c.ginCtx.Request.URL.Query()
}

// Redirect implements fuego.Ctx.
func (c *ContextWithBody[B]) Redirect(code int, url string) (any, error) {
	c.ginCtx.Redirect(code, url)
	return nil, nil
}

// Render implements fuego.Ctx.
func (c *ContextWithBody[B]) Render(templateToExecute string, data any, templateGlobsToOverride ...string) (fuego.CtxRenderer, error) {
	panic("unimplemented")
}

// Request implements fuego.Ctx.
func (c *ContextWithBody[B]) Request() *http.Request {
	return c.ginCtx.Request
}

// Response implements fuego.Ctx.
func (c *ContextWithBody[B]) Response() http.ResponseWriter {
	return c.ginCtx.Writer
}

// SetCookie implements fuego.Ctx.
func (c *ContextWithBody[B]) SetCookie(cookie http.Cookie) {
}

// SetHeader implements fuego.Ctx.
func (c *ContextWithBody[B]) SetHeader(key, value string) {
	c.ginCtx.Header(key, value)
}

// SetStatus implements fuego.Ctx.
func (c *ContextWithBody[B]) SetStatus(code int) {
	c.ginCtx.Status(code)
}
