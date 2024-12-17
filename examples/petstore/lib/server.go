package lib

import (
	"net/http"

	"github.com/go-fuego/fuego"
	controller "github.com/go-fuego/fuego/examples/petstore/controllers"
	"github.com/go-fuego/fuego/examples/petstore/services"
	"github.com/go-fuego/fuego/option"
)

type NoContent struct {
	Empty string `json:"-"`
}

func NewPetStoreServer(options ...func(*fuego.Server)) *fuego.Server {
	options = append(options, fuego.WithRouteOptions(
		option.AddResponse(http.StatusNoContent, "No Content", fuego.Response{Type: NoContent{}}),
	))
	s := fuego.NewServer(options...)

	petsResources := controller.PetsResources{
		PetsService: services.NewInMemoryPetsService(), // Dependency injection: we can pass a service here (for example a database service)
	}
	petsResources.Routes(s)

	return s
}
