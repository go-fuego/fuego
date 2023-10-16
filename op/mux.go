package op

import (
	"fmt"
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
	s.mux.HandleFunc(path, httpHandler[T, B](controller))
	slog.Info(fmt.Sprintf("Registering %s %s", method, path))

	return Route[T, B]{}
}
