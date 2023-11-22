package fuego

import (
	"errors"
	"log/slog"
	"net/http"
)

// ErrorWithStatus is an interface that can be implemented by an error to provide
// additional information about the error.
type ErrorWithStatus interface {
	error
	Status() int
}

// ErrorWithInfo is an interface that can be implemented by an error to provide
// additional information about the error.
type ErrorWithInfo interface {
	error
	Info() map[string]any
}

// HTTPError is the error response used by the serialization part of the framework.
type HTTPError struct {
	Message    string         `json:"error" xml:"Error"`                   // human readable error message
	StatusCode int            `json:"-" xml:"-"`                           // http status code
	MoreInfo   map[string]any `json:"info,omitempty" xml:"Info,omitempty"` // additional info
}

func (e HTTPError) Error() string {
	return e.Message
}

var _ ErrorWithStatus = HTTPError{}

func (e HTTPError) Status() int {
	if e.StatusCode == 0 {
		return http.StatusInternalServerError
	}
	return e.StatusCode
}

var _ ErrorWithInfo = HTTPError{}

func (e HTTPError) Info() map[string]any {
	return e.MoreInfo
}

// BadRequestError is an error used to return a 400 status code.
type BadRequestError struct {
	Message string
}

func (e BadRequestError) Error() string {
	return e.Message
}

func (e BadRequestError) Status() int {
	return http.StatusBadRequest
}

// ErrorHandler is the default error handler used by the framework.
// It transforms any error into the unified error type [HTTPError],
// Using the [ErrorWithStatus] and [ErrorWithInfo] interfaces.
func ErrorHandler(err error) error {
	errResponse := HTTPError{
		Message: err.Error(),
	}

	errResponse.StatusCode = http.StatusInternalServerError
	var errorStatus ErrorWithStatus
	if errors.As(err, &errorStatus) {
		errResponse.StatusCode = errorStatus.Status()
	}

	var errorInfo ErrorWithInfo
	if errors.As(err, &errorInfo) {
		errResponse.MoreInfo = errorInfo.Info()
	}

	slog.Error("Error : "+errResponse.Message, "status:", errResponse.StatusCode, "info:", errResponse.MoreInfo)

	return errResponse
}
