package main

import (
	"flag"
	"log/slog"
	"os"

	"simple-crud/controller"
	"simple-crud/static"
	"simple-crud/store"
	"simple-crud/views"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/middleware/cache"
	"github.com/lmittmann/tint"
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
	db := store.InitDB("./recipe.db")

	// Create ressources that will be available in API controllers
	apiRessources := controller.NewRessource(db)

	// Create ressources that will be available in HTML controllers
	viewsRessources := views.NewRessource(db)

	// Create server with some options
	app := fuego.NewServer(
		fuego.WithPort(":8083"),
		fuego.WithAutoAuth(controller.LoginFunc),
		fuego.WithTemplateGlobs("**/*.html", "**/**/*.html"),
	)

	apiRessources.Security = app.Security

	// Register middlewares (functions that will be executed before AND after the controllers, in the order they are registered)
	// With fuego, you can use any existing middleware that relies on `net/http`, or create your own
	fuego.Use(app, cache.New(cache.Config{}))
	fuego.Use(app, chiMiddleware.Compress(5, "text/html", "text/css", "application/json"))

	app.Mux.Handle("/favicon.ico", static.Handler())

	// Register views (controllers that return HTML pages)
	viewsRessources.Routes(fuego.Group(app, ""))

	// Register API routes (controllers that return JSON)
	apiRessources.MountRoutes(fuego.Group(app, "/api"))

	// Run the server!
	app.Run()
}
