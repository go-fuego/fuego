package fuegogin

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"

	"github.com/go-fuego/fuego"
)

func Get[T, B any](
	s *fuego.OpenAPI,
	e *gin.Engine,
	path string,
	handler func(c *ContextWithBody[B]) (T, error),
	options ...func(*fuego.BaseRoute),
) *fuego.Route[T, B] {
	return Handle(s, e, "GET", path, handler, options...)
}

func Post[T, B any](
	s *fuego.OpenAPI,
	e *gin.Engine,
	path string,
	handler func(c *ContextWithBody[B]) (T, error),
	options ...func(*fuego.BaseRoute),
) *fuego.Route[T, B] {
	return Handle(s, e, "POST", path, handler, options...)
}

func Handle[T, B any](
	openapi *fuego.OpenAPI,
	e *gin.Engine,
	method,
	path string,
	handler func(c *ContextWithBody[B]) (T, error),
	options ...func(*fuego.BaseRoute),
) *fuego.Route[T, B] {
	route := &fuego.Route[T, B]{
		BaseRoute: fuego.BaseRoute{
			Method:    method,
			Path:      path,
			Params:    make(map[string]fuego.OpenAPIParam),
			FullName:  fuego.FuncName(handler),
			Operation: openapi3.NewOperation(),
			OpenAPI:   openapi,
		},
	}

	for _, o := range options {
		o(&route.BaseRoute)
	}

	route.BaseRoute.GenerateDefaultDescription()

	e.Handle(method, path, func(c *gin.Context) {
		context := &ContextWithBody[B]{
			ginCtx: c,
		}

		resp, err := handler(context)
		if err != nil {
			c.Error(err)
			return
		}

		if c.Request.Header.Get("Accept") == "application/xml" {
			c.XML(200, resp)
			return
		}

		c.JSON(200, resp)
	})

	route.RegisterOpenAPIOperation(openapi)

	return route
}
