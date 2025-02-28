package sql

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-fuego/fuego"
)

// SQLErrorHandler maps standard SQL errors to the corresponding fuego errors.
func SQLErrorHandler(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fuego.NotFoundError{
			Err:    err,
			Title:  "Record Not Found",
			Detail: "The requested record was not found in the database",
			Status: http.StatusNotFound,
		}
	}

	if errors.Is(err, sql.ErrConnDone) {
		return fuego.InternalServerError{
			Err:    err,
			Title:  "Connection Closed",
			Detail: "The database connection is already closed",
			Status: http.StatusInternalServerError,
		}
	}

	if errors.Is(err, sql.ErrTxDone) {
		return fuego.ConflictError{
			Err:    err,
			Title:  "Transaction Completed",
			Detail: "The transaction has already been committed or rolled back",
			Status: http.StatusConflict,
		}
	}

	// For any other error, return the original error.
	return err
}