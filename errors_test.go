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

var _ ErrorWithStatus = myError{}

func (e myError) Error() string   { return "test error" }
func (e myError) StatusCode() int { return e.status }

func TestErrorHandler(t *testing.T) {
	t.Run("basic error", func(t *testing.T) {
		err := errors.New("test error")

		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.Contains(t, errResponse.Error(), "Internal Server Error")
		require.Equal(t, http.StatusInternalServerError, errResponse.(HTTPError).StatusCode())
	})

	t.Run("not found error", func(t *testing.T) {
		err := NotFoundError{}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.Contains(t, errResponse.Error(), "Not Found")
		require.Contains(t, errResponse.Error(), "404")
		require.Equal(t, http.StatusNotFound, errResponse.(HTTPError).StatusCode())
	})

	t.Run("error with status ", func(t *testing.T) {
		err := myError{
			status: http.StatusNotFound,
		}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.Contains(t, errResponse.Error(), "Not Found")
		require.Contains(t, errResponse.Error(), "404")
		require.Equal(t, http.StatusNotFound, errResponse.(HTTPError).StatusCode())
	})
}
