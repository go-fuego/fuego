package fuegogin

import (
	"context"
	"net/http"
	"strings"

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
	err := c.ginCtx.BindJSON(&body)
	return body, err
}

func (c ginContext[B]) Context() context.Context {
	return c.ginCtx
}

func (c ginContext[B]) Cookie(name string) (*http.Cookie, error) {
	return c.ginCtx.Request.Cookie(name)
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

func (c ginContext[B]) MainLang() string {
	return strings.Split(c.MainLocale(), "-")[0]
}

func (c ginContext[B]) MainLocale() string {
	return strings.Split(c.Request().Header.Get("Accept-Language"), ",")[0]
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
	c.ginCtx.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
}

func (c ginContext[B]) SetHeader(key, value string) {
	c.ginCtx.Header(key, value)
}

func (c ginContext[B]) SetStatus(code int) {
	c.ginCtx.Status(code)
}
