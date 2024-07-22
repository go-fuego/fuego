package services

import (
	"testing"

	"github.com/go-fuego/fuego/examples/petstore/models"
	"github.com/stretchr/testify/require"
)

func TestInMemoryPets(t *testing.T) {
	service := NewInMemoryPetsService()

	t.Run("can create a pet", func(t *testing.T) {
		newPet, err := service.CreatePets(models.PetsCreate{Name: "kitkat", Age: 1})
		require.NoError(t, err)
		require.Equal(t, "pet-1", newPet.ID)
	})

	t.Run("can get a pet by name", func(t *testing.T) {
		newPet, err := service.GetPetByName("kitkat")
		require.NoError(t, err)
		require.Equal(t, "kitkat", newPet.Name)
		require.Equal(t, 1, newPet.Age)
	})

	t.Run("cannot get a pet by name if it doesn't exists", func(t *testing.T) {
		_, err := service.GetPetByName("snickers")
		require.Error(t, err)
	})

	t.Run("can get a pet by id", func(t *testing.T) {
		newPet, err := service.GetPets("pet-1")
		require.NoError(t, err)
		require.Equal(t, "kitkat", newPet.Name)
		require.Equal(t, 1, newPet.Age)
	})

	t.Run("can get all pets", func(t *testing.T) {
		pets, err := service.GetAllPets()
		require.NoError(t, err)
		require.Len(t, pets, 1)
	})

	t.Run("can update a pet", func(t *testing.T) {
		updatedPet, err := service.UpdatePets("pet-1", models.PetsUpdate{Name: "snickers", Age: 2})
		require.NoError(t, err)
		require.Equal(t, "snickers", updatedPet.Name)
		require.Equal(t, 2, updatedPet.Age)
	})

	t.Run("can delete a pet", func(t *testing.T) {
		_, err := service.DeletePets("pet-1")
		require.NoError(t, err)

		_, err = service.GetPets("pet-1")
		require.Error(t, err)
	})

	t.Run("cannot get a pet that does not exist", func(t *testing.T) {
		_, err := service.GetPets("pet-1")
		require.Error(t, err)
	})

	t.Run("cannot update a pet that does not exist", func(t *testing.T) {
		_, err := service.UpdatePets("pet-1", models.PetsUpdate{Name: "snickers", Age: 2})
		require.Error(t, err)
	})

	t.Run("cannot delete a pet that does not exist", func(t *testing.T) {
		_, err := service.DeletePets("pet-1")
		require.Error(t, err)
	})
}
