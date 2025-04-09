// Package sql provides error handling for SQL operations.
// It maps standard SQL errors to the corresponding fuego errors.
package sql

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-fuego/fuego"
)

// ErrorHandler maps standard SQL errors to the corresponding fuego errors.
func ErrorHandler(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fuego.NotFoundError{
			Err:    err,
			Title:  "Record Not Found",
			Detail: err.Error(),
			Status: http.StatusNotFound,
		}
	}

	if errors.Is(err, sql.ErrConnDone) {
		return fuego.InternalServerError{
			Err:    err,
			Title:  "Connection Closed",
			Detail: err.Error(),
			Status: http.StatusInternalServerError,
		}
	}

	if errors.Is(err, sql.ErrTxDone) {
		return fuego.ConflictError{
			Err:    err,
			Title:  "Transaction Completed",
			Detail: err.Error(),
			Status: http.StatusConflict,
		}
	}

	// For any other error, return the original error.
	return err
}
