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

	t.Run("RegisterParams should not with interfaces", func(t *testing.T) {
		route := NewRoute[struct{}, struct{}, any](
			http.MethodGet,
			"/no-interfaces",
			handler,
			s.Engine,
			OptionDefaultStatusCode(201),
		)

		err := route.RegisterParams()
		require.Error(t, err)
	})
}
