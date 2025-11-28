package fuegoecho

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
)

var pathRegex = regexp.MustCompile(`:([a-zA-Z0-9_]+)`)

type OpenAPIHandler struct {
	Echo *echo.Echo
}

func (o *OpenAPIHandler) SpecHandler(e *fuego.Engine) {
	Get(e, o.Echo, e.OpenAPI.Config.SpecURL, e.SpecHandler(), fuego.OptionHide())
}

func (o *OpenAPIHandler) UIHandler(e *fuego.Engine) {
	GetEcho(
		e,
		o.Echo,
		e.OpenAPI.Config.SwaggerURL,
		echo.WrapHandler(e.OpenAPI.Config.UIHandler(e.OpenAPI.Config.SpecURL)),
		fuego.OptionHide(),
	)
}

type echoIRouter interface {
	Add(method, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route
}

func AddEcho(engine *fuego.Engine, echoRouter echoIRouter,
	method, path string, handler echo.HandlerFunc,
	options ...func(*fuego.BaseRoute),
) *fuego.Route[any, any, any] {
	return handleEcho(engine, echoRouter, method, path, handler, options...)
}

func GetEcho(engine *fuego.Engine, echoRouter echoIRouter,
	path string, handler echo.HandlerFunc,
	options ...func(*fuego.BaseRoute),
) *fuego.Route[any, any, any] {
	return handleEcho(engine, echoRouter, http.MethodGet, path, handler, options...)
}

func PostEcho(engine *fuego.Engine, echoRouter echoIRouter,
	path string, handler echo.HandlerFunc,
	options ...func(*fuego.BaseRoute),
) *fuego.Route[any, any, any] {
	return handleEcho(engine, echoRouter, http.MethodPost, path, handler, options...)
}

func PutEcho(engine *fuego.Engine, echoRouter echoIRouter,
	path string, handler echo.HandlerFunc,
	options ...func(*fuego.BaseRoute),
) *fuego.Route[any, any, any] {
	return handleEcho(engine, echoRouter, http.MethodPut, path, handler, options...)
}

func PatchEcho(engine *fuego.Engine, echoRouter echoIRouter,
	path string, handler echo.HandlerFunc,
	options ...func(*fuego.BaseRoute),
) *fuego.Route[any, any, any] {
	return handleEcho(engine, echoRouter, http.MethodPatch, path, handler, options...)
}

func DeleteEcho(engine *fuego.Engine, echoRouter echoIRouter,
	path string, handler echo.HandlerFunc,
	options ...func(*fuego.BaseRoute),
) *fuego.Route[any, any, any] {
	return handleEcho(engine, echoRouter, http.MethodDelete, path, handler, options...)
}

func Add[T, B, P any](engine *fuego.Engine, echoRouter echoIRouter, method, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, echoRouter, method, path, handler, options...)
}

func Get[T, B, P any](engine *fuego.Engine, echoRouter echoIRouter, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, echoRouter, http.MethodGet, path, handler, options...)
}

func Post[T, B, P any](engine *fuego.Engine, echoRouter echoIRouter, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, echoRouter, http.MethodPost, path, handler, options...)
}

func Put[T, B, P any](engine *fuego.Engine, echoRouter echoIRouter, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, echoRouter, http.MethodPut, path, handler, options...)
}

func Patch[T, B, P any](engine *fuego.Engine, echoRouter echoIRouter, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, echoRouter, http.MethodPatch, path, handler, options...)
}

func Delete[T, B, P any](engine *fuego.Engine, echoRouter echoIRouter, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, echoRouter, http.MethodDelete, path, handler, options...)
}

func handleFuego[T, B, P any](engine *fuego.Engine, echoRouter echoIRouter, method, path string, fuegoHandler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	baseRoute := fuego.NewBaseRoute(method, echoToFuegoRoute(path), fuegoHandler, engine, options...)
	return fuego.Registers(engine, echoRouteRegisterer[T, B, P]{
		echoRouter:   echoRouter,
		route:        fuego.Route[T, B, P]{BaseRoute: baseRoute},
		echoHandler:  EchoHandler(engine, fuegoHandler, baseRoute),
		originalPath: path,
	})
}

func handleEcho(engine *fuego.Engine, echoRouter echoIRouter, method, path string, echoHandler echo.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	baseRoute := fuego.NewBaseRoute(method, echoToFuegoRoute(path), echoHandler, engine, options...)
	return fuego.Registers(engine, echoRouteRegisterer[any, any, any]{
		echoRouter:   echoRouter,
		route:        fuego.Route[any, any, any]{BaseRoute: baseRoute},
		echoHandler:  echoHandler,
		originalPath: path,
	})
}

func echoToFuegoRoute(path string) string {
	return pathRegex.ReplaceAllString(path, `{$1}`)
}

type echoRouteRegisterer[T, B, P any] struct {
	echoRouter   echoIRouter
	echoHandler  echo.HandlerFunc
	route        fuego.Route[T, B, P]
	originalPath string
}

func (a echoRouteRegisterer[T, B, P]) Register() fuego.Route[T, B, P] {
	// We must register the echo handler first, so that the echo router can
	// mutate the route path if it is a Group.
	// This is because echo groups will prepend the group prefix to the route path itself.
	route := a.echoRouter.Add(a.route.Method, a.originalPath, a.echoHandler)
	// Echo returns the full path with group prefix, convert it to Fuego format
	a.route.Path = echoToFuegoRoute(route.Path)

	return a.route
}

// Convert a Fuego handler to an echo handler.
func EchoHandler[B, T, P any](engine *fuego.Engine, handler func(c fuego.Context[B, P]) (T, error), route fuego.BaseRoute) echo.HandlerFunc {
	return func(c echo.Context) error {
		context := &echoContext[B, P]{
			CommonContext: internal.CommonContext[B]{
				CommonCtx:         c.Request().Context(),
				UrlValues:         c.Request().URL.Query(),
				OpenAPIParams:     route.Params,
				DefaultStatusCode: route.DefaultStatusCode,
			},
			echoCtx: c,
		}
		fuego.Flow(engine, context, handler)
		return nil
	}
}
