package sqlite3

import (
	"errors"
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			inputErr:      &sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintUnique},
			expectedTitle: "Duplicate",
		},
		{
			name:          "Duplicate Entry - Primary Key Constraint",
			inputErr:      &sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintPrimaryKey},
			expectedTitle: "Duplicate",
		},
	}

	for _, tc := range tests {
		tc := tc 
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := ErrorHandler(tc.inputErr)
			var conflictErr fuego.ConflictError
			require.ErrorAs(t, result, &conflictErr, "expected ConflictError")
			assert.Equal(t, tc.expectedTitle, conflictErr.Title, "expected title to match")
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
			inputErr:      &sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintForeignKey},
			expectedTitle: "Foreign Key Constraint Failed",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := ErrorHandler(tc.inputErr)
			var badRequestErr fuego.BadRequestError
			require.ErrorAs(t, result, &badRequestErr, "expected BadRequestError")
			assert.Equal(t, tc.expectedTitle, badRequestErr.Title, "expected title to match")
		})
	}
}

func TestNotFoundError(t *testing.T) {
	inputErr := &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrNotFound)}
	result := ErrorHandler(inputErr)
	var notFoundErr fuego.NotFoundError
	require.ErrorAs(t, result, &notFoundErr, "expected NotFoundError")
	assert.Equal(t, "Record Not Found", notFoundErr.Title, "expected title to match")
}

func TestUnauthorizedError(t *testing.T) {
	inputErr := &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrAuth)}
	result := ErrorHandler(inputErr)
	var unauthorizedErr fuego.UnauthorizedError
	require.ErrorAs(t, result, &unauthorizedErr, "expected UnauthorizedError")
	assert.Equal(t, "Authentication Required", unauthorizedErr.Title, "expected title to match")
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
			inputErr:      &sqlite3.Error{ExtendedCode: 9999}, // unknown code
			expectedTitle: "Internal Server Error",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := ErrorHandler(tc.inputErr)
			var internalErr fuego.InternalServerError
			require.ErrorAs(t, result, &internalErr, "expected InternalServerError")
			assert.Equal(t, tc.expectedTitle, internalErr.Title, "expected title to match")
		})
	}
}

func TestNonSqliteError(t *testing.T) {
	inputErr := errors.New("generic error")
	result := ErrorHandler(inputErr)
	assert.Equal(t, inputErr, result, "expected original error to be returned")
}

func TestForbiddenError(t *testing.T) {
	inputErr := &sqlite3.Error{ExtendedCode: sqlite3.ErrNoExtended(sqlite3.ErrPerm)}
	result := ErrorHandler(inputErr)
	var forbiddenErr fuego.ForbiddenError
	require.ErrorAs(t, result, &forbiddenErr, "expected ForbiddenError")
	assert.Equal(t, "Permission Denied", forbiddenErr.Title, "expected title to match")
}