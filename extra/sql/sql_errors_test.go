package sql

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSQLErrorHandler_NotFound verifies that sql.ErrNoRows is correctly mapped
// to a fuego.NotFoundError with the expected title.
func TestSQLErrorHandler_NotFound(t *testing.T) {
	err := sql.ErrNoRows
	result := ErrorHandler(err)
	var notFoundErr fuego.NotFoundError

	// Verify that result can be cast to fuego.NotFoundError.
	require.ErrorAs(t, result, &notFoundErr, "expected NotFoundError")
	// Assert that the error title is correct.
	assert.Equal(t, "Record Not Found", notFoundErr.Title)
}

// TestSQLErrorHandler_ConnDone verifies that sql.ErrConnDone is correctly mapped
// to a fuego.InternalServerError with the expected title.
func TestSQLErrorHandler_ConnDone(t *testing.T) {
	err := sql.ErrConnDone
	result := ErrorHandler(err)
	var internalErr fuego.InternalServerError

	require.ErrorAs(t, result, &internalErr, "expected InternalServerError")
	assert.Equal(t, "Connection Closed", internalErr.Title)
}

// TestSQLErrorHandler_TxDone verifies that sql.ErrTxDone is correctly mapped
// to a fuego.ConflictError with the expected title.
func TestSQLErrorHandler_TxDone(t *testing.T) {
	err := sql.ErrTxDone
	result := ErrorHandler(err)
	var conflictErr fuego.ConflictError

	require.ErrorAs(t, result, &conflictErr, "expected ConflictError")
	assert.Equal(t, "Transaction Completed", conflictErr.Title)
}

// TestSQLErrorHandler_Generic verifies that errors not recognized by the handler
// are returned unchanged.
func TestSQLErrorHandler_Generic(t *testing.T) {
	genericErr := errors.New("generic error")
	result := ErrorHandler(genericErr)
	assert.Equal(t, genericErr, result)
}