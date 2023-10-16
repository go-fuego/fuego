package op

import (
	"encoding/json"
	"fmt"
	"net/http"

	"log/slog"
)

func CtxEmpty() Ctx[any] {
	return Ctx[any]{}
}

// Context for the request. BodyType is the type of the request body.
type Ctx[BodyType any] struct {
	body BodyType
}

func (c Ctx[T]) PathParams() map[string]string {
	return nil
}

func (c Ctx[T]) QueryParams() map[string]string {
	return nil
}

func (c Ctx[T]) Body() T {
	return c.body
}

type RteCtx[T any] struct {
	ReturnType T
	BodyType   any
	ErrorType  error
}

func NewMux[T any, B any]() *Mux[T, B] {
	return &Mux[T, B]{}
}

type Mux[T any, B any] struct {
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

func Register[T any, B any](method string, path string, controller func(Ctx[B]) (T, error)) RteCtx[T] {
	http.HandleFunc(path, HttpHandler[T, B](controller))
	slog.Info(fmt.Sprintf("Registering %s %s", method, path))

	return RteCtx[T]{}
}

func Run(port string) {
	http.ListenAndServe(port, nil)
}

type Controller[ReturnType any, Body any] func(c Ctx[Body]) (ReturnType, error)

func HttpHandler[ReturnType any, Body any](controller any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		f, ok := controller.(func(c Ctx[Body]) (ReturnType, error))
		if !ok {
			var c Controller[ReturnType, Body]
			slog.Info("Controller types not ok",
				"type", fmt.Sprintf("%T", controller),
				"should be", fmt.Sprintf("%T", c))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Deserialize body
		var b Body

		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			slog.Info("Error decoding body", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Info("Decoded body", "body", b)

		ctx := Ctx[Body]{
			body: b,
		}

		ans, err := f(ctx)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(ans)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}
