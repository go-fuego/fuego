package controller

import (
	"strings"

	"simple-crud/store"

	"github.com/go-fuego/fuego"
)

func (rs Ressource) getAllIngredients(c fuego.Ctx[any]) ([]store.Ingredient, error) {
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

func (rs Ressource) newIngredient(c fuego.Ctx[CreateIngredient]) (store.Ingredient, error) {
	body, err := c.Body()
	if err != nil {
		return store.Ingredient{}, err
	}

	payload := store.CreateIngredientParams{
		ID:          generateID(),
		Name:        body.Name,
		Description: body.Description,
	}

	ingredient, err := rs.Queries.CreateIngredient(c.Context(), payload)
	if err != nil {
		return store.Ingredient{}, err
	}

	return ingredient, nil
}
