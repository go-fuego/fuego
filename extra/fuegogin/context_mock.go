package fuegogin

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type ContextTest[B any] struct {
	ginContext[B]
	BodyInjected  B
	ErrorInjected error

	Params url.Values
}

func (c *ContextTest[B]) Body() (B, error) {
	return c.BodyInjected, c.ErrorInjected
}

func (c *ContextTest[B]) Request() *http.Request {
	return c.ginCtx.Request
}

func (c *ContextTest[B]) Response() gin.ResponseWriter {
	return c.ginCtx.Writer
}

// QueryParam implements fuego.Ctx
func (c *ContextTest[B]) QueryParam(key string) string {
	return c.Params.Get(key)
}
