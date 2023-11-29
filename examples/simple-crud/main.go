package main

import (
	"flag"
	"log/slog"
	"os"

	"simple-crud/controller"
	"simple-crud/static"
	"simple-crud/store"
	"simple-crud/views"

	"simple-crud/templates"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-fuego/fuego"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
)

//go:generate sqlc generate

func main() {
	// Load .env.local then .env files
	err := godotenv.Load(".env.local", ".env")
	if err != nil {
		wd, _ := os.Getwd()
		slog.Error("Error loading .env files: %s in dir %s", err, wd)
		return
	}

	// Flags
	port := flag.String("port", ":8083", "port to listen to")
	dbPath := flag.String("db", "./recipe.db", "path to database file")
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
	db := store.InitDB(*dbPath)

	// Create ressources that will be available in API controllers
	apiRessources := controller.NewRessource(db)

	// Create ressources that will be available in HTML controllers
	viewsRessources := views.NewRessource(db)

	// Create server with some options
	app := fuego.NewServer(
		fuego.WithPort(*port),
		fuego.WithAutoAuth(controller.LoginFunc),
		fuego.WithTemplateFS(templates.FS),
		fuego.WithTemplateGlobs("**/*.html", "**/**/*.html"),
	)

	apiRessources.Security = app.Security

	// Register middlewares (functions that will be executed before AND after the controllers, in the order they are registered)
	// With fuego, you can use any existing middleware that relies on `net/http`, or create your own
	fuego.Use(app, chiMiddleware.Compress(5, "text/html", "text/css", "application/json"))

	fuego.Handle(app, "/tailwind.min.css", static.Handler())
	fuego.Handle(app, "/favicon.ico", static.Handler())

	// Register views (controllers that return HTML pages)
	viewsRessources.Routes(fuego.Group(app, ""))

	// Register API routes (controllers that return JSON)
	apiRessources.MountRoutes(fuego.Group(app, "/api"))

	// Run the server!
	app.Run()
}
