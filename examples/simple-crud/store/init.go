package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	_ "embed"

	"modernc.org/sqlite"
)

//go:embed queries/_schema.sql
var schema []byte

// InitDB initialize the database.
// SchemaPath is imported from schema.sql
func InitDB(path string) *sql.DB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		slog.Error("cannot open db connexion", "err", err)
	}

	err = db.Ping()
	if err != nil {
		slog.Error("cannot ping db", "err", err)
	}

	slog.Info("Database connected", "address", path)

	_, err = db.Exec(string(schema))
	if err != nil {
		var sqliteError *sqlite.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.Code() == 1 {
				slog.Info("Database already initialized")
				return db
			}
		}

		slog.Error(fmt.Sprintf("%q: %s\n", err, string(schema)))
		return nil
	}

	slog.Info("Database initialized")

	return db
}
