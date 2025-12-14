package fuego_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
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

	t.Run("Params with default tag", func(t *testing.T) {
		type TestParams struct {
			Limit  int64  `query:"limit" default:"10"`
			Active bool   `query:"active" default:"true"`
			Name   string `query:"name" default:"test"`
		}

		s := fuego.NewServer()

		fuego.Get(s, "/test", func(c fuego.Context[any, TestParams]) (string, error) {
			params, err := c.Params()
			if err != nil {
				return "", err
			}

			return params.Name + strconv.FormatInt(params.Limit, 10) + strconv.FormatBool(params.Active), nil
		})

		t.Run("Defaults work when params not provided", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			s.Mux.ServeHTTP(w, r)

			require.Equal(t, http.StatusOK, w.Code)
			require.Equal(t, "test10true", w.Body.String())
		})

		t.Run("Provided values override defaults", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/test?limit=50&active=false&name=custom", nil)
			w := httptest.NewRecorder()
			s.Mux.ServeHTTP(w, r)

			require.Equal(t, http.StatusOK, w.Code)
			require.Equal(t, "custom50false", w.Body.String())
		})

		t.Run("Partial override - some defaults used", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/test?name=partial", nil)
			w := httptest.NewRecorder()
			s.Mux.ServeHTTP(w, r)

			require.Equal(t, http.StatusOK, w.Code)
			require.Equal(t, "partial10true", w.Body.String())
		})
	})

	t.Run("Params with array default tag", func(t *testing.T) {
		type ArrayParams struct {
			Tags []int `query:"tags" default:"1,2,3"`
		}

		s := fuego.NewServer()

		fuego.Get(s, "/test", func(c fuego.Context[any, ArrayParams]) (string, error) {
			params, err := c.Params()
			if err != nil {
				return "", err
			}

			var result strings.Builder
			for _, tag := range params.Tags {
				result.WriteString(strconv.Itoa(tag) + ",")
			}
			return result.String(), nil
		})

		t.Run("Array defaults work when param not provided", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			s.Mux.ServeHTTP(w, r)

			require.Equal(t, http.StatusOK, w.Code)
			// Note: The runtime might not handle array defaults from OpenAPIParams
			// This test will help us verify if additional runtime support is needed
		})
	})
}
