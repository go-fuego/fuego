// Package fuegomux provides a gorilla/mux adapter for the Fuego web framework.
package fuegomux

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
)

type muxContext[B, P any] struct {
	internal.CommonContext[B]
	req *http.Request
	res http.ResponseWriter
}

var (
	_ fuego.Context[any, any]         = &muxContext[any, any]{}
	_ fuego.ContextWithBody[any]      = &muxContext[any, any]{}
	_ fuego.ContextFlowable[any, any] = &muxContext[any, any]{}
)

func (c *muxContext[B, P]) Body() (B, error) {
	var body B
	err := json.NewDecoder(c.req.Body).Decode(&body)
	if err != nil {
		return body, err
	}
	return fuego.TransformAndValidate(c, body)
}

func (c *muxContext[B, P]) MustBody() B {
	body, err := c.Body()
	if err != nil {
		panic(err)
	}
	return body
}

func (c *muxContext[B, P]) Params() (P, error) {
	var params P
	return params, nil
}

func (c *muxContext[B, P]) MustParams() P {
	params, err := c.Params()
	if err != nil {
		panic(err)
	}
	return params
}

func (c *muxContext[B, P]) Context() context.Context {
	return c.req.Context()
}

func (c *muxContext[B, P]) Cookie(name string) (*http.Cookie, error) {
	return c.req.Cookie(name)
}

func (c *muxContext[B, P]) HasCookie(name string) bool {
	_, err := c.Cookie(name)
	return err == nil
}

func (c *muxContext[B, P]) Header(key string) string {
	return c.req.Header.Get(key)
}

func (c *muxContext[B, P]) HasHeader(key string) bool {
	_, ok := c.req.Header[key]
	return ok
}

func (c *muxContext[B, P]) SetHeader(key, value string) {
	c.res.Header().Set(key, value)
}

func (c *muxContext[B, P]) SetCookie(cookie http.Cookie) {
	http.SetCookie(c.res, &cookie)
}

func (c *muxContext[B, P]) PathParam(name string) string {
	return mux.Vars(c.req)[name]
}

func (c *muxContext[B, P]) PathParamIntErr(name string) (int, error) {
	return fuego.PathParamIntErr(c, name)
}

func (c *muxContext[B, P]) PathParamInt(name string) int {
	param, _ := fuego.PathParamIntErr(c, name)
	return param
}

func (c *muxContext[B, P]) MainLang() string {
	return strings.Split(c.MainLocale(), "-")[0]
}

func (c *muxContext[B, P]) MainLocale() string {
	return strings.Split(c.req.Header.Get("Accept-Language"), ",")[0]
}

func (c *muxContext[B, P]) Redirect(code int, url string) (any, error) {
	http.Redirect(c.res, c.req, url, code)
	return nil, nil
}

func (c *muxContext[B, P]) Render(templateToExecute string, data any, templateGlobsToOverride ...string) (fuego.CtxRenderer, error) {
	panic("unimplemented")
}

func (c *muxContext[B, P]) Request() *http.Request {
	return c.req
}

func (c *muxContext[B, P]) Response() http.ResponseWriter {
	return c.res
}

func (c *muxContext[B, P]) SetStatus(code int) {
	c.res.WriteHeader(code)
}

func (c *muxContext[B, P]) Serialize(data any) error {
	return fuego.Send(c.res, c.req, data)
}

func (c *muxContext[B, P]) SerializeError(err error) {
	fuego.SendError(c.res, c.req, err)
}

func (c *muxContext[B, P]) SetDefaultStatusCode() {
	if c.DefaultStatusCode != 0 {
		c.SetStatus(c.DefaultStatusCode)
	}
}
