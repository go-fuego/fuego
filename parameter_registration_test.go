package fuego

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RegisterOpenAPIOperation(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	s := NewServer()

	t.Run("Nil operation handling", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct{}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		route.Operation = nil
		err := route.RegisterParams()
		require.NoError(t, err)
		assert.NotNil(t, route.Operation)
	})

	t.Run("Register with params", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			QueryParam  string `query:"queryParam"`
			HeaderParam string `header:"headerParam"`
		}](
			http.MethodGet,
			"/some/path/{pathParam}",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.NoError(t, err)
		operation := route.Operation
		assert.NotNil(t, operation)
		assert.Len(t, operation.Parameters, 2)

		queryParam := operation.Parameters.GetByInAndName("query", "queryParam")
		assert.NotNil(t, queryParam)
		assert.Equal(t, "queryParam", queryParam.Name)

		headerParam := operation.Parameters.GetByInAndName("header", "headerParam")
		assert.NotNil(t, headerParam)
		assert.Equal(t, "headerParam", headerParam.Name)
	})

	t.Run("RegisterParams do not raise error with interface types", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, any](
			http.MethodGet,
			"/no-interfaces",
			handler,
			s.Engine,
			OptionDefaultStatusCode(201),
		)

		err := route.RegisterParams()
		require.NoError(t, err)
	})
}

func TestRegisterParams_DefaultTag(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	s := NewServer()

	t.Run("String default", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Name string `query:"name" default:"test"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.NoError(t, err)

		param := route.Operation.Parameters.GetByInAndName("query", "name")
		require.NotNil(t, param)
		require.NotNil(t, param.Schema)
		require.NotNil(t, param.Schema.Value)
		assert.Equal(t, "test", param.Schema.Value.Default)
	})

	t.Run("Int default", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Age int `query:"age" default:"25"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.NoError(t, err)

		param := route.Operation.Parameters.GetByInAndName("query", "age")
		require.NotNil(t, param)
		require.NotNil(t, param.Schema)
		require.NotNil(t, param.Schema.Value)
		assert.Equal(t, 25, param.Schema.Value.Default)
	})

	t.Run("Int64 default", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Limit int64 `query:"limit" default:"100"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.NoError(t, err)

		param := route.Operation.Parameters.GetByInAndName("query", "limit")
		require.NotNil(t, param)
		require.NotNil(t, param.Schema)
		require.NotNil(t, param.Schema.Value)
		// OpenAPI normalizes all integer types to int
		assert.Equal(t, 100, param.Schema.Value.Default)
	})

	t.Run("Uint default", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Count uint `query:"count" default:"50"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.NoError(t, err)

		param := route.Operation.Parameters.GetByInAndName("query", "count")
		require.NotNil(t, param)
		require.NotNil(t, param.Schema)
		require.NotNil(t, param.Schema.Value)
		// OpenAPI normalizes all integer types to int
		assert.Equal(t, 50, param.Schema.Value.Default)
	})

	// Note: Float types are currently registered as integer parameters in RegisterParams
	// (see openapi.go:310-313), which is a pre-existing limitation.
	// Default tags for float parameters are not tested here due to this issue.

	t.Run("Bool default", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Active bool `query:"active" default:"true"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.NoError(t, err)

		param := route.Operation.Parameters.GetByInAndName("query", "active")
		require.NotNil(t, param)
		require.NotNil(t, param.Schema)
		require.NotNil(t, param.Schema.Value)
		assert.Equal(t, true, param.Schema.Value.Default)
	})

	t.Run("Int slice default", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Tags []int `query:"tags" default:"1,2,3"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.NoError(t, err)

		param := route.Operation.Parameters.GetByInAndName("query", "tags")
		require.NotNil(t, param)
		require.NotNil(t, param.Schema)
		require.NotNil(t, param.Schema.Value)
		// Array defaults are stored as actual arrays in OpenAPI
		assert.Equal(t, []any{1, 2, 3}, param.Schema.Value.Default)
	})

	t.Run("String slice default", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Categories []string `query:"categories" default:"a,b,c"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.NoError(t, err)

		param := route.Operation.Parameters.GetByInAndName("query", "categories")
		require.NotNil(t, param)
		require.NotNil(t, param.Schema)
		require.NotNil(t, param.Schema.Value)
		// Array defaults are stored as actual arrays in OpenAPI
		assert.Equal(t, []any{"a", "b", "c"}, param.Schema.Value.Default)
	})

	t.Run("Header with default", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Auth string `header:"Authorization" default:"Bearer token"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.NoError(t, err)

		param := route.Operation.Parameters.GetByInAndName("header", "Authorization")
		require.NotNil(t, param)
		require.NotNil(t, param.Schema)
		require.NotNil(t, param.Schema.Value)
		assert.Equal(t, "Bearer token", param.Schema.Value.Default)
	})

	t.Run("Cookie with default", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			SessionID string `cookie:"session_id" default:"guest"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.NoError(t, err)

		param := route.Operation.Parameters.GetByInAndName("cookie", "session_id")
		require.NotNil(t, param)
		require.NotNil(t, param.Schema)
		require.NotNil(t, param.Schema.Value)
		assert.Equal(t, "guest", param.Schema.Value.Default)
	})

	t.Run("Empty default value", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Name string `query:"name" default:""`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.NoError(t, err)

		param := route.Operation.Parameters.GetByInAndName("query", "name")
		require.NotNil(t, param)
		// Empty default should not be set
		if param.Schema != nil && param.Schema.Value != nil {
			assert.Nil(t, param.Schema.Value.Default)
		}
	})

	t.Run("Invalid int default", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Age int `query:"age" default:"not-a-number"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid default value for field Age")
	})

	t.Run("Invalid bool default", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Active bool `query:"active" default:"yes"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid default value for field Active")
	})

	t.Run("Invalid array element", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Tags []int `query:"tags" default:"1,abc,3"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid default value for field Tags")
	})

	t.Run("Multiple params with defaults", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, struct {
			Limit  int64  `query:"limit" default:"10"`
			Active bool   `query:"active" default:"true"`
			Name   string `query:"name" default:"test"`
		}](
			http.MethodGet,
			"/test",
			handler,
			s.Engine,
		)
		err := route.RegisterParams()
		require.NoError(t, err)

		limitParam := route.Operation.Parameters.GetByInAndName("query", "limit")
		require.NotNil(t, limitParam)
		// OpenAPI normalizes all integer types to int
		assert.Equal(t, 10, limitParam.Schema.Value.Default)

		activeParam := route.Operation.Parameters.GetByInAndName("query", "active")
		require.NotNil(t, activeParam)
		assert.Equal(t, true, activeParam.Schema.Value.Default)

		nameParam := route.Operation.Parameters.GetByInAndName("query", "name")
		require.NotNil(t, nameParam)
		assert.Equal(t, "test", nameParam.Schema.Value.Default)
	})
}
