package main

import (
	"flag"
	"log/slog"
	"os"

	"simple-crud/controller"
	"simple-crud/server"
	"simple-crud/store"
	"simple-crud/views"

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

	store := store.New(db)

	// Create ressources that will be available in API controllers
	apiRessources := controller.Ressource{
		RecipesQueries:     store,
		IngredientsQueries: store,
		DosingQueries:      store,
	}

	// Create ressources that will be available in HTML controllers
	viewsRessources := views.Ressource{
		RecipesQueries:     store,
		IngredientsQueries: store,
		DosingQueries:      store,
	}

	rs := server.Ressources{
		API:   apiRessources,
		Views: viewsRessources,
	}

	app := rs.Setup(fuego.WithPort(*port))

	// Run the server!
	app.Run()
}
