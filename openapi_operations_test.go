package fuego

import (
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
	require.Equal(t, "controller: `github.com/go-fuego/fuego.testController`\n\n---\n\nmy description", route.Operation.Description)
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
	route := Get(s, "/test", testController).
		AddError(400, "My Validation Error", MyError{})

	require.Equal(t, "My Validation Error", *route.Operation.Responses.Map()["400"].Value.Description)
}

func TestCustomErrorGlobalAndOnRoute(t *testing.T) {
	type MyGlobalError struct {
		Message string
	}
	s := NewServer(
		WithGlobalResponseTypes(400, "My Global Error", MyGlobalError{}),
		WithGlobalResponseTypes(501, "Another Global Error", MyGlobalError{}),
	)

	type MyLocalError struct {
		Message string
	}

	routeGlobal := Get(s, "/test-global", testController)
	routeCustom := Get(s, "/test-custom", testController).
		AddError(400, "My Local Error", MyLocalError{}).
		AddError(419, "My Local Teapot")

	require.Equal(t, "My Global Error", *routeGlobal.Operation.Responses.Map()["400"].Value.Description, "Overrides Fuego's default 400 error")
	require.Equal(t, "Another Global Error", *routeGlobal.Operation.Responses.Map()["501"].Value.Description)

	require.Equal(t, "My Local Error", *routeCustom.Operation.Responses.Map()["400"].Value.Description, "Local error overrides global error")
	require.Equal(t, "My Local Teapot", *routeCustom.Operation.Responses.Map()["419"].Value.Description)
	require.Equal(t, "Internal Server Error _(panics)_", *routeCustom.Operation.Responses.Map()["500"].Value.Description, "Global error set by default by Fuego")
	require.Equal(t, "Another Global Error", *routeCustom.Operation.Responses.Map()["501"].Value.Description, "Global error is still available")
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
