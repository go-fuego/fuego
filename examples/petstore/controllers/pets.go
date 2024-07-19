//coverage:ignore
package controller

import (
	"github.com/go-fuego/fuego"
)

type PetsRessources struct {
	PetsService PetsService
}

type Pets struct {
	ID   string `json:"id"`
	Name string `json:"name" example:"Napoleon"`
	Age  int    `json:"age" example:"18"`
}

type PetsCreate struct {
	Name string `json:"name" validate:"required,min=1,max=100" example:"Napoleon"`
	Age  int    `json:"age" validate:"max=100" example:"18"`
}

type PetsUpdate struct {
	Name string `json:"name" validate:"min=1,max=100" example:"Napoleon"`
	Age  int    `json:"age" validate:"max=100" example:"18"`
}

func (rs PetsRessources) Routes(s *fuego.Server) {
	petsGroup := fuego.Group(s, "/pets")

	fuego.Get(petsGroup, "/", rs.getAllPets)
	fuego.Post(petsGroup, "/", rs.postPets)

	fuego.Get(petsGroup, "/{id}", rs.getPets)
	fuego.Put(petsGroup, "/{id}", rs.putPets)
	fuego.Delete(petsGroup, "/{id}", rs.deletePets)
}

func (rs PetsRessources) getAllPets(c fuego.ContextNoBody) ([]Pets, error) {
	return rs.PetsService.GetAllPets()
}

func (rs PetsRessources) postPets(c *fuego.ContextWithBody[PetsCreate]) (Pets, error) {
	body, err := c.Body()
	if err != nil {
		return Pets{}, err
	}

	new, err := rs.PetsService.CreatePets(body)
	if err != nil {
		return Pets{}, err
	}

	return new, nil
}

func (rs PetsRessources) getPets(c fuego.ContextNoBody) (Pets, error) {
	id := c.PathParam("id")

	return rs.PetsService.GetPets(id)
}

func (rs PetsRessources) putPets(c *fuego.ContextWithBody[PetsUpdate]) (Pets, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return Pets{}, err
	}

	new, err := rs.PetsService.UpdatePets(id, body)
	if err != nil {
		return Pets{}, err
	}

	return new, nil
}

func (rs PetsRessources) deletePets(c *fuego.ContextNoBody) (any, error) {
	return rs.PetsService.DeletePets(c.PathParam("id"))
}

type PetsService interface {
	GetPets(id string) (Pets, error)
	CreatePets(PetsCreate) (Pets, error)
	GetAllPets() ([]Pets, error)
	UpdatePets(id string, input PetsUpdate) (Pets, error)
	DeletePets(id string) (any, error)
}
