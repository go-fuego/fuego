package sqlite

import (

	"github.com/mattn/go-sqlite3"
	"github.com/go-fuego/fuego"
)

func SQLiteErrorHandler(err error) error {
	if err == nil {
		return nil
	}
	sqliteErr, ok := err.(*sqlite3.Error)
	if !ok {
		return err
	}

	switch sqliteErr.ExtendedCode {
	case sqlite3.ErrNoExtended(sqlite3.ErrConstraint):
		return fuego.ConflictError{Title: "Constraint Violation"}
	case sqlite3.ErrNoExtended(sqlite3.ErrConstraintUnique), sqlite3.ErrNoExtended(sqlite3.ErrConstraintPrimaryKey):
		return fuego.ConflictError{Title: "Duplicate Entry"}
	case sqlite3.ErrNoExtended(sqlite3.ErrConstraintForeignKey):
		return fuego.BadRequestError{Title: "Foreign Key Constraint Failed"}
	case sqlite3.ErrNoExtended(sqlite3.ErrNotFound):
		return fuego.NotFoundError{Title: "Record Not Found"}
	case sqlite3.ErrNoExtended(sqlite3.ErrPerm):
		return fuego.ForbiddenError{Title: "Permission Denied"}
	case sqlite3.ErrNoExtended(sqlite3.ErrAuth):
		return fuego.UnauthorizedError{Title: "Authentication Required"}
	case sqlite3.ErrNoExtended(sqlite3.ErrBusy), sqlite3.ErrNoExtended(sqlite3.ErrLocked):
		return fuego.InternalServerError{Title: "Database Locked"}
	case sqlite3.ErrNoExtended(sqlite3.ErrIoErr):
		return fuego.InternalServerError{Title: "I/O Error"}
	default:
		return fuego.InternalServerError{Title: "Internal Server Error"}
	}
}