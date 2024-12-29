package fuego

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTags(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", testController,
		OptionTags("my-tag"),
		OptionDescription("my description"),
		OptionSummary("my summary"),
		OptionDeprecated(),
	)

	require.Equal(t, []string{"my-tag"}, route.Operation.Tags)
	require.Equal(t, "#### Controller: \n\n`github.com/go-fuego/fuego.testController`\n\n#### Middlewares:\n\n- `github.com/go-fuego/fuego.defaultLogger.middleware`\n\n---\n\nmy description", route.Operation.Description)
	require.Equal(t, "my summary", route.Operation.Summary)
	require.Equal(t, true, route.Operation.Deprecated)
}

func TestAddTags(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", dummyController,
		OptionTags("my-tag"),
		OptionTags("my-other-tag"),
	)

	require.Equal(t, route.Operation.Tags, []string{"my-tag", "my-other-tag"})
}

func TestQuery(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", dummyController,
		OptionQuery("my-param", "my description"),
	)

	require.Equal(t, "my description", route.Operation.Parameters.GetByInAndName("query", "my-param").Description)
}

func TestHeaderParams(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", testController,
		OptionHeader("my-header", "my description"),
	)

	require.Equal(t, "my description", route.Operation.Parameters.GetByInAndName("header", "my-header").Description)
}

func TestCustomError(t *testing.T) {
	type MyError struct {
		Message string
	}
	s := NewServer()
	route := Get(s, "/test", testController,
		OptionAddError(400, "My Validation Error", MyError{}),
	)

	require.Equal(t, "My Validation Error", *route.Operation.Responses.Map()["400"].Value.Description)
}

func TestWithGlobalResponseType(t *testing.T) {
	type MyGlobalResponse struct {
		Message string
	}
	type MyLocalResponse struct {
		Message string
	}
	t.Run("base", func(t *testing.T) {
		s := NewServer(
			WithGlobalResponseTypes(http.StatusNotImplemented, "My Global Error", Response{Type: MyGlobalResponse{}}),
		)
		routeGlobal := Get(s, "/test-global", testController)
		require.Equal(t, "My Global Error", *routeGlobal.Operation.Responses.Value("501").Value.Description)
	})

	t.Run("base with custom contents", func(t *testing.T) {
		s := NewServer(
			WithGlobalResponseTypes(http.StatusNotImplemented, "My Global Error", Response{
				Type:         MyGlobalResponse{},
				ContentTypes: []string{"application/x-yaml"},
			}),
		)
		routeGlobal := Get(s, "/test-global", testController)
		require.NotNil(t, routeGlobal.Operation.Responses.Value("501").Value.Content.Get("application/x-yaml"))
		require.Nil(t, routeGlobal.Operation.Responses.Value("501").Value.Content.Get("application/xml"))
	})

	t.Run("errors with route overrides", func(t *testing.T) {
		s := NewServer(
			WithGlobalResponseTypes(http.StatusBadRequest, "My Global Error", Response{Type: MyGlobalResponse{}}),
			WithGlobalResponseTypes(http.StatusNotImplemented, "Another Global Error", Response{Type: MyGlobalResponse{}}),
		)

		routeGlobal := Get(s, "/test-global", testController)
		routeCustom := Get(s, "/test-custom", testController,
			OptionAddResponse(http.StatusBadRequest, "My Local Error", Response{Type: MyLocalResponse{}}),
			OptionAddResponse(http.StatusTeapot, "My Local Teapot", Response{Type: HTTPError{}}),
		)

		require.Equal(t, "My Global Error", *routeGlobal.Operation.Responses.Value("400").Value.Description, "Overrides Fuego's default 400 error")
		require.Equal(t, "Another Global Error", *routeGlobal.Operation.Responses.Value("501").Value.Description)

		require.Equal(t, "My Local Error", *routeCustom.Operation.Responses.Map()["400"].Value.Description, "Local error overrides global error")
		require.Equal(t, "My Local Teapot", *routeCustom.Operation.Responses.Map()["418"].Value.Description)
		require.Equal(t, "Internal Server Error _(panics)_", *routeCustom.Operation.Responses.Map()["500"].Value.Description, "Global error set by default by Fuego")
		require.Equal(t, "Another Global Error", *routeCustom.Operation.Responses.Map()["501"].Value.Description, "Global error is still available")
	})

	t.Run("200 responses with overrides", func(t *testing.T) {
		s := NewServer(
			WithGlobalResponseTypes(http.StatusCreated, "A Global Response", Response{Type: MyGlobalResponse{}}),
			WithGlobalResponseTypes(http.StatusAccepted, "My 202 response with content", Response{
				Type: MyGlobalResponse{}, ContentTypes: []string{"application/x-yaml"},
			}),
		)

		t.Run("routeGlobal", func(t *testing.T) {
			routeGlobal := Get(s, "/test-global", testController)
			require.Equal(t,
				"#/components/schemas/ans",
				routeGlobal.Operation.Responses.Value("200").Value.Content.Get("application/json").Schema.Ref,
			)
			require.Equal(t,
				"#/components/schemas/ans",
				routeGlobal.Operation.Responses.Value("200").Value.Content.Get("application/xml").Schema.Ref,
			)
			require.Equal(t, "A Global Response", *routeGlobal.Operation.Responses.Value("201").Value.Description)
			require.Equal(t, "My 202 response with content", *routeGlobal.Operation.Responses.Value("202").Value.Description)
			require.Equal(t,
				"#/components/schemas/MyGlobalResponse",
				routeGlobal.Operation.Responses.Value("202").Value.Content.Get("application/x-yaml").Schema.Ref,
			)
		})

		t.Run("routeCustom", func(t *testing.T) {
			routeCustom := Get(s, "/test-custom", testController,
				OptionAddResponse(http.StatusOK, "My Local Response", Response{Type: MyLocalResponse{}}),
				OptionAddResponse(http.StatusNoContent, "My No Content", Response{Type: struct{}{}}),
			)
			require.Equal(t,
				"#/components/schemas/MyLocalResponse",
				routeCustom.Operation.Responses.Value("200").Value.Content.Get("application/json").Schema.Ref,
			)
			require.Equal(t,
				"#/components/schemas/MyLocalResponse",
				routeCustom.Operation.Responses.Value("200").Value.Content.Get("application/xml").Schema.Ref,
			)
			require.Equal(t, "My No Content", *routeCustom.Operation.Responses.Value("204").Value.Description)
			require.Equal(t, "My 202 response with content", *routeCustom.Operation.Responses.Value("202").Value.Description)
			require.Equal(t,
				"#/components/schemas/MyGlobalResponse",
				routeCustom.Operation.Responses.Value("202").Value.Content.Get("application/x-yaml").Schema.Ref,
			)
		})
	})

	t.Run("should be fatal", func(t *testing.T) {
		s := NewServer(
			WithGlobalResponseTypes(http.StatusNotImplemented, "My Global Error", Response{}),
		)
		require.Panics(t, func() {
			Get(s, "/test-global", testController)
		})
	})
}

func TestCookieParams(t *testing.T) {
	t.Run("basic cookie", func(t *testing.T) {
		s := NewServer()
		route := Get(s, "/test", testController,
			OptionCookie("my-cookie", "my description"),
		)

		require.Equal(t, "my description", route.Operation.Parameters.GetByInAndName("cookie", "my-cookie").Description)
	})

	t.Run("with more parameters", func(t *testing.T) {
		s := NewServer()
		route := Get(s, "/test", testController,
			OptionCookie("my-cookie", "my description", ParamRequired(), ParamExample("example", "my-example")),
		)

		cookieParam := route.Operation.Parameters.GetByInAndName("cookie", "my-cookie")
		t.Logf("%#v", cookieParam.Examples["example"].Value)
		require.Equal(t, "my description", cookieParam.Description)
		require.Equal(t, true, cookieParam.Required)
		require.Equal(t, "my-example", cookieParam.Examples["example"].Value.Value)
	})
}
