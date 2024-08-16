package controller

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/petstore/models"
)

type PetsRessources struct {
	PetsService PetsService
}

func (rs PetsRessources) Routes(s *fuego.Server) {
	petsGroup := fuego.Group(s, "/pets")

	fuego.Get(petsGroup, "/", rs.getAllPets)
	fuego.Get(petsGroup, "/by-age", rs.getAllPetsByAge)
	fuego.Post(petsGroup, "/", rs.postPets)

	fuego.Get(petsGroup, "/{id}", rs.getPets)
	fuego.Get(petsGroup, "/by-name/{name...}", rs.getPetByName)
	fuego.Put(petsGroup, "/{id}", rs.putPets)
	fuego.Delete(petsGroup, "/{id}", rs.deletePets)
}

func (rs PetsRessources) getAllPets(c fuego.ContextNoBody) ([]models.Pets, error) {
	return rs.PetsService.GetAllPets()
}

func (rs PetsRessources) getAllPetsByAge(c fuego.ContextNoBody) ([][]models.Pets, error) {
	return rs.PetsService.GetAllPetsByAge()
}

func (rs PetsRessources) postPets(c *fuego.ContextWithBody[models.PetsCreate]) (models.Pets, error) {
	body, err := c.Body()
	if err != nil {
		return models.Pets{}, err
	}

	return rs.PetsService.CreatePets(body)
}

func (rs PetsRessources) getPets(c fuego.ContextNoBody) (models.Pets, error) {
	id := c.PathParam("id")

	return rs.PetsService.GetPets(id)
}

func (rs PetsRessources) getPetByName(c fuego.ContextNoBody) (models.Pets, error) {
	name := c.PathParam("name")

	return rs.PetsService.GetPetByName(name)
}

func (rs PetsRessources) putPets(c *fuego.ContextWithBody[models.PetsUpdate]) (models.Pets, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return models.Pets{}, err
	}

	return rs.PetsService.UpdatePets(id, body)
}

func (rs PetsRessources) deletePets(c *fuego.ContextNoBody) (any, error) {
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
