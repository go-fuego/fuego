package fuegogin

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-fuego/fuego"
	"github.com/stretchr/testify/assert"
)

type wrappedRouter struct {
	router gin.IRouter
}

func (m *wrappedRouter) Use(handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.router.Use(handlers...)
}

func (m *wrappedRouter) Handle(method, path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.router.Handle(method, path, handlers...)
}

func (m *wrappedRouter) Any(path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.router.Any(path, handlers...)
}

func (m *wrappedRouter) GET(path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.router.GET(path, handlers...)
}

func (m *wrappedRouter) POST(path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.router.POST(path, handlers...)
}

func (m *wrappedRouter) DELETE(path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.router.DELETE(path, handlers...)
}

func (m *wrappedRouter) PATCH(path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.router.PATCH(path, handlers...)
}

func (m *wrappedRouter) PUT(path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.router.PUT(path, handlers...)
}

func (m *wrappedRouter) OPTIONS(path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.router.OPTIONS(path, handlers...)
}

func (m *wrappedRouter) HEAD(path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.router.HEAD(path, handlers...)
}

func (m *wrappedRouter) Match(methods []string, path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.router.Match(methods, path, handlers...)
}

func (m *wrappedRouter) StaticFile(path, filepath string) gin.IRoutes {
	return m.router.StaticFile(path, filepath)
}

func (m *wrappedRouter) StaticFileFS(path, filepath string, fs http.FileSystem) gin.IRoutes {
	return m.router.StaticFileFS(path, filepath, fs)
}

func (m *wrappedRouter) Static(path, filepath string) gin.IRoutes {
	return m.router.Static(path, filepath)
}

func (m *wrappedRouter) StaticFS(path string, fs http.FileSystem) gin.IRoutes {
	return m.router.StaticFS(path, fs)
}

func (m *wrappedRouter) Group(path string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return m.router.Group(path, handlers...)
}

func (m *wrappedRouter) BasePath() string {
	if grouped, ok := m.router.(GroupedRouter); ok {
		return grouped.BasePath()
	}

	return ""
}

func dummyHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}

func TestWrappedRouterGroup(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		e := fuego.NewEngine()
		r := gin.New()
		groupV1 := r.Group("/api/v1")

		PostGin(e, &wrappedRouter{router: groupV1}, "/test", dummyHandler)

		assert.NotNil(t, e.OpenAPI.Description().Paths.Value("/api/v1/test"))
	})

	t.Run("withParam", func(t *testing.T) {
		e := fuego.NewEngine()
		r := gin.New()
		groupV1 := r.Group("/api/v1")

		PostGin(e, &wrappedRouter{router: groupV1}, "/test/:id", dummyHandler)

		assert.NotNil(t, e.OpenAPI.Description().Paths.Value("/api/v1/test/{id}"))
	})

	t.Run("withGroupParam", func(t *testing.T) {
		e := fuego.NewEngine()
		r := gin.New()
		groupV1 := r.Group("/api/v1/:id")

		PostGin(e, &wrappedRouter{router: groupV1}, "/test", dummyHandler)

		assert.NotNil(t, e.OpenAPI.Description().Paths.Value("/api/v1/{id}/test"))
	})

	t.Run("defaultGin", func(t *testing.T) {
		e := fuego.NewEngine()
		r := gin.New()

		PostGin(e, r, "/test", dummyHandler)

		// Basic gin engine can be casted to fuegogin.GroupedRouter too
		// so we expect no extra / in path.
		assert.Nil(t, e.OpenAPI.Description().Paths.Value("//test"))
		assert.NotNil(t, e.OpenAPI.Description().Paths.Value("/test"))
	})

	t.Run("emptyGroup", func(t *testing.T) {
		e := fuego.NewEngine()
		r := gin.New()
		emptyGroup := r.Group("")

		PostGin(e, emptyGroup, "/test", dummyHandler)

		assert.Nil(t, e.OpenAPI.Description().Paths.Value("//test"))
		assert.NotNil(t, e.OpenAPI.Description().Paths.Value("/test"))
	})

	t.Run("slashGroup", func(t *testing.T) {
		e := fuego.NewEngine()
		r := gin.New()
		slashGroup := r.Group("/")

		PostGin(e, slashGroup, "/test", dummyHandler)

		assert.Nil(t, e.OpenAPI.Description().Paths.Value("//test"))
		assert.NotNil(t, e.OpenAPI.Description().Paths.Value("/test"))
	})
}
