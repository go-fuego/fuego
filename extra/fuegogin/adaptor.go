package fuegogin

import (
	"github.com/gin-gonic/gin"

	"github.com/go-fuego/fuego"
)

func Get[T, B any](
	s *fuego.Server,
	e *gin.Engine,
	path string,
	handler func(c *ContextWithBody[T]) (B, error),
	options ...func(*fuego.BaseRoute),
) *fuego.Route[B, T] {
	return Handle(s, e, "GET", path, handler)
}

func Post[T, B any](
	s *fuego.Server,
	e *gin.Engine,
	path string,
	handler func(c *ContextWithBody[T]) (B, error),
	options ...func(*fuego.BaseRoute),
) *fuego.Route[B, T] {
	return Handle(s, e, "GET", path, handler)
}

func Handle[T, B any](
	s *fuego.Server,
	e *gin.Engine,
	method,
	path string,
	handler func(c *ContextWithBody[T]) (B, error),
	options ...func(*fuego.BaseRoute),
) *fuego.Route[B, T] {
	e.Handle(method, path, func(c *gin.Context) {
		ans, err := handler(&ContextWithBody[T]{})
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, ans)
	})

	// Also register the route with fuego!!!
	// Useful for the OpenAPI spec but also allows for to run Fuego in parallel.
	return fuego.Get(s, path, handler, options...)
}
