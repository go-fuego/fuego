package server

import (
	"errors"

	"modernc.org/sqlite"

	"github.com/go-fuego/fuego"
	fuegosql "github.com/go-fuego/fuego/extra/sql"
)

func SQLiteErrorHandler(err error) error {
	var sqliteError *sqlite.Error
	if errors.As(err, &sqliteError) {
		sqliteErrorCode := sqliteError.Code()
		switch sqliteErrorCode {
		case 1555, 2067 /* UNIQUE constraint failed */ :
			return fuego.ConflictError{Title: "Duplicate", Detail: sqliteError.Error(), Err: sqliteError}
		default:
			return fuego.InternalServerError{Title: "Internal Server Error", Detail: sqliteError.Error(), Err: sqliteError}
		}
	}

	return err
}

func customErrorHandler(err error) error {
	return fuego.ErrorHandler(SQLiteErrorHandler(fuegosql.ErrorHandler(err)))
}
