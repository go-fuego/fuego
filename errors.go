package fuego

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

// ErrorWithStatus is an interface that can be implemented by an error to provide
// a status code
type ErrorWithStatus interface {
	error
	StatusCode() int
}

// ErrorWithDetail is an interface that can be implemented by an error to provide
// an additional detail message about the error
type ErrorWithDetail interface {
	error
	DetailMsg() string
}

// HTTPError is the error response used by the serialization part of the framework.
type HTTPError struct {
	// Developer readable error message. Not shown to the user to avoid security leaks.
	Err error `json:"-" xml:"-" yaml:"-"`
	// URL of the error type. Can be used to lookup the error in a documentation
	Type string `json:"type,omitempty" xml:"type,omitempty" yaml:"type,omitempty" description:"URL of the error type. Can be used to lookup the error in a documentation"`
	// Short title of the error
	Title string `json:"title,omitempty" xml:"title,omitempty" yaml:"title,omitempty" description:"Short title of the error"`
	// HTTP status code. If using a different type than [HTTPError], for example [BadRequestError], this will be automatically overridden after Fuego error handling.
	Status int `json:"status,omitempty" xml:"status,omitempty" yaml:"status,omitempty" description:"HTTP status code" example:"403"`
	// Human readable error message
	Detail   string      `json:"detail,omitempty" xml:"detail,omitempty" yaml:"detail,omitempty" description:"Human readable error message"`
	Instance string      `json:"instance,omitempty" xml:"instance,omitempty" yaml:"instance,omitempty"`
	Errors   []ErrorItem `json:"errors,omitempty" xml:"errors,omitempty" yaml:"errors,omitempty"`
}

type ErrorItem struct {
	More   map[string]any `json:"more,omitempty" xml:"more,omitempty" description:"Additional information about the error"`
	Name   string         `json:"name" xml:"name" description:"For example, name of the parameter that caused the error"`
	Reason string         `json:"reason" xml:"reason" description:"Human readable error message"`
}

// PublicError returns a human readable error message.
// It ignores the underlying error for security and only returns the status code, title and detail.
func (e HTTPError) PublicError() string {
	var msgBuilder strings.Builder

	code := e.StatusCode()
	msgBuilder.WriteString(strconv.Itoa(code))

	title := e.Title
	if title == "" {
		title = http.StatusText(code)
		if title == "" {
			title = "HTTP Error"
		}
	}
	msgBuilder.WriteString(" ")
	msgBuilder.WriteString(title)

	if e.Detail != "" {
		msgBuilder.WriteString(" (")
		msgBuilder.WriteString(e.Detail)
		msgBuilder.WriteString(")")
	}

	return msgBuilder.String()
}

func (e HTTPError) Error() string {
	msg := e.PublicError()

	if e.Err != nil {
		msg = msg + ": " + e.Err.Error()
	}

	return msg
}

func (e HTTPError) StatusCode() int {
	if e.Status == 0 {
		return http.StatusInternalServerError
	}
	return e.Status
}

func (e HTTPError) DetailMsg() string {
	return e.Detail
}

func (e HTTPError) Unwrap() error { return e.Err }

// BadRequestError is an error used to return a 400 status code.
type BadRequestError HTTPError

var _ ErrorWithStatus = BadRequestError{}

func (e BadRequestError) Error() string {
	e.Status = http.StatusBadRequest
	return HTTPError(e).Error()
}

func (e BadRequestError) StatusCode() int { return http.StatusBadRequest }

func (e BadRequestError) Unwrap() error { return HTTPError(e) }

// NotFoundError is an error used to return a 404 status code.
type NotFoundError HTTPError

var _ ErrorWithStatus = NotFoundError{}

func (e NotFoundError) Error() string {
	e.Status = http.StatusNotFound
	return HTTPError(e).Error()
}

func (e NotFoundError) StatusCode() int { return http.StatusNotFound }

func (e NotFoundError) Unwrap() error { return HTTPError(e) }

// UnauthorizedError is an error used to return a 401 status code.
type UnauthorizedError HTTPError

var _ ErrorWithStatus = UnauthorizedError{}

func (e UnauthorizedError) Error() string {
	e.Status = http.StatusUnauthorized
	return HTTPError(e).Error()
}

func (e UnauthorizedError) StatusCode() int { return http.StatusUnauthorized }

func (e UnauthorizedError) Unwrap() error { return HTTPError(e) }

// ForbiddenError is an error used to return a 403 status code.
type ForbiddenError HTTPError

var _ ErrorWithStatus = ForbiddenError{}

func (e ForbiddenError) Error() string {
	e.Status = http.StatusForbidden
	return HTTPError(e).Error()
}

func (e ForbiddenError) StatusCode() int { return http.StatusForbidden }

func (e ForbiddenError) Unwrap() error { return HTTPError(e) }

// ConflictError is an error used to return a 409 status code.
type ConflictError HTTPError

var _ ErrorWithStatus = ConflictError{}

func (e ConflictError) Error() string {
	e.Status = http.StatusConflict
	return HTTPError(e).Error()
}

func (e ConflictError) StatusCode() int { return http.StatusConflict }

func (e ConflictError) Unwrap() error { return HTTPError(e) }

// InternalServerError is an error used to return a 500 status code.
type InternalServerError = HTTPError

// NotAcceptableError is an error used to return a 406 status code.
type NotAcceptableError HTTPError

var _ ErrorWithStatus = NotAcceptableError{}

func (e NotAcceptableError) Error() string {
	e.Status = http.StatusNotAcceptable
	return HTTPError(e).Error()
}

func (e NotAcceptableError) StatusCode() int { return http.StatusNotAcceptable }

func (e NotAcceptableError) Unwrap() error { return HTTPError(e) }

// ErrorHandler is the default error handler used by the framework.
// If the error is an [HTTPError] that error is returned.
// If the error adheres to the [ErrorWithStatus] interface
// the error is transformed to a [HTTPError] using [HandleHTTPError].
// If the error is not an [HTTPError] nor does it adhere to an
// interface the error is returned as is.
func ErrorHandler(ctx context.Context, err error) error {
	var errorStatus ErrorWithStatus
	switch {
	case errors.As(err, &HTTPError{}),
		errors.As(err, &errorStatus):
		return HandleHTTPError(ctx, err)
	}

	slog.ErrorContext(ctx, "Error in controller", "error", err.Error())

	return err
}

// HandleHTTPError is the core logic
// of handling fuego [HTTPError]'s. This
// function takes any error and coerces it into a fuego HTTPError.
// This can be used override the default handler:
//
//	engine := fuego.NewEngine(
//		WithErrorHandler(HandleHTTPError),
//	)
//
// or
//
//	server := fuego.NewServer(
//		fuego.WithEngineOptions(
//			fuego.WithErrorHandler(HandleHTTPError),
//		),
//	)
func HandleHTTPError(ctx context.Context, err error) error {
	errResponse := HTTPError{
		Err: err,
	}

	if v := castHTTPError(err); v != nil {
		errResponse = *v
	}

	// Check status code
	var errorStatus ErrorWithStatus
	if errors.As(err, &errorStatus) {
		errResponse.Status = errorStatus.StatusCode()
	}

	// Check for detail
	var errorDetail ErrorWithDetail
	if errors.As(err, &errorDetail) {
		errResponse.Detail = errorDetail.DetailMsg()
	}

	if errResponse.Title == "" {
		errResponse.Title = http.StatusText(errResponse.Status)
	}

	slog.ErrorContext(ctx, "Error "+errResponse.Title, "status", errResponse.StatusCode(), "detail", errResponse.DetailMsg(), "error", errResponse.Err)

	return errResponse
}

func castHTTPError(err error) *HTTPError {
	var errorInfo *HTTPError
	if errors.As(err, &errorInfo) {
		return errorInfo
	}
	var errorInfoValue HTTPError
	if errors.As(err, &errorInfoValue) {
		return &errorInfoValue
	}
	return nil
}
