package sqlite

import (
	"errors"
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/mattn/go-sqlite3"
)

func TestConflictErrors(t *testing.T) {
	tests := []struct {
		name          string
		inputErr      error
		expectedTitle string
	}{
		{
			name:          "Constraint Violation",
			inputErr:      &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrConstraint)},
			expectedTitle: "Constraint Violation",
		},
		{
			name:          "Duplicate Entry - Unique Constraint",
			inputErr:      &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrConstraintUnique)},
			expectedTitle: "Duplicate Entry",
		},
		{
			name:          "Duplicate Entry - Primary Key Constraint",
			inputErr:      &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrConstraintPrimaryKey)},
			expectedTitle: "Duplicate Entry",
		},
	}

	for _, tc := range tests {
		tc := tc // capturar variable para ejecución paralela
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := SQLiteErrorHandler(tc.inputErr)
			var err fuego.ConflictError
			if !errors.As(result, &err) {
				t.Fatalf("expected ConflictError, got %T", result)
			}
			if err.Title != tc.expectedTitle {
				t.Errorf("expected title %q, got %q", tc.expectedTitle, err.Title)
			}
		})
	}
}

func TestBadRequestErrors(t *testing.T) {
	tests := []struct {
		name          string
		inputErr      error
		expectedTitle string
	}{
		{
			name:          "Foreign Key Constraint Failed",
			inputErr:      &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrConstraintForeignKey)},
			expectedTitle: "Foreign Key Constraint Failed",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := SQLiteErrorHandler(tc.inputErr)
			var err fuego.BadRequestError
			if !errors.As(result, &err) {
				t.Fatalf("expected BadRequestError, got %T", result)
			}
			if err.Title != tc.expectedTitle {
				t.Errorf("expected title %q, got %q", tc.expectedTitle, err.Title)
			}
		})
	}
}

func TestNotFoundError(t *testing.T) {
	inputErr := &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrNotFound)}
	result := SQLiteErrorHandler(inputErr)
	var err fuego.NotFoundError
	if !errors.As(result, &err) {
		t.Fatalf("expected NotFoundError, got %T", result)
	}
	if err.Title != "Record Not Found" {
		t.Errorf("expected title %q, got %q", "Record Not Found", err.Title)
	}
}

func TestUnauthorizedError(t *testing.T) {
	inputErr := &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrAuth)}
	result := SQLiteErrorHandler(inputErr)
	var err fuego.UnauthorizedError
	if !errors.As(result, &err) {
		t.Fatalf("expected UnauthorizedError, got %T", result)
	}
	if err.Title != "Authentication Required" {
		t.Errorf("expected title %q, got %q", "Authentication Required", err.Title)
	}
}

func TestInternalServerErrors(t *testing.T) {
	tests := []struct {
		name          string
		inputErr      error
		expectedTitle string
	}{
		{
			name:          "Database Locked - Busy",
			inputErr:      &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrBusy)},
			expectedTitle: "Database Locked",
		},
		{
			name:          "Database Locked - Locked",
			inputErr:      &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrLocked)},
			expectedTitle: "Database Locked",
		},
		{
			name:          "I/O Error",
			inputErr:      &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrIoErr)},
			expectedTitle: "I/O Error",
		},
		{
			name:          "Default Internal Server Error",
			inputErr:      &sqlite3.Error{ExtendedCode: 9999}, // código desconocido
			expectedTitle: "Internal Server Error",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := SQLiteErrorHandler(tc.inputErr)
			var err fuego.InternalServerError
			if !errors.As(result, &err) {
				t.Fatalf("expected InternalServerError, got %T", result)
			}
			if err.Title != tc.expectedTitle {
				t.Errorf("expected title %q, got %q", tc.expectedTitle, err.Title)
			}
		})
	}
}

func TestNonSqliteError(t *testing.T) {
	inputErr := errors.New("generic error")
	result := SQLiteErrorHandler(inputErr)
	if result != inputErr {
		t.Errorf("expected original error, got %v", result)
	}
}

func TestForbiddenError(t *testing.T) {
	inputErr := &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrPerm)}
	result := SQLiteErrorHandler(inputErr)
	var err fuego.ForbiddenError
	if !errors.As(result, &err) {
		t.Fatalf("expected ForbiddenError, got %T", result)
	}
	if err.Title != "Permission Denied" {
		t.Errorf("expected title %q, got %q", "Permission Denied", err.Title)
	}
}