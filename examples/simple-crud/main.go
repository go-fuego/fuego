package main

import (
	"flag"
	"log/slog"
	"os"

	"simple-crud/controller"
	"simple-crud/store"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-op/op"
	"github.com/lmittmann/tint"
	"github.com/rs/cors"
)

//go:generate sqlc generate

func main() {
	// Flags
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	logLevel := slog.LevelInfo
	if *debug {
		logLevel = slog.LevelDebug
	}

	// Set my custom colored logger
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			AddSource:  true,
			Level:      logLevel,
			TimeFormat: "15:04:05",
		}),
	))

	// Connect to database
	db := store.InitDB("/tmp/recipe.db")

	// Create queries
	queries := store.New(db)

	// Create ressources that will be available in controllers
	rs := controller.NewRessource(*queries)

	// Create server with some options
	app := op.NewServer(
		op.WithPort(":8080"),
	)

	// Register middlewares (functions that will be executed before AND after the controllers, in the order they are registered)
	// With op, you can use any existing middleware that relies on `net/http`, or create your own
	op.UseStd(app, cors.Default().Handler)
	op.UseStd(app, chiMiddleware.Compress(5, "text/html", "text/css", "application/json"))

	// Register routes
	rs.Routes(app)

	// Run the server!
	app.Run()
}
