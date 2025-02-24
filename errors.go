package fuego

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-fuego/fuego/internal"
)

type (
	ErrorWithStatus     = internal.ErrorWithStatus
	ErrorWithDetail     = internal.ErrorWithDetail
	HTTPError           = internal.HTTPError
	ErrorItem           = internal.ErrorItem
	BadRequestError     = internal.BadRequestError
	NotFoundError       = internal.NotFoundError
	UnauthorizedError   = internal.UnauthorizedError
	ForbiddenError      = internal.ForbiddenError
	ConflictError       = internal.ConflictError
	InternalServerError = internal.HTTPError
	NotAcceptableError  = internal.NotAcceptableError
)

// ErrorHandler is the default error handler used by the framework.
// If the error is an [HTTPError] that error is returned.
// If the error adheres to the [ErrorWithStatus] interface
// the error is transformed to a [HTTPError] using [HandleHTTPError].
// If the error is not an [HTTPError] nor does it adhere to an
// interface the error is returned as is.
func ErrorHandler(err error) error {
	var errorStatus ErrorWithStatus
	switch {
	case errors.As(err, &HTTPError{}),
		errors.As(err, &errorStatus):
		return HandleHTTPError(err)
	}

	slog.Error("Error in controller", "error", err.Error())

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
func HandleHTTPError(err error) error {
	errResponse := HTTPError{
		Err: err,
	}

	var errorInfo HTTPError
	if errors.As(err, &errorInfo) {
		errResponse = errorInfo
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

	slog.Error("Error "+errResponse.Title, "status", errResponse.StatusCode(), "detail", errResponse.DetailMsg(), "error", errResponse.Err)

	return errResponse
}
