package fuegomux

import (
	"net/http"
	"regexp"

	"github.com/gorilla/mux"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
)

// pathRegex matches gorilla/mux path params with regex constraints: {name:pattern}
var pathRegex = regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*):([^}]+)\}`)

// muxToFuegoRoute strips regex constraints from gorilla/mux paths.
// {id:[0-9]+} -> {id}
func muxToFuegoRoute(path string) string {
	return pathRegex.ReplaceAllString(path, `{$1}`)
}

// MuxRouter is the interface that gorilla/mux routers must satisfy.
// Both *mux.Router and subrouters from PathPrefix().Subrouter() satisfy this.
type MuxRouter interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
}

// OpenAPIHandler implements fuego.OpenAPIServable for gorilla/mux routers.
type OpenAPIHandler struct {
	Router MuxRouter
}

func (o *OpenAPIHandler) SpecHandler(e *fuego.Engine) {
	Get(e, o.Router, e.OpenAPI.Config.SpecURL, e.SpecHandler(), fuego.OptionHide(), fuego.OptionMiddleware(e.OpenAPI.Config.SwaggerMiddlewares...))
}

func (o *OpenAPIHandler) UIHandler(e *fuego.Engine) {
	GetMux(
		e,
		o.Router,
		e.OpenAPI.Config.SwaggerURL,
		e.OpenAPI.Config.UIHandler(e.OpenAPI.Config.SpecURL).ServeHTTP,
		fuego.OptionHide(),
		fuego.OptionMiddleware(e.OpenAPI.Config.SwaggerMiddlewares...),
	)
}

// --- Native mux handler registration (Level 1 & 2) ---

func AddMux(engine *fuego.Engine, muxRouter MuxRouter, method, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, method, path, handler, options...)
}

func GetMux(engine *fuego.Engine, muxRouter MuxRouter, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, http.MethodGet, path, handler, options...)
}

func PostMux(engine *fuego.Engine, muxRouter MuxRouter, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, http.MethodPost, path, handler, options...)
}

func PutMux(engine *fuego.Engine, muxRouter MuxRouter, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, http.MethodPut, path, handler, options...)
}

func DeleteMux(engine *fuego.Engine, muxRouter MuxRouter, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, http.MethodDelete, path, handler, options...)
}

func PatchMux(engine *fuego.Engine, muxRouter MuxRouter, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, http.MethodPatch, path, handler, options...)
}

func OptionsMux(engine *fuego.Engine, muxRouter MuxRouter, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, http.MethodOptions, path, handler, options...)
}

// --- Fuego typed handler registration (Level 3 & 4) ---

func Add[T, B, P any](engine *fuego.Engine, muxRouter MuxRouter, method, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, method, path, handler, options...)
}

func Get[T, B, P any](engine *fuego.Engine, muxRouter MuxRouter, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, http.MethodGet, path, handler, options...)
}

func Post[T, B, P any](engine *fuego.Engine, muxRouter MuxRouter, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, http.MethodPost, path, handler, options...)
}

func Put[T, B, P any](engine *fuego.Engine, muxRouter MuxRouter, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, http.MethodPut, path, handler, options...)
}

func Delete[T, B, P any](engine *fuego.Engine, muxRouter MuxRouter, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, http.MethodDelete, path, handler, options...)
}

func Patch[T, B, P any](engine *fuego.Engine, muxRouter MuxRouter, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, http.MethodPatch, path, handler, options...)
}

func Options[T, B, P any](engine *fuego.Engine, muxRouter MuxRouter, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, http.MethodOptions, path, handler, options...)
}

// --- Internal registration ---

func handleFuego[T, B, P any](engine *fuego.Engine, muxRouter MuxRouter, method, path string, fuegoHandler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	baseRoute := fuego.NewBaseRoute(method, muxToFuegoRoute(path), fuegoHandler, engine, options...)
	return fuego.Registers(engine, muxRouteRegisterer[T, B, P]{
		muxRouter:    muxRouter,
		route:        fuego.Route[T, B, P]{BaseRoute: baseRoute},
		httpHandler:  MuxHandler(engine, fuegoHandler, baseRoute),
		originalPath: path,
	})
}

func handleMux(engine *fuego.Engine, muxRouter MuxRouter, method, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	baseRoute := fuego.NewBaseRoute(method, muxToFuegoRoute(path), handler, engine, options...)
	return fuego.Registers(engine, muxRouteRegisterer[any, any, any]{
		muxRouter:    muxRouter,
		route:        fuego.Route[any, any, any]{BaseRoute: baseRoute},
		httpHandler:  handler,
		originalPath: path,
	})
}

// --- Route registerer ---

type muxRouteRegisterer[T, B, P any] struct {
	muxRouter    MuxRouter
	httpHandler  http.HandlerFunc
	route        fuego.Route[T, B, P]
	originalPath string
}

func (a muxRouteRegisterer[T, B, P]) Register() fuego.Route[T, B, P] {
	// Apply route-level middlewares (e.g. from fuego.OptionMiddleware).
	// Global/group middlewares are handled by gorilla/mux's router.Use().
	handler := applyMiddlewares(http.HandlerFunc(a.httpHandler), a.route.Middlewares...)
	muxRoute := a.muxRouter.HandleFunc(a.originalPath, handler.ServeHTTP).Methods(a.route.Method)

	// Get the full path template including any subrouter prefix.
	// If GetPathTemplate fails (which shouldn't happen for properly registered routes),
	// fall back to the pre-converted path already set on the route.
	if tpl, err := muxRoute.GetPathTemplate(); err == nil {
		a.route.Path = muxToFuegoRoute(tpl)
	}

	return a.route
}

// applyMiddlewares wraps an http.Handler with middlewares in the correct order
// (last middleware is innermost, first is outermost).
func applyMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// MuxHandler converts a Fuego handler to an http.HandlerFunc.
func MuxHandler[B, T, P any](engine *fuego.Engine, handler func(c fuego.Context[B, P]) (T, error), route fuego.BaseRoute) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &muxContext[B, P]{
			CommonContext: internal.CommonContext[B]{
				CommonCtx:         r.Context(),
				UrlValues:         r.URL.Query(),
				OpenAPIParams:     route.Params,
				DefaultStatusCode: route.DefaultStatusCode,
			},
			req: r,
			res: w,
		}
		fuego.Flow(engine, ctx, handler)
	}
}
