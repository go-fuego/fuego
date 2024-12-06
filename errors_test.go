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
		require.ErrorContains(t, errResponse, "test error")
	})

	t.Run("not found error", func(t *testing.T) {
		err := NotFoundError{
			Err: errors.New("Not Found :c"),
		}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.ErrorContains(t, err, "Not Found :c")
		require.ErrorContains(t, errResponse, "Not Found")
		require.ErrorContains(t, errResponse, "404")
		require.Equal(t, http.StatusNotFound, errResponse.(HTTPError).StatusCode())
	})

	t.Run("not duplicate HTTPError", func(t *testing.T) {
		err := HTTPError{
			Err: errors.New("HTTPError"),
		}
		errResponse := ErrorHandler(err)

		var httpError HTTPError
		require.ErrorAs(t, errResponse, &httpError)
		require.NotErrorAs(t, httpError.Err, &HTTPError{})
		require.ErrorContains(t, err, "Internal Server Error")
	})

	t.Run("error with status ", func(t *testing.T) {
		err := myError{
			status: http.StatusNotFound,
		}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.ErrorContains(t, errResponse, "Not Found")
		require.ErrorContains(t, errResponse, "404")
		require.Equal(t, http.StatusNotFound, errResponse.(HTTPError).StatusCode())
	})

	t.Run("conflict error", func(t *testing.T) {
		err := ConflictError{
			Err: errors.New("Conflict"),
		}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.ErrorContains(t, err, "Conflict")
		require.ErrorContains(t, errResponse, "Conflict")
		require.ErrorContains(t, errResponse, "409")
		require.Equal(t, http.StatusConflict, errResponse.(HTTPError).StatusCode())
	})

	t.Run("unauthorized error", func(t *testing.T) {
		err := UnauthorizedError{
			Err: errors.New("coucou"),
		}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.ErrorContains(t, err, "coucou")
		require.ErrorContains(t, errResponse, "Unauthorized")
		require.ErrorContains(t, errResponse, "401")
		require.Equal(t, http.StatusUnauthorized, errResponse.(HTTPError).StatusCode())
	})

	t.Run("forbidden error", func(t *testing.T) {
		err := ForbiddenError{
			Err: errors.New("Forbidden"),
		}
		errResponse := ErrorHandler(err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.ErrorContains(t, err, "Forbidden")
		require.ErrorContains(t, errResponse, "Forbidden")
		require.ErrorContains(t, errResponse, "403")
		require.Equal(t, http.StatusForbidden, errResponse.(HTTPError).StatusCode())
	})
}

func TestHTTPError_Error(t *testing.T) {
	t.Run("title", func(t *testing.T) {
		t.Run("custom title", func(t *testing.T) {
			err := HTTPError{
				Title: "Custom Title",
			}
			require.ErrorContains(t, err, "Custom Title")
		})
		t.Run("title from status", func(t *testing.T) {
			err := HTTPError{Status: http.StatusNotFound}
			require.ErrorContains(t, err, "Not Found")
		})
		t.Run("default title", func(t *testing.T) {
			err := HTTPError{}
			require.ErrorContains(t, err, "Internal Server Error")
		})
	})
}

func TestHTTPError_Unwrap(t *testing.T) {
	err := myError{status: 999}

	errResponse := HTTPError{
		Err: err,
	}

	var unwrapped myError
	require.ErrorAs(t, errResponse.Unwrap(), &unwrapped)
	require.Equal(t, 999, unwrapped.status)
}
