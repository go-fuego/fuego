package fuego

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type myError struct {
	status int
}

func (e myError) Error() string { return "test error" }
func (e myError) Status() int   { return e.status }

func TestErrorHandler(t *testing.T) {
	t.Run("basic error", func(t *testing.T) {
		err := errors.New("test error")

		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.Equal(t, "test error", errResponse.Error())
		require.Equal(t, http.StatusInternalServerError, errResponse.(HTTPError).Status())
		require.Nil(t, errResponse.(HTTPError).Info())
	})

	t.Run("error with status ", func(t *testing.T) {
		err := myError{
			status: http.StatusNotFound,
		}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.Equal(t, "test error", errResponse.Error())
		require.Equal(t, http.StatusNotFound, errResponse.(HTTPError).Status())
		require.Nil(t, errResponse.(HTTPError).Info())
	})

	t.Run("error with status and info", func(t *testing.T) {
		err := HTTPError{
			Message:    "test error",
			StatusCode: http.StatusNotFound,
			MoreInfo: map[string]any{
				"test": "info",
			},
		}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.Equal(t, "test error", errResponse.Error())
		require.Equal(t, http.StatusNotFound, errResponse.(HTTPError).Status())
		require.NotNil(t, errResponse.(HTTPError).Info())
	})
}
