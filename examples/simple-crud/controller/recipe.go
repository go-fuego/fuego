package controller

import (
	"simple-crud/store"

	"github.com/go-op/op"
)

func (rs Ressource) getAllRecipes(c op.Ctx[any]) ([]store.Recipe, error) {
	recipes, err := rs.Queries.GetRecipes(c.Context())
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (rs Ressource) newRecipe(c op.Ctx[store.CreateRecipeParams]) (store.Recipe, error) {
	body, err := c.Body()
	if err != nil {
		return store.Recipe{}, err
	}

	recipe, err := rs.Queries.CreateRecipe(c.Context(), body)
	if err != nil {
		return store.Recipe{}, err
	}

	return recipe, nil
}

func (rs Ressource) getRecipeWithIngredients(c op.Ctx[any]) ([]store.GetIngredientsOfRecipeRow, error) {

	recipe, err := rs.Queries.GetIngredientsOfRecipe(c.Context(), "uggjghj")
	if err != nil {
		return nil, err
	}

	return recipe, nil
}
