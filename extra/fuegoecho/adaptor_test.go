package fuegoecho

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func TestFuegoPathWithEchoPathParam(t *testing.T) {
	basePath := "/api"
	apiPath := "/path/:id"
	e := fuego.NewEngine()
	echoRouter := echo.New()
	group := echoRouter.Group(basePath)

	Get(e, group, apiPath, func(c fuego.Context[any, any]) (string, error) { return "ok", nil },
		option.Path("id", "ID"),
	)

	spec := e.OutputOpenAPISpec()
	specPath := spec.Paths.Find("/api/path/{id}")

	routes := echoRouter.Routes()
	assert.Len(t, routes, 1, "Expected exactly one route registered in Echo")
	assert.Equal(t, routes[0].Path, basePath+apiPath)

	assert.NotNil(t, specPath,
		"Expected path '/api/path/{id}' to be registered in OpenAPI spec")
}

func TestFuegoPathWithTwoEchoPathParams(t *testing.T) {
	basePath := "/api"
	apiPath := "/path/:id1/foo/:id2"
	e := fuego.NewEngine()
	echoRouter := echo.New()
	group := echoRouter.Group(basePath)

	Get(e, group, apiPath, func(c fuego.Context[any, any]) (string, error) { return "ok", nil },
		option.Path("id1", "First ID"),
		option.Path("id2", "Second ID"),
	)

	spec := e.OutputOpenAPISpec()
	specPath := spec.Paths.Find("/api/path/{id1}/foo/{id2}")

	routes := echoRouter.Routes()
	assert.Len(t, routes, 1, "Expected exactly one route registered in Echo")
	assert.Equal(t, routes[0].Path, basePath+apiPath)

	assert.NotNil(t, specPath,
		"Expected path '/api/path/{id1}/foo/{id2}' to be registered in OpenAPI spec")
}
