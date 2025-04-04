package fuegogin

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-fuego/fuego"
	"gotest.tools/v3/assert"
)

func TestGinPathWithGinGroup(t *testing.T) {
	basePath := "/api"
	apiPath := "/path"
	completePath := basePath + apiPath

	e := fuego.NewEngine()
	ginRouter := gin.New()
	group := ginRouter.Group(basePath)

	GetGin(e, group, apiPath, gin.HandlerFunc(func(ctx *gin.Context) {
		result, err := func(ctx *gin.Context) (string, error) { return "ok", nil }(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"result": result})
	}))

	routes := ginRouter.Routes()
	for _, route := range routes {
		assert.Equal(t, route.Path, completePath)
	}
}

func TestFuegoPathWithGinPathParam(t *testing.T) {
	basePath := "/api"
	apiPath := "/path/:id"
	e := fuego.NewEngine()
	ginRouter := gin.New()
	group := ginRouter.Group(basePath)

	Get(e, group, apiPath, func(c fuego.ContextWithBody[any]) (string, error) { return "ok", nil })

	spec := e.OutputOpenAPISpec()
	specPath := spec.Paths.Find("/api/path/{id}")

	routes := ginRouter.Routes()
	assert.Assert(t, len(routes) == 1, "Expected exactly one route registered in Gin")
	assert.Equal(t, routes[0].Path, basePath+apiPath)

	assert.Check(t, specPath != nil,
		"Expected path '/api/path/{id}' to be registered in OpenAPI spec")
}

func TestFuegoPathWithGinGroup(t *testing.T) {
	basePath := "/api"
	apiPath := "/path"
	completePath := basePath + apiPath

	e := fuego.NewEngine()
	ginRouter := gin.New()
	group := ginRouter.Group(basePath)

	Get(e, group, apiPath, func(c fuego.ContextWithBody[any]) (string, error) { return "ok", nil })

	routes := ginRouter.Routes()
	for _, route := range routes {
		assert.Equal(t, route.Path, completePath)
	}
}
