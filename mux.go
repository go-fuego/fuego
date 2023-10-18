package op

import (
	"log/slog"
	"net/http"
)

type Route[ResponseBody any, RequestBody any] struct {
	ReturnType ResponseBody
	BodyType   ResponseBody
	ErrorType  error
}

func Get[T any, B any](s *Server, path string, controller func(Ctx[B]) (T, error)) Route[T, B] {
	return Register[T](s, http.MethodGet, path, controller)
}

func Post[T any, B any](s *Server, path string, controller func(Ctx[B]) (T, error)) Route[T, B] {
	return Register[T](s, http.MethodPost, path, controller)
}

// Registers route into the default mux.
func Register[T any, B any](s *Server, method string, path string, controller func(Ctx[B]) (T, error)) Route[T, B] {
	// fullRegistration := method + " " + path // TODO: switch when go 1.22 is released
	fullRegistration := path
	slog.Debug("registering openapi controller " + fullRegistration)
	s.mux.Handle(fullRegistration, withMiddlewares(http.HandlerFunc(httpHandler[T, B](s, controller)), s.middlewares...))

	RegisterOpenAPIOperation(s, method, fullRegistration)

	return Route[T, B]{}
}

func UseStd(s *Server, middlewares ...func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func GetStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request)) Route[any, any] {
	return RegisterStd(s, http.MethodGet, path, controller)
}

// RegisterStd registers a standard http handler into the default mux.
func RegisterStd(s *Server, method string, path string, controller func(http.ResponseWriter, *http.Request)) Route[any, any] {
	// fullRegistration := method + " " + path // TODO: switch when go 1.22 is released
	fullRegistration := path
	slog.Debug("registering standard controller " + fullRegistration)
	s.mux.Handle(fullRegistration, withMiddlewares(http.HandlerFunc(controller), s.middlewares...))

	RegisterOpenAPIOperation(s, method, fullRegistration)

	return Route[any, any]{}
}

func withMiddlewares(controller http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		controller = middleware(controller)
	}
	return controller
}
