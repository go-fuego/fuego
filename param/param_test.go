package param_test

import (
	"strconv"
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
	"github.com/stretchr/testify/require"
)

func TestParam(t *testing.T) {
	t.Run("All options", func(t *testing.T) {
		s := fuego.NewServer()

		route := fuego.Get(s, "/test", func(c fuego.ContextNoBody) (string, error) {
			name := c.QueryParam("name")
			age := c.QueryParamInt("age")
			isok := c.QueryParamBool("is_ok")

			return name + strconv.Itoa(age) + strconv.FormatBool(isok), nil
		},
			option.Query("name", "Name", param.Required(), param.Default("hey"), param.Example("example1", "you")),
			option.QueryInt("age", "Age", param.Nullable(), param.Default(18), param.Example("example1", 1)),
			option.QueryBool("is_ok", "Is OK?", param.Default(true), param.Example("example1", true)),
		)

		require.NotNil(t, route)
		require.NotNil(t, route.Params)
		require.Len(t, route.Params, 3)
		require.Equal(t, "Name", route.Params["name"].Description)
		require.True(t, route.Params["name"].Required)
		require.Equal(t, "hey", route.Params["name"].Default)
		require.Equal(t, "you", route.Params["name"].Examples["example1"])
		require.Equal(t, "string", route.Params["name"].GoType)

		require.Equal(t, "Age", route.Params["age"].Description)
		require.True(t, route.Params["age"].Nullable)
		require.Equal(t, 18, route.Params["age"].Default)
		require.Equal(t, "integer", route.Params["age"].GoType)
	})
}
