package controller

import (
	"crypto/rand"
	"strings"

	"simple-crud/store"

	"github.com/go-op/op"
)

func (rs Ressource) getAllIngredients(c op.Ctx[any]) ([]store.Ingredient, error) {
	ingredients, err := rs.Queries.GetIngredients(c.Context())
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

type CreateIngredient struct {
	Name        string `json:"name" validate:"required,min=3,max=20"`
	Description string `json:"description"`
}

func (ci *CreateIngredient) InTransform() error {
	if ci.Description == "" {
		ci.Description = "No description"
	}
	ci.Name = strings.TrimSpace(ci.Name)
	return nil
}

func (rs Ressource) newIngredient(c op.Ctx[CreateIngredient]) (store.Ingredient, error) {
	body, err := c.Body()
	if err != nil {
		return store.Ingredient{}, err
	}

	// Generate random string
	id := make([]byte, 10)
	rand.Read(id)

	payload := store.CreateIngredientParams{
		ID:          string(id),
		Name:        body.Name,
		Description: body.Description,
	}

	ingredient, err := rs.Queries.CreateIngredient(c.Context(), payload)
	if err != nil {
		return store.Ingredient{}, err
	}

	return ingredient, nil
}
