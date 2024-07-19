package main

import (
	"github.com/go-fuego/fuego"
	controller "github.com/go-fuego/fuego/examples/petstore/controllers"
)

func newPetStoreServer(options ...func(*fuego.Server)) *fuego.Server {
	s := fuego.NewServer(options...)

	petsRessources := controller.PetsRessources{
		PetsService: nil, // Dependency injection: we can pass a service here (for example a database service)
	}
	petsRessources.Routes(s)

	return s
}

func main() {
	newPetStoreServer().Run()
}
