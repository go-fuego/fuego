package services

import (
	"errors"
	"fmt"
	"log/slog"

	controller "github.com/go-fuego/fuego/examples/petstore/controllers"
	"github.com/go-fuego/fuego/examples/petstore/models"
)

func NewInMemoryPetsService() *InMemoryPetsService {
	return &InMemoryPetsService{
		Pets: []models.Pets{},
		Incr: new(int),
	}
}

type InMemoryPetsService struct {
	Pets []models.Pets
	Incr *int
}

// GetPetByName implements controller.PetsService.
func (petService *InMemoryPetsService) GetPetByName(name string) (models.Pets, error) {
	for _, p := range petService.Pets {
		if p.Name == name {
			return p, nil
		}
	}
	return models.Pets{}, errors.New("pet not found")
}

// CreatePets implements controller.PetsService.
func (petService *InMemoryPetsService) CreatePets(c models.PetsCreate) (models.Pets, error) {
	*petService.Incr++
	newPet := models.Pets{
		ID:   fmt.Sprintf("pet-%d", *petService.Incr),
		Name: c.Name,
		Age:  c.Age,
	}
	petService.Pets = append(petService.Pets, newPet)
	slog.Info("Created pet", "id", newPet.ID)

	return newPet, nil
}

// DeletePets implements controller.PetsService.
func (petService *InMemoryPetsService) DeletePets(id string) (any, error) {
	for i, p := range petService.Pets {
		if p.ID == id {
			petService.Pets = append(petService.Pets[:i], petService.Pets[i+1:]...)
			return nil, nil
		}
	}
	return nil, errors.New("pet not found")
}

// GetAllPets implements controller.PetsService.
func (petService *InMemoryPetsService) GetAllPets() ([]models.Pets, error) {
	return petService.Pets, nil
}

// GetAllPetsByAge implements controller.PetsService.
func (petService *InMemoryPetsService) GetAllPetsByAge() ([][]models.Pets, error) {
	maxAge := 0
	for _, p := range petService.Pets {
		if maxAge < p.Age {
			maxAge = p.Age
		}
	}
	pets := make([][]models.Pets, maxAge+1)
	for _, p := range petService.Pets {
		pets[p.Age] = append(pets[p.Age], p)
	}
	return pets, nil
}

// GetPets implements controller.PetsService.
func (petService *InMemoryPetsService) GetPets(id string) (models.Pets, error) {
	for _, p := range petService.Pets {
		if p.ID == id {
			return p, nil
		}
	}
	return models.Pets{}, errors.New("pet not found")
}

// UpdatePets implements controller.PetsService.
func (petService *InMemoryPetsService) UpdatePets(id string, input models.PetsUpdate) (models.Pets, error) {
	for i, p := range petService.Pets {
		if p.ID == id {
			if input.Name != "" {
				p.Name = input.Name
			}
			if input.Age != 0 {
				p.Age = input.Age
			}
			petService.Pets[i] = p
			return p, nil
		}
	}
	return models.Pets{}, errors.New("pet not found")
}

var _ controller.PetsService = &InMemoryPetsService{}
