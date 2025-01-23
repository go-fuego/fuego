package main

import (
	"github.com/go-fuego/fuego"
	"github.com/sonkeydotcom/fuego/examples/crud-gorm/handlers"
	"github.com/sonkeydotcom/fuego/examples/crud-gorm/models"
	"github.com/sonkeydotcom/fuego/examples/crud-gorm/queries"

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
