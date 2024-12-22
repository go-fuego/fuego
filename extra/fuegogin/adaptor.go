package fuegogin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
)

func GetGin(engine *fuego.Engine, ginRouter gin.IRouter, path string, handler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	return handleGin(engine, ginRouter, http.MethodGet, path, handler, options...)
}

func PostGin(engine *fuego.Engine, ginRouter gin.IRouter, path string, handler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	return handleGin(engine, ginRouter, http.MethodPost, path, handler, options...)
}

func Get[T, B any](engine *fuego.Engine, ginRouter gin.IRouter, path string, handler func(c fuego.ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	return handleFuego(engine, ginRouter, http.MethodGet, path, handler, options...)
}

func Post[T, B any](engine *fuego.Engine, ginRouter gin.IRouter, path string, handler func(c fuego.ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	return handleFuego(engine, ginRouter, http.MethodPost, path, handler, options...)
}

func handleFuego[T, B any](engine *fuego.Engine, ginRouter gin.IRouter, method, path string, fuegoHandler func(c fuego.ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	baseRoute := fuego.NewBaseRoute(method, path, fuegoHandler, engine.OpenAPI, options...)
	return handle(engine, ginRouter, &fuego.Route[T, B]{BaseRoute: baseRoute}, GinHandler(engine, fuegoHandler, baseRoute))
}

func handleGin(engine *fuego.Engine, ginRouter gin.IRouter, method, path string, ginHandler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	baseRoute := fuego.NewBaseRoute(method, path, ginHandler, engine.OpenAPI, options...)
	return handle(engine, ginRouter, &fuego.Route[any, any]{BaseRoute: baseRoute}, ginHandler)
}

func handle[T, B any](engine *fuego.Engine, ginRouter gin.IRouter, route *fuego.Route[T, B], ginHandler gin.HandlerFunc) *fuego.Route[T, B] {
	if _, ok := ginRouter.(*gin.RouterGroup); ok {
		route.Path = ginRouter.(*gin.RouterGroup).BasePath() + route.Path
	}

	ginRouter.Handle(route.Method, route.Path, ginHandler)

	err := route.RegisterOpenAPIOperation(engine.OpenAPI)
	if err != nil {
		slog.Warn("error documenting openapi operation", "error", err)
	}

	return route
}

// Convert a Fuego handler to a Gin handler.
func GinHandler[B, T any](engine *fuego.Engine, handler func(c fuego.ContextWithBody[B]) (T, error), route fuego.BaseRoute) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := &ginContext[B]{
			CommonContext: internal.CommonContext[B]{
				CommonCtx:     c,
				UrlValues:     c.Request.URL.Query(),
				OpenAPIParams: route.Params,
			},
			ginCtx: c,
		}

		resp, err := handler(context)
		if err != nil {
			err = engine.ErrorHandler(err)
			c.JSON(getErrorCode(err), err)
			return
		}

		if c.Request.Header.Get("Accept") == "application/xml" {
			c.XML(200, resp)
			return
		}

		c.JSON(200, resp)
	}
}

func getErrorCode(err error) int {
	var status fuego.ErrorWithStatus
	if errors.As(err, &status) {
		return status.StatusCode()
	}
	return 500
}
