package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/controller"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/server"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/store"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/views"
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
	port := flag.Int("port", 8083, "port to listen to")
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

	store := store.New(db)

	// Create resources that will be available in API controllers
	apiResources := controller.Resource{
		RecipesQueries:     store,
		IngredientsQueries: store,
		DosingQueries:      store,
	}

	// Create resources that will be available in HTML controllers
	viewsResources := views.Resource{
		RecipesQueries:     store,
		IngredientsQueries: store,
		DosingQueries:      store,
	}

	rs := server.Resources{
		API:   apiResources,
		Views: viewsResources,
	}

	app := rs.Setup(fuego.WithAddr(fmt.Sprintf("localhost:%d", *port)))

	// Run the server!
	err = app.Run()
	if err != nil {
		slog.Error("Error running server: %s", "err", err)
	}
}
