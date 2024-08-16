package lib

import (
	"github.com/go-fuego/fuego"
	controller "github.com/go-fuego/fuego/examples/petstore/controllers"
	"github.com/go-fuego/fuego/examples/petstore/services"
)

func NewPetStoreServer(options ...func(*fuego.Server)) *fuego.Server {
	s := fuego.NewServer(options...)

	petsRessources := controller.PetsRessources{
		PetsService: services.NewInMemoryPetsService(), // Dependency injection: we can pass a service here (for example a database service)
	}
	petsRessources.Routes(s)

	return s
}
