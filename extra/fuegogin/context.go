package fuegogin

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
)

type ContextWithBody[B any] interface {
	internal.CommonCtx[B]

	Request() *http.Request
	Response() gin.ResponseWriter

	// Original Gin context
	Context() *gin.Context
}

var _ internal.CommonCtx[string] = (ContextWithBody[string])(nil) // Check that ContextWithBody[string] implements CommonCtx.

type ContextNoBody = ContextWithBody[any]

type contextWithBody[B any] struct {
	ginCtx *gin.Context
}

var _ ContextWithBody[any] = &contextWithBody[any]{}

func (c *contextWithBody[B]) Body() (B, error) {
	var body B
	err := c.ginCtx.Bind(&body)
	return body, err
}

func (c *contextWithBody[B]) Context() *gin.Context {
	return c.ginCtx
}

func (c *contextWithBody[B]) Cookie(name string) (*http.Cookie, error) {
	panic("unimplemented")
}

func (c *contextWithBody[B]) Header(key string) string {
	return c.ginCtx.GetHeader(key)
}

func (c *contextWithBody[B]) MustBody() B {
	body, err := c.Body()
	if err != nil {
		panic(err)
	}
	return body
}

func (c *contextWithBody[B]) PathParam(name string) string {
	return c.ginCtx.Param(name)
}

func (c *contextWithBody[B]) QueryParam(name string) string {
	return c.ginCtx.Query(name)
}

func (c *contextWithBody[B]) QueryParamArr(name string) []string {
	panic("unimplemented")
}

func (c *contextWithBody[B]) QueryParamBool(name string) bool {
	panic("unimplemented")
}

func (c *contextWithBody[B]) QueryParamBoolErr(name string) (bool, error) {
	panic("unimplemented")
}

func (c *contextWithBody[B]) QueryParamInt(name string) int {
	panic("unimplemented")
}

func (c *contextWithBody[B]) QueryParamIntErr(name string) (int, error) {
	panic("unimplemented")
}

func (c *contextWithBody[B]) QueryParams() url.Values {
	return c.ginCtx.Request.URL.Query()
}

func (c *contextWithBody[B]) MainLang() string {
	panic("unimplemented")
}

func (c *contextWithBody[B]) MainLocale() string {
	panic("unimplemented")
}

func (c *contextWithBody[B]) Redirect(code int, url string) (any, error) {
	c.ginCtx.Redirect(code, url)
	return nil, nil
}

func (c *contextWithBody[B]) Render(templateToExecute string, data any, templateGlobsToOverride ...string) (fuego.CtxRenderer, error) {
	panic("unimplemented")
}

func (c *contextWithBody[B]) Request() *http.Request {
	return c.ginCtx.Request
}

func (c *contextWithBody[B]) Response() gin.ResponseWriter {
	return c.ginCtx.Writer
}

func (c *contextWithBody[B]) SetCookie(cookie http.Cookie) {
}

func (c *contextWithBody[B]) SetHeader(key, value string) {
	c.ginCtx.Header(key, value)
}

func (c *contextWithBody[B]) SetStatus(code int) {
	c.ginCtx.Status(code)
}
