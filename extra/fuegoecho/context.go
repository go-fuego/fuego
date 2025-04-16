package fuegoecho

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
)

type echoContext[B, P any] struct {
	internal.CommonContext[B]
	echoCtx echo.Context
}

var (
	_ fuego.ContextWithBody[any]      = &echoContext[any, any]{}
	_ fuego.ContextFlowable[any, any] = &echoContext[any, any]{}
)

func (c echoContext[B, P]) Body() (B, error) {
	var body B
	err := c.echoCtx.Bind(&body)
	if err != nil {
		return body, err
	}

	return fuego.TransformAndValidate(c, body)
}

func (c echoContext[B, P]) Context() context.Context {
	return c.echoCtx.Request().Context()
}

func (c echoContext[B, P]) Cookie(name string) (*http.Cookie, error) {
	return c.echoCtx.Request().Cookie(name)
}

func (c echoContext[B, P]) Header(key string) string {
	return c.echoCtx.Request().Header.Get(key)
}

func (c echoContext[B, P]) MustBody() B {
	body, err := c.Body()
	if err != nil {
		panic(err)
	}
	return body
}

func (c echoContext[B, P]) Params() (P, error) {
	var params P
	err := c.echoCtx.Bind(&params)
	if err != nil {
		return params, err
	}
	return params, nil
}

func (c echoContext[B, P]) MustParams() P {
	params, err := c.Params()
	if err != nil {
		panic(err)
	}
	return params
}

func (c echoContext[B, P]) PathParam(name string) string {
	return c.echoCtx.Param(name)
}

func (c echoContext[B, P]) PathParamIntErr(name string) (int, error) {
	return fuego.PathParamIntErr(c, name)
}

func (c echoContext[B, P]) PathParamInt(name string) int {
	param, _ := fuego.PathParamIntErr(c, name)
	return param
}

func (c echoContext[B, P]) MainLang() string {
	return strings.Split(c.MainLocale(), "-")[0]
}

func (c echoContext[B, P]) MainLocale() string {
	return strings.Split(c.Request().Header.Get("Accept-Language"), ",")[0]
}

func (c echoContext[B, P]) Redirect(code int, url string) (any, error) {
	c.echoCtx.Redirect(code, url)
	return nil, nil
}

func (c echoContext[B, P]) Render(templateToExecute string, data any, templateGlobsToOverride ...string) (fuego.CtxRenderer, error) {
	panic("unimplemented")
}

func (c echoContext[B, P]) Request() *http.Request {
	return c.echoCtx.Request()
}

func (c echoContext[B, P]) Response() http.ResponseWriter {
	return c.echoCtx.Response()
}

func (c echoContext[B, P]) SetCookie(cookie http.Cookie) {
	c.echoCtx.SetCookie(&cookie)
}

func (c echoContext[B, P]) HasCookie(name string) bool {
	_, err := c.Cookie(name)
	return err == nil
}

func (c echoContext[B, P]) HasHeader(key string) bool {
	_, ok := c.echoCtx.Request().Header[key]
	return ok
}

func (c echoContext[B, P]) SetHeader(key, value string) {
	c.echoCtx.Response().Header().Add(key, value)
}

func (c echoContext[B, P]) SetStatus(code int) {
	c.echoCtx.Response().WriteHeader(code)
}

func (c echoContext[B, P]) Serialize(data any) error {
	status := c.echoCtx.Response().Status
	if status == 0 {
		status = http.StatusOK
	}
	c.echoCtx.JSON(status, data)
	return nil
}

func (c echoContext[B, P]) SerializeError(err error) {
	statusCode := http.StatusInternalServerError
	var errorWithStatusCode fuego.ErrorWithStatus
	if errors.As(err, &errorWithStatusCode) {
		statusCode = errorWithStatusCode.StatusCode()
	}
	c.echoCtx.JSON(statusCode, err)
}

func (c echoContext[B, P]) SetDefaultStatusCode() {
	if c.DefaultStatusCode == 0 {
		c.DefaultStatusCode = http.StatusOK
	}
	c.echoCtx.Response().Status = c.DefaultStatusCode
}
