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

	require.Equal(t, route.operation.Tags, []string{"my-tag"})
	require.Equal(t, route.operation.Description, "my description")
	require.Equal(t, route.operation.Summary, "my summary")
	require.Equal(t, route.operation.Deprecated, true)
}

func TestAddTags(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	}).
		AddTags("my-tag").
		AddTags("my-other-tag")

	require.Equal(t, route.operation.Tags, []string{"string", "my-tag", "my-other-tag"})
}

func TestRemoveTags(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	}).
		AddTags("my-tag").
		RemoveTags("my-tag", "string").
		AddTags("my-other-tag")

	require.Equal(t, route.operation.Tags, []string{"my-other-tag"})
}

func TestQueryParams(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", func(ctx *ContextNoBody) (string, error) {
		return "test", nil
	}).
		QueryParam("my-param", "my description")

	require.Equal(t, "my description", route.operation.Parameters[0].Description)
}

func TestHeaderParams(t *testing.T) {
	s := NewServer()
	route := Get(s, "/test", testController).
		Header("my-header", "my description")

	require.Equal(t, "my description", route.operation.Parameters[0].Description)
}

func TestCookieParams(t *testing.T) {
	t.Run("basic cookie", func(t *testing.T) {
		s := NewServer()
		route := Get(s, "/test", testController).
			Cookie("my-cookie", "my description")

		require.Equal(t, "my description", route.operation.Parameters[0].Description)
	})

	t.Run("with more parameters", func(t *testing.T) {
		s := NewServer()
		route := Get(s, "/test", testController).
			Cookie("my-cookie", "my description", OpenAPIParam{Required: true, Example: "my-example"})

		require.Equal(t, "my description", route.operation.Parameters[0].Description)
		require.Equal(t, true, route.operation.Parameters[0].Required)
		require.Equal(t, "my-example", route.operation.Parameters[0].Example)
	})
}
