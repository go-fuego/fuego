package fuego

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func controller(c Ctx[any]) (testStruct, error) {
	return testStruct{"Ewen", 23}, nil
}

func controllerWithError(c Ctx[any]) (testStruct, error) {
	return testStruct{}, errors.New("error")
}

func TestNewServer(t *testing.T) {
	s := NewServer()

	t.Run("can register controller", func(t *testing.T) {
		Get(s, "/", controller)

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		s.Mux.ServeHTTP(recorder, req)

		require.Equal(t, 200, recorder.Code)
	})
}

func TestWithXML(t *testing.T) {
	s := NewServer(
		WithXML(),
	)
	Get(s, "/", controller)
	Get(s, "/error", controllerWithError)

	t.Run("response is XML", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		s.Mux.ServeHTTP(recorder, req)

		require.Equal(t, 200, recorder.Code)
		require.Equal(t, "application/xml", recorder.Header().Get("Content-Type"))
		require.Equal(t, "<testStruct><Name>Ewen</Name><Age>23</Age></testStruct>", recorder.Body.String())
	})

	t.Run("error response is XML", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/error", nil)

		s.Mux.ServeHTTP(recorder, req)

		require.Equal(t, 500, recorder.Code)
		require.Equal(t, "application/xml", recorder.Header().Get("Content-Type"))
		require.Equal(t, "<HTTPError><Error>error</Error></HTTPError>", recorder.Body.String())
	})
}

func TestWithOpenAPIConfig(t *testing.T) {
	s := NewServer(
		WithOpenapiConfig(OpenapiConfig{}),
	)

	require.Equal(t, "/swagger", s.OpenapiConfig.SwaggerUrl)
	require.Equal(t, "/swagger/openapi.json", s.OpenapiConfig.JsonSpecUrl)
	require.Equal(t, "doc/openapi.json", s.OpenapiConfig.JsonSpecLocalPath)
}

func TestWithBasePath(t *testing.T) {
	s := NewServer(
		WithBasePath("/api"),
	)

	require.Equal(t, "/api", s.basePath)
}

func TestWithMaxBodySize(t *testing.T) {
	s := NewServer(
		WithMaxBodySize(1024),
	)

	require.Equal(t, int64(1024), s.maxBodySize)
}

func TestWithAutoAuth(t *testing.T) {
	s := NewServer(
		WithAutoAuth(nil),
	)

	require.NotNil(t, s.autoAuth)
	require.True(t, s.autoAuth.Enabled)
	// The authoauth is tested in security_test.go,
	// this is just an option to enable it.
}
