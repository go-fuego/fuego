package sql

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/go-fuego/fuego"
)

// TestSQLErrorHandler_NotFound verifies that sql.ErrNoRows is correctly mapped
// to a fuego.NotFoundError with the expected title and HTTP status.
func TestSQLErrorHandler_NotFound(t *testing.T) {
	err := sql.ErrNoRows
	result := SQLErrorHandler(err)
	var notFoundErr fuego.NotFoundError
	if !errors.As(result, &notFoundErr) {
		t.Fatalf("expected NotFoundError, got %T", result)
	}
	if notFoundErr.Title != "Record Not Found" {
		t.Errorf("expected title %q, got %q", "Record Not Found", notFoundErr.Title)
	}
}

// TestSQLErrorHandler_ConnDone verifies that sql.ErrConnDone is correctly mapped
// to a fuego.InternalServerError with the expected title and HTTP status.
func TestSQLErrorHandler_ConnDone(t *testing.T) {
	err := sql.ErrConnDone
	result := SQLErrorHandler(err)
	var internalErr fuego.InternalServerError
	if !errors.As(result, &internalErr) {
		t.Fatalf("expected InternalServerError, got %T", result)
	}
	if internalErr.Title != "Connection Closed" {
		t.Errorf("expected title %q, got %q", "Connection Closed", internalErr.Title)
	}
}

// TestSQLErrorHandler_TxDone verifies that sql.ErrTxDone is correctly mapped
// to a fuego.ConflictError with the expected title and HTTP status.
func TestSQLErrorHandler_TxDone(t *testing.T) {
	err := sql.ErrTxDone
	result := SQLErrorHandler(err)
	var conflictErr fuego.ConflictError
	if !errors.As(result, &conflictErr) {
		t.Fatalf("expected ConflictError, got %T", result)
	}
	if conflictErr.Title != "Transaction Completed" {
		t.Errorf("expected title %q, got %q", "Transaction Completed", conflictErr.Title)
	}
}

// TestSQLErrorHandler_Generic verifies that errors that are not standard SQL errors
// are returned unchanged.
func TestSQLErrorHandler_Generic(t *testing.T) {
	genericErr := errors.New("generic error")
	result := SQLErrorHandler(genericErr)
	if result != genericErr {
		t.Errorf("expected original error, got %v", result)
	}
}