package fuego

import (
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thejerf/slogassert"
)

func Test_RegisterOpenAPIOperation(t *testing.T) {
    openapi := NewOpenAPI()
    handler := slogassert.New(t, slog.LevelWarn, nil)
    s := NewServer()
    t.Run("Register with params", func(t *testing.T) {
        route := NewRoute[struct {
            QueryParam  string `query:"queryParam"`
            HeaderParam string `header:"headerParam"`
            PathParam   string `path:"pathParam"`
        }, struct{}](
            http.MethodGet,
            "/some/path/{pathParam}",
            handler,
            s.Engine,
        )

        operation, err := RegisterOpenAPIOperation(openapi, route)
        require.NoError(t, err)
        assert.NotNil(t, operation)
        assert.Equal(t, route.Method+"_"+"/some/path/:pathParam", operation.OperationID)
        assert.Equal(t, http.MethodGet, strings.ToUpper(route.Method))
        assert.Len(t, operation.Parameters, 3)
        assert.Equal(t, "queryParam", operation.Parameters.GetByInAndName("query", "queryParam").Name)
        assert.Equal(t, "headerParam", operation.Parameters.GetByInAndName("header", "headerParam").Name)
        defaultResponse := operation.Responses
        assert.NotNil(t, defaultResponse)
    })
    t.Run("Automatically add path parameters", func(t *testing.T) {
        route := NewRoute[struct {
            QueryParam  string `query:"queryParam"`
            HeaderParam string `header:"headerParam"`
        }, struct{}](
            "GET",
            "/some/path/{pathParam}",
            handler,
            s.Engine,
        )
        operation, err := RegisterOpenAPIOperation(openapi, route)
        require.NoError(t, err)
        assert.Len(t, operation.Parameters, 3) // Should add the path parameter with others
        assert.Equal(t, "pathParam", operation.Parameters[2].Value.Name)
        assert.Equal(t, "path", operation.Parameters[2].Value.In)
    })
}