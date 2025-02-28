package sqlite

import (
	"github.com/mattn/go-sqlite3"

	"github.com/go-fuego/fuego"
)

// ErrorHandler maps sqlite3.Error errors to specific fuego errors.
// It returns a detailed fuego error including Title, Detail (the original error message),
// and Err (the original sqlite3.Error).
func ErrorHandler(err error) error {
	if err == nil {
		return nil
	}

	sqliteErr, ok := err.(*sqlite3.Error)
	if !ok {
		return err
	}

	switch sqliteErr.ExtendedCode {
	case sqlite3.ErrNoExtended(sqlite3.ErrConstraint):
		return fuego.ConflictError{
			Title:  "Constraint Violation",
			Detail: sqliteErr.Error(),
			Err:    sqliteErr,
		}
	case sqlite3.ErrConstraintUnique, sqlite3.ErrConstraintPrimaryKey:
		return fuego.ConflictError{
			Title:  "Duplicate",
			Detail: sqliteErr.Error(),
			Err:    sqliteErr,
		}
	case sqlite3.ErrConstraintForeignKey:
		return fuego.BadRequestError{
			Title:  "Foreign Key Constraint Failed",
			Detail: sqliteErr.Error(),
			Err:    sqliteErr,
		}
	case sqlite3.ErrNoExtended(sqlite3.ErrNotFound):
		return fuego.NotFoundError{
			Title:  "Record Not Found",
			Detail: sqliteErr.Error(),
			Err:    sqliteErr,
		}
	case sqlite3.ErrNoExtended(sqlite3.ErrPerm):
		return fuego.ForbiddenError{
			Title:  "Permission Denied",
			Detail: sqliteErr.Error(),
			Err:    sqliteErr,
		}
	case sqlite3.ErrNoExtended(sqlite3.ErrAuth):
		return fuego.UnauthorizedError{
			Title:  "Authentication Required",
			Detail: sqliteErr.Error(),
			Err:    sqliteErr,
		}
	case sqlite3.ErrNoExtended(sqlite3.ErrBusy), sqlite3.ErrNoExtended(sqlite3.ErrLocked):
		return fuego.InternalServerError{
			Title:  "Database Locked",
			Detail: sqliteErr.Error(),
			Err:    sqliteErr,
		}
	case sqlite3.ErrNoExtended(sqlite3.ErrIoErr):
		return fuego.InternalServerError{
			Title:  "I/O Error",
			Detail: sqliteErr.Error(),
			Err:    sqliteErr,
		}
	default:
		return fuego.InternalServerError{
			Title:  "Internal Server Error",
			Detail: sqliteErr.Error(),
			Err:    sqliteErr,
		}
	}
}
