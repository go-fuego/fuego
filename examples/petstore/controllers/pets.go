//coverage:ignore
package controller

import (
	"context"
	"errors"

	"github.com/go-fuego/fuego"
)

type PetsRessources struct {
	PetsService PetsService
}

type Pets struct {
	ID   string `json:"id" validate:"required" example:"pet-123456"`
	Name string `json:"name" validate:"required" example:"Napoleon"`
	Age  int    `json:"age" example:"18" description:"Age of the pet, in years"`
}

type PetsCreate struct {
	Name string `json:"name" validate:"required,min=1,max=100" example:"Napoleon"`
	Age  int    `json:"age" validate:"max=100" example:"18" description:"Age of the pet, in years"`
}

type PetsUpdate struct {
	Name string `json:"name,omitempty" validate:"min=1,max=100" example:"Napoleon" description:"Name of the pet"`
	Age  int    `json:"age,omitempty" validate:"max=100" example:"18"`
}

var _ fuego.InTransformer = &Pets{}

func (*Pets) InTransform(context.Context) error {
	return errors.New("pets must only be used as output")
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
