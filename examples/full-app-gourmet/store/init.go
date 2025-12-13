package store

import (
	"database/sql"
	_ "embed"
	"errors"
	"log"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite" // SQLite driver for migration
	_ "github.com/golang-migrate/migrate/v4/source/file"     // Migration files
	"github.com/golang-migrate/migrate/v4/source/iofs"       // Migration files
	_ "modernc.org/sqlite"

	"github.com/go-fuego/fuego/examples/full-app-gourmet/store/migrations"
)

// InitDB initialize the database.
// SchemaPath is imported from schema.sql
func InitDB(path string) *sql.DB {
	// Open database with busy timeout parameter
	connectionString := path + "?_busy_timeout=5000"

	db, err := sql.Open("sqlite", connectionString)
	if err != nil {
		slog.Error("cannot open db connection", "err", err)
	}

	// Configure connection pool for better concurrency
	db.SetMaxOpenConns(25)   // Limit total connections
	db.SetMaxIdleConns(5)    // Keep some connections ready
	db.SetConnMaxLifetime(0) // Connections live forever (SQLite is local file)

	err = db.Ping()
	if err != nil {
		slog.Error("cannot ping db", "err", err)
	}

	// Explicitly enable WAL mode via PRAGMA (more reliable than connection string)
	// WAL mode allows concurrent readers and one writer without blocking
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		slog.Error("failed to enable WAL mode", "err", err)
	} else {
		slog.Info("WAL mode enabled for database")
	}

	// Set synchronous mode for better performance while maintaining safety
	_, err = db.Exec("PRAGMA synchronous=NORMAL;")
	if err != nil {
		slog.Error("failed to set synchronous mode", "err", err)
	}

	// Verify WAL mode is actually enabled
	var journalMode string
	err = db.QueryRow("PRAGMA journal_mode;").Scan(&journalMode)
	if err != nil {
		slog.Error("failed to query journal mode", "err", err)
	} else {
		slog.Info("Database journal mode", "mode", journalMode)
	}

	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithSourceInstance("embed://", d, "sqlite://"+path)
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
