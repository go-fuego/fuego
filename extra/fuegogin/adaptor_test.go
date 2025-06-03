package fuegogin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/go-fuego/fuego"
)

func setupTestEngine() (*fuego.Engine, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	engine := fuego.NewEngine()
	ginEngine := gin.New()
	return engine, ginEngine
}

func TestOpenAPIHandler_SpecHandler(t *testing.T) {
	engine, ginEngine := setupTestEngine()
	handler := &OpenAPIHandler{GinEngine: ginEngine}

	// This should not panic
	require.NotPanics(t, func() {
		handler.SpecHandler(engine)
	})
}

func TestOpenAPIHandler_UIHandler(t *testing.T) {
	engine, ginEngine := setupTestEngine()
	handler := &OpenAPIHandler{GinEngine: ginEngine}

	// This should not panic
	require.NotPanics(t, func() {
		handler.UIHandler(engine)
	})
}

func TestGetGin(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	ginHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hello"})
	}

	route := GetGin(engine, ginEngine, "/test", ginHandler)

	require.NotNil(t, route)
	require.Equal(t, http.MethodGet, route.Method)
	require.Equal(t, "/test", route.Path)
}

func TestPostGin(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	ginHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "created"})
	}

	route := PostGin(engine, ginEngine, "/test", ginHandler)

	require.NotNil(t, route)
	require.Equal(t, http.MethodPost, route.Method)
	require.Equal(t, "/test", route.Path)
}

func TestPutGin(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	ginHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "updated"})
	}

	route := PutGin(engine, ginEngine, "/test", ginHandler)

	require.NotNil(t, route)
	require.Equal(t, http.MethodPut, route.Method)
	require.Equal(t, "/test", route.Path)
}

func TestDeleteGin(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	ginHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "deleted"})
	}

	route := DeleteGin(engine, ginEngine, "/test", ginHandler)

	require.NotNil(t, route)
	require.Equal(t, http.MethodDelete, route.Method)
	require.Equal(t, "/test", route.Path)
}

func TestPatchGin(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	ginHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "patched"})
	}

	route := PatchGin(engine, ginEngine, "/test", ginHandler)

	require.NotNil(t, route)
	require.Equal(t, http.MethodPatch, route.Method)
	require.Equal(t, "/test", route.Path)
}

func TestOptionsGin(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	ginHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "options"})
	}

	route := OptionsGin(engine, ginEngine, "/test", ginHandler)

	require.NotNil(t, route)
	require.Equal(t, http.MethodOptions, route.Method)
	require.Equal(t, "/test", route.Path)
}

func TestGet(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	handler := func(c fuego.Context[any, any]) (string, error) {
		return "hello", nil
	}

	route := Get(engine, ginEngine, "/test", handler)

	require.NotNil(t, route)
	require.Equal(t, http.MethodGet, route.Method)
	require.Equal(t, "/test", route.Path)
}

func TestPost(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	handler := func(c fuego.Context[any, any]) (string, error) {
		return "created", nil
	}

	route := Post(engine, ginEngine, "/test", handler)

	require.NotNil(t, route)
	require.Equal(t, http.MethodPost, route.Method)
	require.Equal(t, "/test", route.Path)
}

func TestPut(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	handler := func(c fuego.Context[any, any]) (string, error) {
		return "updated", nil
	}

	route := Put(engine, ginEngine, "/test", handler)

	require.NotNil(t, route)
	require.Equal(t, http.MethodPut, route.Method)
	require.Equal(t, "/test", route.Path)
}

func TestDelete(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	handler := func(c fuego.Context[any, any]) (string, error) {
		return "deleted", nil
	}

	route := Delete(engine, ginEngine, "/test", handler)

	require.NotNil(t, route)
	require.Equal(t, http.MethodDelete, route.Method)
	require.Equal(t, "/test", route.Path)
}

func TestPatch(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	handler := func(c fuego.Context[any, any]) (string, error) {
		return "patched", nil
	}

	route := Patch(engine, ginEngine, "/test", handler)

	require.NotNil(t, route)
	require.Equal(t, http.MethodPatch, route.Method)
	require.Equal(t, "/test", route.Path)
}

func TestOptions(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	handler := func(c fuego.Context[any, any]) (string, error) {
		return "options", nil
	}

	route := Options(engine, ginEngine, "/test", handler)

	require.NotNil(t, route)
	require.Equal(t, http.MethodOptions, route.Method)
	require.Equal(t, "/test", route.Path)
}

func TestGinToFuegoRoute(t *testing.T) {
	tests := []struct {
		name     string
		ginPath  string
		expected string
	}{
		{
			name:     "simple path",
			ginPath:  "/users",
			expected: "/users",
		},
		{
			name:     "path with single parameter",
			ginPath:  "/users/:id",
			expected: "/users/{id}",
		},
		{
			name:     "path with multiple parameters",
			ginPath:  "/users/:id/posts/:postId",
			expected: "/users/{id}/posts/{postId}",
		},
		{
			name:     "path with parameter at end",
			ginPath:  "/api/v1/users/:userId",
			expected: "/api/v1/users/{userId}",
		},
		{
			name:     "path with underscore in parameter",
			ginPath:  "/users/:user_id",
			expected: "/users/{user_id}",
		},
		{
			name:     "path with numbers in parameter",
			ginPath:  "/api/:version123/users",
			expected: "/api/{version123}/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ginToFuegoRoute(tt.ginPath)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGinRouteRegisterer_Register(t *testing.T) {
	t.Run("basic registration", func(t *testing.T) {
		engine, ginEngine := setupTestEngine()

		ginHandler := func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test"})
		}

		baseRoute := fuego.NewBaseRoute(http.MethodGet, "/test", ginHandler, engine)
		registerer := ginRouteRegisterer[any, any, any]{
			ginRouter:    ginEngine,
			ginHandler:   ginHandler,
			route:        fuego.Route[any, any, any]{BaseRoute: baseRoute},
			originalPath: "/test",
		}

		route := registerer.Register()
		require.Equal(t, "/test", route.Path)
		require.Equal(t, http.MethodGet, route.Method)
	})

	t.Run("registration with router group", func(t *testing.T) {
		engine, ginEngine := setupTestEngine()

		// Create a router group
		group := ginEngine.Group("/api/v1")

		ginHandler := func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test"})
		}

		baseRoute := fuego.NewBaseRoute(http.MethodGet, "/test", ginHandler, engine)
		registerer := ginRouteRegisterer[any, any, any]{
			ginRouter:    group,
			ginHandler:   ginHandler,
			route:        fuego.Route[any, any, any]{BaseRoute: baseRoute},
			originalPath: "/test",
		}

		route := registerer.Register()
		require.Equal(t, "/api/v1/test", route.Path)
		require.Equal(t, http.MethodGet, route.Method)
	})
}

func TestGinHandler(t *testing.T) {
	engine, _ := setupTestEngine()

	handler := func(c fuego.Context[any, any]) (string, error) {
		return "hello world", nil
	}

	baseRoute := fuego.NewBaseRoute(http.MethodGet, "/test", handler, engine)
	ginHandler := GinHandler(engine, handler, baseRoute)

	require.NotNil(t, ginHandler)

	// Test that the gin handler can be called without panicking
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	require.NotPanics(t, func() {
		ginHandler(c)
	})
}

// Test custom router that implements GroupedRouter interface
type testGroupedRouter struct {
	*gin.RouterGroup
	basePath string
}

func (t *testGroupedRouter) BasePath() string {
	return t.basePath
}

func TestGinRouteRegisterer_RegisterWithCustomGroupedRouter(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	customRouter := &testGroupedRouter{
		RouterGroup: ginEngine.Group("/custom"),
		basePath:    "/custom",
	}

	ginHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	}

	baseRoute := fuego.NewBaseRoute(http.MethodGet, "/test", ginHandler, engine)
	registerer := ginRouteRegisterer[any, any, any]{
		ginRouter:    customRouter,
		ginHandler:   ginHandler,
		route:        fuego.Route[any, any, any]{BaseRoute: baseRoute},
		originalPath: "/test",
	}

	route := registerer.Register()
	require.Equal(t, "/custom/test", route.Path)
	require.Equal(t, http.MethodGet, route.Method)
}

func TestGinRouteRegisterer_RegisterWithRootGroup(t *testing.T) {
	engine, ginEngine := setupTestEngine()

	// Test with root group (should not modify path)
	customRouter := &testGroupedRouter{
		RouterGroup: ginEngine.Group("/"),
		basePath:    "/",
	}

	ginHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	}

	baseRoute := fuego.NewBaseRoute(http.MethodGet, "/test", ginHandler, engine)
	registerer := ginRouteRegisterer[any, any, any]{
		ginRouter:    customRouter,
		ginHandler:   ginHandler,
		route:        fuego.Route[any, any, any]{BaseRoute: baseRoute},
		originalPath: "/test",
	}

	route := registerer.Register()
	require.Equal(t, "/test", route.Path) // Should not be modified
	require.Equal(t, http.MethodGet, route.Method)
}
