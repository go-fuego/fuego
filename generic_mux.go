package fuego

import (
	"log/slog"
)

// Registerer is an interface that allows registering routes.
// It can be implementable by any router.
type Registerer[T, B, P any] interface {
	Register() Route[T, B, P]
}

func Registers[B, T, P any](engine *Engine, a Registerer[B, T, P]) *Route[B, T, P] {
	route := a.Register()
	err := route.RegisterOpenAPIOperation(engine.OpenAPI)
	if err != nil {
		slog.Warn("error documenting openapi operation", "error", err)
	}
	return &route
}
