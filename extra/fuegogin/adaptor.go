package fuegogin

import (
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
	return fuego.Registers(engine, ginRouteRegisterer[T, B]{
		ginRouter:  ginRouter,
		route:      fuego.Route[T, B]{BaseRoute: baseRoute},
		ginHandler: GinHandler(engine, fuegoHandler, baseRoute),
	})
}

func handleGin(engine *fuego.Engine, ginRouter gin.IRouter, method, path string, ginHandler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	baseRoute := fuego.NewBaseRoute(method, path, ginHandler, engine.OpenAPI, options...)
	return fuego.Registers(engine, ginRouteRegisterer[any, any]{
		ginRouter:  ginRouter,
		route:      fuego.Route[any, any]{BaseRoute: baseRoute},
		ginHandler: ginHandler,
	})
}

type ginRouteRegisterer[T, B any] struct {
	ginRouter  gin.IRouter
	ginHandler gin.HandlerFunc
	route      fuego.Route[T, B]
}

func (a ginRouteRegisterer[T, B]) Register() fuego.Route[T, B] {
	if _, ok := a.ginRouter.(*gin.RouterGroup); ok {
		a.route.Path = a.ginRouter.(*gin.RouterGroup).BasePath() + a.route.Path
	}

	a.ginRouter.Handle(a.route.Method, a.route.Path, a.ginHandler)

	return a.route
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

		fuego.Flow(engine, context, handler)
	}
}
