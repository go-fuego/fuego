package fuegogin

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
)

type ginContext[B any] struct {
	internal.CommonContext[B]
	ginCtx *gin.Context
}

var _ fuego.ContextWithBody[any] = &ginContext[any]{}

func (c ginContext[B]) Body() (B, error) {
	var body B
	err := c.ginCtx.Bind(&body)
	return body, err
}

func (c ginContext[B]) Context() context.Context {
	return c.ginCtx
}

func (c ginContext[B]) Cookie(name string) (*http.Cookie, error) {
	panic("unimplemented")
}

func (c ginContext[B]) Header(key string) string {
	return c.ginCtx.GetHeader(key)
}

func (c ginContext[B]) MustBody() B {
	body, err := c.Body()
	if err != nil {
		panic(err)
	}
	return body
}

func (c ginContext[B]) PathParam(name string) string {
	return c.ginCtx.Param(name)
}

func (c ginContext[B]) QueryParam(name string) string {
	return c.ginCtx.Query(name)
}

func (c ginContext[B]) QueryParamArr(name string) []string {
	panic("unimplemented")
}

func (c ginContext[B]) QueryParamBool(name string) bool {
	panic("unimplemented")
}

func (c ginContext[B]) QueryParamBoolErr(name string) (bool, error) {
	panic("unimplemented")
}

func (c ginContext[B]) QueryParamInt(name string) int {
	panic("unimplemented")
}

func (c ginContext[B]) QueryParamIntErr(name string) (int, error) {
	panic("unimplemented")
}

func (c ginContext[B]) QueryParams() url.Values {
	return c.ginCtx.Request.URL.Query()
}

func (c ginContext[B]) MainLang() string {
	panic("unimplemented")
}

func (c ginContext[B]) MainLocale() string {
	panic("unimplemented")
}

func (c ginContext[B]) Redirect(code int, url string) (any, error) {
	c.ginCtx.Redirect(code, url)
	return nil, nil
}

func (c ginContext[B]) Render(templateToExecute string, data any, templateGlobsToOverride ...string) (fuego.CtxRenderer, error) {
	panic("unimplemented")
}

func (c ginContext[B]) Request() *http.Request {
	return c.ginCtx.Request
}

func (c ginContext[B]) Response() http.ResponseWriter {
	return c.ginCtx.Writer
}

func (c ginContext[B]) SetCookie(cookie http.Cookie) {
}

func (c ginContext[B]) SetHeader(key, value string) {
	c.ginCtx.Header(key, value)
}

func (c ginContext[B]) SetStatus(code int) {
	c.ginCtx.Status(code)
}
