package op

import (
	"fmt"
	"log/slog"
	"net/http"
)

func NewMux[ResponseBody any, RequestBody any]() *Mux[ResponseBody, RequestBody] {
	mux := http.NewServeMux()
	return &Mux[ResponseBody, RequestBody]{
		mux: mux,
	}
}

var defaultMux = http.NewServeMux()

type Mux[ResponseBody any, RequestBody any] struct {
	mux *http.ServeMux
}

func (m Mux[T, B]) Get(path string, controller func(Ctx[B]) (T, error)) RteCtx[T] {
	return Register[T](http.MethodGet, path, controller)
}

func Get[T any, B any](path string, controller func(Ctx[B]) (T, error)) RteCtx[T] {
	return Register[T](http.MethodGet, path, controller)
}

func Post[T any, B any](path string, controller func(Ctx[B]) (T, error)) RteCtx[T] {
	return Register[T](http.MethodPost, path, controller)
}

// Registers route into the default mux.
func Register[T any, B any](method string, path string, controller func(Ctx[B]) (T, error)) RteCtx[T] {
	defaultMux.HandleFunc(path, httpHandler[T, B](controller))
	slog.Info(fmt.Sprintf("Registering %s %s", method, path))

	return RteCtx[T]{}
}
