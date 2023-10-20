package controller

import (
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

func (rs Ressource) newIngredient(c op.Ctx[store.CreateIngredientParams]) (store.Ingredient, error) {
	body, err := c.Body()
	if err != nil {
		return store.Ingredient{}, err
	}

	ingredient, err := rs.Queries.CreateIngredient(c.Context(), body)
	if err != nil {
		return store.Ingredient{}, err
	}

	return ingredient, nil
}
