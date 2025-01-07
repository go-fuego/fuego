package fuego

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithErrorHandler(t *testing.T) {
	t.Run("default engine", func(t *testing.T) {
		e := NewEngine()
		err := NotFoundError{
			Err: errors.New("Not Found :c"),
		}
		errResponse := e.ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
	})
	t.Run("custom handler", func(t *testing.T) {
		e := NewEngine(
			WithErrorHandler(func(err error) error {
				return fmt.Errorf("%w foobar", err)
			}),
		)
		err := NotFoundError{
			Err: errors.New("Not Found :c"),
		}
		errResponse := e.ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.ErrorContains(t, errResponse, "Not Found :c foobar")
	})
	t.Run("should be fatal", func(t *testing.T) {
		require.Panics(t, func() {
			NewEngine(
				WithErrorHandler(nil),
			)
		})
	})

	t.Run("disable error handler", func(t *testing.T) {
		e := NewEngine(DisableErrorHandler())
		err := NotFoundError{
			Err: errors.New("Not Found"),
		}
		errResponse := e.ErrorHandler(err)
		require.Equal(t, "Not Found", errResponse.Error())
	})
}

func TestWithRequestContentType(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		e := NewEngine()
		require.Nil(t, e.acceptedContentTypes)
	})

	t.Run("input", func(t *testing.T) {
		arr := []string{"application/json", "application/xml"}
		e := NewEngine(WithRequestContentType("application/json", "application/xml"))
		require.ElementsMatch(t, arr, e.acceptedContentTypes)
	})

	t.Run("ensure applied to route", func(t *testing.T) {
		s := NewServer(WithEngineOptions(
			WithRequestContentType("application/json", "application/xml")),
		)
		route := Post(s, "/test", dummyController)

		content := route.Operation.RequestBody.Value.Content
		require.NotNil(t, content.Get("application/json"))
		require.NotNil(t, content.Get("application/xml"))
		require.Equal(t, "#/components/schemas/ReqBody", content.Get("application/json").Schema.Ref)
		require.Equal(t, "#/components/schemas/ReqBody", content.Get("application/xml").Schema.Ref)
		_, ok := s.OpenAPI.Description().Components.RequestBodies["ReqBody"]
		require.False(t, ok)
	})
}
