package main

import (
	"crud-gorm/handlers"
	"crud-gorm/models"
	"crud-gorm/queries"

	"github.com/go-fuego/fuego"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, error := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if error != nil {
		panic("error connecting to  database")
	}

	db.AutoMigrate(&models.User{})

	server := fuego.NewServer()

	userQueries := &queries.UserQueries{DB: db}
	handlers := &handlers.Handlers{UserQueries: userQueries}

	fuego.Get(server, "/users", handlers.GetUsers)
	fuego.Post(server, "/users", handlers.CreateUser)
	fuego.Get(server, "/users/:id", handlers.GetUserByID)
	fuego.Put(server, "/users/:id", handlers.UpdateUser)
	fuego.Delete(server, "/users/:id", handlers.DeleteUser)

	server.Run()

}
