package fuegogin

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/go-fuego/fuego"
)

func GetGin(s *fuego.OpenAPI, e gin.IRouter, path string, handler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	return handleGin(s, e, "GET", path, handler, options...)
}

func PostGin(s *fuego.OpenAPI, e gin.IRouter, path string, handler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	return handleGin(s, e, "POST", path, handler, options...)
}

func Get[T, B any](s *fuego.OpenAPI, e gin.IRouter, path string, handler func(c ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	return handleFuego(s, e, "GET", path, handler, options...)
}

func Post[T, B any](s *fuego.OpenAPI, e gin.IRouter, path string, handler func(c ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	return handleFuego(s, e, "POST", path, handler, options...)
}

func handleFuego[T, B any](openapi *fuego.OpenAPI, e gin.IRouter, method, path string, fuegoHandler func(c ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	baseRoute := fuego.NewBaseRoute(method, path, fuegoHandler, openapi, options...)
	return handle(openapi, e, &fuego.Route[T, B]{BaseRoute: baseRoute}, GinHandler(fuegoHandler))
}

func handleGin(openapi *fuego.OpenAPI, e gin.IRouter, method, path string, ginHandler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	baseRoute := fuego.NewBaseRoute(method, path, ginHandler, openapi, options...)
	return handle(openapi, e, &fuego.Route[any, any]{BaseRoute: baseRoute}, ginHandler)
}

func handle[T, B any](openapi *fuego.OpenAPI, e gin.IRouter, route *fuego.Route[T, B], fuegoHandler gin.HandlerFunc) *fuego.Route[T, B] {
	if _, ok := e.(*gin.RouterGroup); ok {
		route.Path = e.(*gin.RouterGroup).BasePath() + route.Path
	}

	e.Handle(route.Method, route.Path, fuegoHandler)

	err := route.RegisterOpenAPIOperation(openapi)
	if err != nil {
		slog.Warn("error documenting openapi operation", "error", err)
	}

	return route
}

// Convert a Fuego handler to a Gin handler.
func GinHandler[B, T any](handler func(c ContextWithBody[B]) (T, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := &contextWithBody[B]{
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
	}
}
