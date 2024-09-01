package controller

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/petstore/models"
)

type PetsResources struct {
	PetsService PetsService
}

func (rs PetsResources) Routes(s *fuego.Server) {
	petsGroup := fuego.Group(s, "/pets").Header("X-Header", "header description")

	fuego.Get(petsGroup, "/", rs.getAllPets)
	fuego.Get(petsGroup, "/by-age", rs.getAllPetsByAge)
	fuego.Post(petsGroup, "/", rs.postPets)

	fuego.Get(petsGroup, "/{id}", rs.getPets)
	fuego.Get(petsGroup, "/by-name/{name...}", rs.getPetByName)
	fuego.Put(petsGroup, "/{id}", rs.putPets)
	fuego.Put(petsGroup, "/{id}/json", rs.putPets).
		RequestContentType("application/json")
	fuego.Delete(petsGroup, "/{id}", rs.deletePets)
}

func (rs PetsResources) getAllPets(c fuego.ContextNoBody) ([]models.Pets, error) {
	return rs.PetsService.GetAllPets()
}

func (rs PetsResources) getAllPetsByAge(c fuego.ContextNoBody) ([][]models.Pets, error) {
	return rs.PetsService.GetAllPetsByAge()
}

func (rs PetsResources) postPets(c *fuego.ContextWithBody[models.PetsCreate]) (models.Pets, error) {
	body, err := c.Body()
	if err != nil {
		return models.Pets{}, err
	}

	return rs.PetsService.CreatePets(body)
}

func (rs PetsResources) getPets(c fuego.ContextNoBody) (models.Pets, error) {
	id := c.PathParam("id")

	return rs.PetsService.GetPets(id)
}

func (rs PetsResources) getPetByName(c fuego.ContextNoBody) (models.Pets, error) {
	name := c.PathParam("name")

	return rs.PetsService.GetPetByName(name)
}

func (rs PetsResources) putPets(c *fuego.ContextWithBody[models.PetsUpdate]) (models.Pets, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return models.Pets{}, err
	}

	return rs.PetsService.UpdatePets(id, body)
}

func (rs PetsResources) deletePets(c *fuego.ContextNoBody) (any, error) {
	return rs.PetsService.DeletePets(c.PathParam("id"))
}

type PetsService interface {
	GetPets(id string) (models.Pets, error)
	GetPetByName(name string) (models.Pets, error)
	CreatePets(models.PetsCreate) (models.Pets, error)
	GetAllPets() ([]models.Pets, error)
	GetAllPetsByAge() ([][]models.Pets, error)
	UpdatePets(id string, input models.PetsUpdate) (models.Pets, error)
	DeletePets(id string) (any, error)
}
