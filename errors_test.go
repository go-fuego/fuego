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
	t.Run("basic", func(t *testing.T) {
		err := errors.New("test error")

		errResponse := ErrorHandler(err)
		require.Contains(t, errResponse.Error(), "test error")
	})

	t.Run("not found error", func(t *testing.T) {
		err := NotFoundError{
			Err: errors.New("Not Found :c"),
		}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.Contains(t, err.Error(), "Not Found :c")
		require.Contains(t, errResponse.Error(), "Not Found")
		require.Contains(t, errResponse.Error(), "404")
		require.Equal(t, http.StatusNotFound, errResponse.(HTTPError).StatusCode())
	})

	t.Run("not duplicate HTTPError", func(t *testing.T) {
		err := HTTPError{
			Err: errors.New("HTTPError"),
		}
		errResponse := ErrorHandler(err)

		var httpError HTTPError
		require.ErrorAs(t, errResponse, &httpError)
		require.False(t, errors.As(httpError.Err, &HTTPError{}))
		require.Contains(t, err.Error(), "Internal Server Error")
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

	t.Run("conflict error", func(t *testing.T) {
		err := ConflictError{
			Err: errors.New("Conflict"),
		}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.Contains(t, err.Error(), "Conflict")
		require.Contains(t, errResponse.Error(), "Conflict")
		require.Contains(t, errResponse.Error(), "409")
		require.Equal(t, http.StatusConflict, errResponse.(HTTPError).StatusCode())
	})

	t.Run("unauthorized error", func(t *testing.T) {
		err := UnauthorizedError{
			Err: errors.New("coucou"),
		}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.Contains(t, err.Error(), "coucou")
		require.Contains(t, errResponse.Error(), "Unauthorized")
		require.Contains(t, errResponse.Error(), "401")
		require.Equal(t, http.StatusUnauthorized, errResponse.(HTTPError).StatusCode())
	})

	t.Run("forbidden error", func(t *testing.T) {
		err := ForbiddenError{
			Err: errors.New("Forbidden"),
		}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.Contains(t, err.Error(), "Forbidden")
		require.Contains(t, errResponse.Error(), "Forbidden")
		require.Contains(t, errResponse.Error(), "403")
		require.Equal(t, http.StatusForbidden, errResponse.(HTTPError).StatusCode())
	})
}

func TestHTTPError_Error(t *testing.T) {
	t.Run("title", func(t *testing.T) {
		t.Run("custom title", func(t *testing.T) {
			err := HTTPError{
				Title: "Custom Title",
			}
			require.Contains(t, err.Error(), "Custom Title")
		})
		t.Run("title from status", func(t *testing.T) {
			err := HTTPError{Status: http.StatusNotFound}
			require.Contains(t, err.Error(), "Not Found")
		})
		t.Run("default title", func(t *testing.T) {
			err := HTTPError{}
			require.Contains(t, err.Error(), "Internal Server Error")
		})
	})
}

func TestHTTPError_Unwrap(t *testing.T) {
	err := myError{status: 999}

	errResponse := HTTPError{
		Err: err,
	}

	var unwrapped myError
	require.True(t, errors.As(errResponse.Unwrap(), &unwrapped))
	require.Equal(t, 999, unwrapped.status)
}
