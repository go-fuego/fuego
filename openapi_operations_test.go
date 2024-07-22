package fuego

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTags(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", testController).
		Tags("my-tag").
		Description("my description").
		Summary("my summary").
		Deprecated()

	require.Equal(t, route.Operation.Tags, []string{"my-tag"})
	require.Equal(t, route.Operation.Description, "my description")
	require.Equal(t, route.Operation.Summary, "my summary")
	require.Equal(t, route.Operation.Deprecated, true)
}

func TestAddTags(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	}).
		AddTags("my-tag").
		AddTags("my-other-tag")

	require.Equal(t, route.Operation.Tags, []string{"my-tag", "my-other-tag"})
}

func TestRemoveTags(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	}).
		AddTags("my-tag").
		RemoveTags("my-tag", "string").
		AddTags("my-other-tag")

	require.Equal(t, route.Operation.Tags, []string{"my-other-tag"})
}

func TestQueryParams(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	}).
		QueryParam("my-param", "my description")

	require.Equal(t, "my description", route.Operation.Parameters.GetByInAndName("query", "my-param").Description)
}

func TestHeaderParams(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", testController).
		Header("my-header", "my description")

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

func TestCookieParams(t *testing.T) {
	t.Run("basic cookie", func(t *testing.T) {
		s := NewServer()
		route := Get(s, "/test", testController).
			Cookie("my-cookie", "my description")

		require.Equal(t, "my description", route.Operation.Parameters.GetByInAndName("cookie", "my-cookie").Description)
	})

	t.Run("with more parameters", func(t *testing.T) {
		s := NewServer()
		route := Get(s, "/test", testController).
			Cookie("my-cookie", "my description", OpenAPIParam{Required: true, Example: "my-example"})

		require.Equal(t, "my description", route.Operation.Parameters.GetByInAndName("cookie", "my-cookie").Description)
		require.Equal(t, true, route.Operation.Parameters.GetByInAndName("cookie", "my-cookie").Required)
		require.Equal(t, "my-example", route.Operation.Parameters.GetByInAndName("cookie", "my-cookie").Example)
	})
}
