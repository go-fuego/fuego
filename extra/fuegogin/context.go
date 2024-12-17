package fuegogin

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/go-fuego/fuego"
)

type ContextWithBody[B any] interface {
	fuego.CommonCtx[B]

	Request() *http.Request
	Response() gin.ResponseWriter

	// Original Gin context
	Context() *gin.Context
}

type ContextNoBody = ContextWithBody[any]

type contextWithBody[B any] struct {
	ginCtx *gin.Context
}

// Body satisfies [fuego.Ctx].
func (c *contextWithBody[B]) Body() (B, error) {
	var body B
	err := c.ginCtx.Bind(&body)
	return body, err
}

// Context satisfies [fuego.Ctx].
func (c *contextWithBody[B]) Context() *gin.Context {
	return c.ginCtx
}

// Cookie satisfies [fuego.Ctx].
func (c *contextWithBody[B]) Cookie(name string) (*http.Cookie, error) {
	panic("unimplemented")
}

// Header satisfies [fuego.Ctx].
func (c *contextWithBody[B]) Header(key string) string {
	return c.ginCtx.GetHeader(key)
}

// MustBody satisfies [fuego.Ctx].
func (c *contextWithBody[B]) MustBody() B {
	body, err := c.Body()
	if err != nil {
		panic(err)
	}
	return body
}

// PathParam satisfies [fuego.Ctx].
func (c *contextWithBody[B]) PathParam(name string) string {
	return c.ginCtx.Param(name)
}

// QueryParam satisfies [fuego.Ctx].
func (c *contextWithBody[B]) QueryParam(name string) string {
	return c.ginCtx.Query(name)
}

// QueryParamArr satisfies [fuego.Ctx].
func (c *contextWithBody[B]) QueryParamArr(name string) []string {
	panic("unimplemented")
}

// QueryParamBool satisfies [fuego.Ctx].
func (c *contextWithBody[B]) QueryParamBool(name string) bool {
	panic("unimplemented")
}

// QueryParamBoolErr satisfies [fuego.Ctx].
func (c *contextWithBody[B]) QueryParamBoolErr(name string) (bool, error) {
	panic("unimplemented")
}

// QueryParamInt satisfies [fuego.Ctx].
func (c *contextWithBody[B]) QueryParamInt(name string) int {
	panic("unimplemented")
}

// QueryParamIntErr satisfies [fuego.Ctx].
func (c *contextWithBody[B]) QueryParamIntErr(name string) (int, error) {
	panic("unimplemented")
}

// QueryParams satisfies [fuego.Ctx].
func (c *contextWithBody[B]) QueryParams() url.Values {
	return c.ginCtx.Request.URL.Query()
}

// Redirect satisfies [fuego.Ctx].
func (c *contextWithBody[B]) Redirect(code int, url string) (any, error) {
	c.ginCtx.Redirect(code, url)
	return nil, nil
}

// Render satisfies [fuego.Ctx].
func (c *contextWithBody[B]) Render(templateToExecute string, data any, templateGlobsToOverride ...string) (fuego.CtxRenderer, error) {
	panic("unimplemented")
}

// Request satisfies [fuego.Ctx].
func (c *contextWithBody[B]) Request() *http.Request {
	return c.ginCtx.Request
}

// Response satisfies [fuego.Ctx].
func (c *contextWithBody[B]) Response() gin.ResponseWriter {
	return c.ginCtx.Writer
}

// SetCookie satisfies [fuego.Ctx].
func (c *contextWithBody[B]) SetCookie(cookie http.Cookie) {
}

// SetHeader satisfies [fuego.Ctx].
func (c *contextWithBody[B]) SetHeader(key, value string) {
	c.ginCtx.Header(key, value)
}

// SetStatus satisfies [fuego.Ctx].
func (c *contextWithBody[B]) SetStatus(code int) {
	c.ginCtx.Status(code)
}
