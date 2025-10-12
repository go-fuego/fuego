package fuego

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thejerf/slogassert"
)

type myError struct {
	status int
	err    HTTPError
	detail string
}

var (
	_ ErrorWithStatus = myError{}
	_ ErrorWithDetail = myError{}
)

func (e myError) Error() string     { return "test error" }
func (e myError) StatusCode() int   { return e.status }
func (e myError) DetailMsg() string { return e.detail }
func (e myError) Unwrap() error     { return e.err }

func TestErrorHandler(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		err := errors.New("test error")

		handler := slogassert.NewDefault(t)

		errResponse := ErrorHandler(context.Background(), err)
		require.ErrorContains(t, errResponse, "test error")

		handler.AssertMessage("Error in controller")

		handler.AssertEmpty()
	})

	t.Run("not found error", func(t *testing.T) {
		err := NotFoundError{
			Err: errors.New("Not Found :c"),
		}
		var errHTTP HTTPError
		require.ErrorAs(t, ErrorHandler(context.Background(), err), &errHTTP)
		require.ErrorContains(t, err, "Not Found :c")
		require.ErrorContains(t, errHTTP, "Not Found")
		require.ErrorContains(t, errHTTP, "404")
		assert.Equal(t, http.StatusNotFound, errHTTP.StatusCode())
	})

	t.Run("not duplicate HTTPError", func(t *testing.T) {
		err := HTTPError{
			Err: errors.New("HTTPError"),
		}

		var errHTTP HTTPError
		require.ErrorAs(t, ErrorHandler(context.Background(), err), &errHTTP)
		require.NotErrorAs(t, errHTTP.Err, &HTTPError{})
		require.ErrorContains(t, err, "Internal Server Error")
	})

	t.Run("error with status", func(t *testing.T) {
		err := myError{
			status: http.StatusNotFound,
		}
		var errHTTP HTTPError
		require.ErrorAs(t, ErrorHandler(context.Background(), err), &errHTTP)
		require.ErrorContains(t, errHTTP, "Not Found")
		require.ErrorContains(t, errHTTP, "404")
		require.Equal(t, http.StatusNotFound, errHTTP.StatusCode())
	})

	t.Run("error with detail", func(t *testing.T) {
		err := myError{
			detail: "my detail",
		}
		var errHTTP HTTPError
		require.ErrorAs(t, ErrorHandler(context.Background(), err), &errHTTP)
		require.Contains(t, errHTTP.Error(), "Internal Server Error")
		require.Contains(t, errHTTP.Error(), "500")
		require.Contains(t, errHTTP.Error(), "my detail")
		require.Equal(t, http.StatusInternalServerError, errHTTP.StatusCode())
	})

	t.Run("conflict error", func(t *testing.T) {
		err := ConflictError{
			Err: errors.New("Conflict"),
		}
		var errHTTP HTTPError
		require.ErrorAs(t, ErrorHandler(context.Background(), err), &errHTTP)
		require.ErrorContains(t, err, "Conflict")
		require.ErrorContains(t, errHTTP, "Conflict")
		require.ErrorContains(t, errHTTP, "409")
		require.Equal(t, http.StatusConflict, errHTTP.StatusCode())
	})

	t.Run("unauthorized error", func(t *testing.T) {
		err := UnauthorizedError{
			Err: errors.New("coucou"),
		}
		var errHTTP HTTPError
		require.ErrorAs(t, ErrorHandler(context.Background(), err), &errHTTP)
		require.ErrorContains(t, err, "coucou")
		require.ErrorContains(t, errHTTP, "Unauthorized")
		require.ErrorContains(t, errHTTP, "401")
		require.Equal(t, http.StatusUnauthorized, errHTTP.StatusCode())
	})

	t.Run("forbidden error", func(t *testing.T) {
		err := ForbiddenError{
			Err: errors.New("Forbidden"),
		}
		var errHTTP HTTPError
		require.ErrorAs(t, ErrorHandler(context.Background(), err), &errHTTP)
		require.ErrorContains(t, err, "Forbidden")
		require.ErrorContains(t, errHTTP, "Forbidden")
		require.ErrorContains(t, errHTTP, "403")
		require.Equal(t, http.StatusForbidden, errHTTP.StatusCode())
	})
}

func TestHandleHTTPError(t *testing.T) {
	t.Run("should always be HTTPError", func(t *testing.T) {
		err := errors.New("test error")

		errResponse := HandleHTTPError(context.Background(), err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.ErrorContains(t, errResponse, "500 Internal Server Error")
	})

	t.Run("not found error", func(t *testing.T) {
		err := NotFoundError{
			Err: errors.New("Not Found :c"),
		}
		errResponse := HandleHTTPError(context.Background(), err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.ErrorContains(t, err, "Not Found :c")
		require.ErrorContains(t, errResponse, "Not Found")
		require.ErrorContains(t, errResponse, "404")
		require.Equal(t, http.StatusNotFound, errResponse.(HTTPError).StatusCode())
	})

	t.Run("error is a reference to HTTPError", func(t *testing.T) {
		err := &HTTPError{
			Title:  "HTTPError",
			Errors: []ErrorItem{{Name: "my name"}},
			Err:    errors.New("my new error"),
		}
		errResponse := HandleHTTPError(context.Background(), err)
		require.ErrorAs(t, errResponse, &HTTPError{})
		require.ErrorContains(t, err, "my new error")
		require.ErrorContains(t, errResponse, "HTTPError")
		require.ErrorContains(t, errResponse, "500")
		require.Equal(t, http.StatusInternalServerError, errResponse.(HTTPError).StatusCode())
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

func TestUnauthorizedError(t *testing.T) {
	t.Run("without error", func(t *testing.T) {
		err := UnauthorizedError{Title: "Unauthorized"}
		assert.EqualError(t, err, "401 Unauthorized")
	})

	t.Run("with error", func(t *testing.T) {
		err := UnauthorizedError{Title: "Unauthorized", Err: errors.New("error message")}
		assert.EqualError(t, err, "401 Unauthorized: error message")
	})

	t.Run("with error and detail", func(t *testing.T) {
		err := UnauthorizedError{Title: "Unauthorized", Detail: "detail message", Err: errors.New("error message")}
		assert.EqualError(t, err, "401 Unauthorized (detail message): error message")
	})
}

func TestForbiddenError(t *testing.T) {
	t.Run("without error", func(t *testing.T) {
		err := ForbiddenError{Title: "Access forbidden"}
		assert.EqualError(t, err, "403 Access forbidden")
	})

	t.Run("with error", func(t *testing.T) {
		err := ForbiddenError{Title: "Access forbidden", Err: errors.New("error message")}
		assert.EqualError(t, err, "403 Access forbidden: error message")
	})
}

func BenchmarkHTTPError_PublicError(b *testing.B) {
	err := HTTPError{
		Title:  "Custom Title",
		Detail: "Custom Detail",
		Status: http.StatusNotFound,
	}

	for range b.N {
		_ = err.PublicError()
	}
}
