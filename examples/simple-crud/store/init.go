package store

import (
	"database/sql"
	"errors"
	"log/slog"

	_ "embed"

	_ "modernc.org/sqlite"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite" // SQLite driver for migration
	_ "github.com/golang-migrate/migrate/v4/source/file"     // Migration files
)

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

	m, err := migrate.New("file://./store/migrations/", "sqlite://"+path)
	if err != nil {
		slog.Error("cannot migrate db", "err", err)
		panic("cannot migrate db")
	}

	err = m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		if err != nil {
			slog.Error("database migration failed",
				slog.Any("error", err),
				slog.String("database", "migrating: failure"))

			panic("database migration failed")
		}
		slog.Info("", slog.String("database", "migrating: success"))
	} else {
		slog.Info("", slog.String("database", "migrating: no change required"))
	}

	slog.Info("Database connected", "address", path)

	slog.Info("Database initialized")

	return db
}
