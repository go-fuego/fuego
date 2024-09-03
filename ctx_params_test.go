package fuego_test

import (
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
	"github.com/stretchr/testify/require"
)

func TestParam(t *testing.T) {
	t.Run("Query params default values", func(t *testing.T) {
		s := fuego.NewServer()

		fuego.Get(s, "/test", func(c fuego.ContextNoBody) (string, error) {
			name := c.QueryParam("name")
			age := c.QueryParamInt("age")
			isok := c.QueryParamBool("is_ok")

			return name + strconv.Itoa(age) + strconv.FormatBool(isok), nil
		},
			option.Query("name", "Name", param.Required(), param.Default("hey"), param.Example("example1", "you")),
			option.QueryInt("age", "Age", param.Nullable(), param.Default(18), param.Example("example1", 1)),
			option.QueryBool("is_ok", "Is OK?", param.Default(true), param.Example("example1", true)),
		)

		t.Run("Default should correctly set parameter in controller", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			s.Mux.ServeHTTP(w, r)

			require.Equal(t, http.StatusOK, w.Code)
			require.Equal(t, "hey18true", w.Body.String())
		})
	})
}
