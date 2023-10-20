package op

import (
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
	fullPath := path
	if isGo1_22 {
		fullPath = method + " " + path
	}
	s.logger.Debug("registering openapi controller " + fullPath)
	s.mux.Handle(fullPath, withMiddlewares(http.HandlerFunc(httpHandler[T, B](s, controller)), s.middlewares...))

	RegisterOpenAPIOperation[T, B](s, method, path)

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
	fullPath := path
	if isGo1_22 {
		fullPath = method + " " + path
	}
	s.logger.Debug("registering standard controller " + fullPath)
	s.mux.Handle(fullPath, withMiddlewares(http.HandlerFunc(controller), s.middlewares...))

	RegisterOpenAPIOperation[any, any](s, method, path)

	return Route[any, any]{}
}

func withMiddlewares(controller http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		controller = middleware(controller)
	}
	return controller
}
