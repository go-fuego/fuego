package fuegogin

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
)

var pathRegex = regexp.MustCompile(`:([a-zA-Z0-9_]+)`)

type OpenAPIHandler struct {
	GinEngine *gin.Engine
}

func (o *OpenAPIHandler) SpecHandler(e *fuego.Engine) {
	Get(e, o.GinEngine, e.OpenAPI.Config.SpecURL, e.SpecHandler(), fuego.OptionHide())
}

func (o *OpenAPIHandler) UIHandler(e *fuego.Engine) {
	GetGin(
		e,
		o.GinEngine,
		e.OpenAPI.Config.SwaggerURL+"/",
		gin.WrapH(e.OpenAPI.Config.UIHandler(e.OpenAPI.Config.SpecURL)),
		fuego.OptionHide(),
	)
}

func GetGin(engine *fuego.Engine, ginRouter gin.IRouter, path string, handler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	return handleGin(engine, ginRouter, http.MethodGet, path, handler, options...)
}

func PostGin(engine *fuego.Engine, ginRouter gin.IRouter, path string, handler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	return handleGin(engine, ginRouter, http.MethodPost, path, handler, options...)
}

func PutGin(engine *fuego.Engine, ginRouter gin.IRouter, path string, handler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	return handleGin(engine, ginRouter, http.MethodPut, path, handler, options...)
}

func DeleteGin(engine *fuego.Engine, ginRouter gin.IRouter, path string, handler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	return handleGin(engine, ginRouter, http.MethodDelete, path, handler, options...)
}

func PatchGin(engine *fuego.Engine, ginRouter gin.IRouter, path string, handler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	return handleGin(engine, ginRouter, http.MethodPatch, path, handler, options...)
}

func OptionsGin(engine *fuego.Engine, ginRouter gin.IRouter, path string, handler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	return handleGin(engine, ginRouter, http.MethodOptions, path, handler, options...)
}

func Get[T, B any](engine *fuego.Engine, ginRouter gin.IRouter, path string, handler func(c fuego.ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	return handleFuego(engine, ginRouter, http.MethodGet, path, handler, options...)
}

func Post[T, B any](engine *fuego.Engine, ginRouter gin.IRouter, path string, handler func(c fuego.ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	return handleFuego(engine, ginRouter, http.MethodPost, path, handler, options...)
}

func Put[T, B any](engine *fuego.Engine, ginRouter gin.IRouter, path string, handler func(c fuego.ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	return handleFuego(engine, ginRouter, http.MethodPut, path, handler, options...)
}

func Delete[T, B any](engine *fuego.Engine, ginRouter gin.IRouter, path string, handler func(c fuego.ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	return handleFuego(engine, ginRouter, http.MethodDelete, path, handler, options...)
}

func Patch[T, B any](engine *fuego.Engine, ginRouter gin.IRouter, path string, handler func(c fuego.ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	return handleFuego(engine, ginRouter, http.MethodPatch, path, handler, options...)
}

func Options[T, B any](engine *fuego.Engine, ginRouter gin.IRouter, path string, handler func(c fuego.ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	return handleFuego(engine, ginRouter, http.MethodOptions, path, handler, options...)
}

func handleFuego[T, B any](engine *fuego.Engine, ginRouter gin.IRouter, method, path string, fuegoHandler func(c fuego.ContextWithBody[B]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B] {
	baseRoute := fuego.NewBaseRoute(method, ginToFuegoRoute(path), fuegoHandler, engine, options...)
	return fuego.Registers(engine, ginRouteRegisterer[T, B]{
		ginRouter:    ginRouter,
		route:        fuego.Route[T, B]{BaseRoute: baseRoute},
		ginHandler:   GinHandler(engine, fuegoHandler, baseRoute),
		originalPath: path,
	})
}

func handleGin(engine *fuego.Engine, ginRouter gin.IRouter, method, path string, ginHandler gin.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any] {
	baseRoute := fuego.NewBaseRoute(method, ginToFuegoRoute(path), ginHandler, engine, options...)
	return fuego.Registers(engine, ginRouteRegisterer[any, any]{
		ginRouter:    ginRouter,
		route:        fuego.Route[any, any]{BaseRoute: baseRoute},
		ginHandler:   ginHandler,
		originalPath: path,
	})
}

func ginToFuegoRoute(path string) string {
	return pathRegex.ReplaceAllString(path, `{$1}`)
}

type ginRouteRegisterer[T, B any] struct {
	ginRouter    gin.IRouter
	ginHandler   gin.HandlerFunc
	route        fuego.Route[T, B]
	originalPath string
}

func (a ginRouteRegisterer[T, B]) Register() fuego.Route[T, B] {
	// We must register the gin handler first, so that the gin router can
	// mutate the route path if it is a RouterGroup.
	// This is because gin groups will prepend the group path to the route path itself.
	a.ginRouter.Handle(a.route.Method, a.originalPath, a.ginHandler)

	if _, ok := a.ginRouter.(*gin.RouterGroup); ok {
		a.route.Path = a.ginRouter.(*gin.RouterGroup).BasePath() + a.route.Path
	}

	return a.route
}

// Convert a Fuego handler to a Gin handler.
func GinHandler[B, T any](engine *fuego.Engine, handler func(c fuego.ContextWithBody[B]) (T, error), route fuego.BaseRoute) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := &ginContext[B]{
			CommonContext: internal.CommonContext[B]{
				CommonCtx:         c,
				UrlValues:         c.Request.URL.Query(),
				OpenAPIParams:     route.Params,
				DefaultStatusCode: route.DefaultStatusCode,
			},
			ginCtx: c,
		}

		fuego.Flow(engine, context, handler)
	}
}
