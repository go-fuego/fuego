package fuego

import (
	"log/slog"
)

// Registerer is an interface that allows registering routes.
// It can be implementable by any router.
type Registerer[T, B any] interface {
	Register() Route[T, B]
}

func Registers[B, T any](engine *Engine, a Registerer[B, T]) *Route[B, T] {
	route := a.Register()
	err := route.RegisterOpenAPIOperation(engine.OpenAPI)
	if err != nil {
		slog.Warn("error documenting openapi operation", "error", err)
	}
	return &route
}
