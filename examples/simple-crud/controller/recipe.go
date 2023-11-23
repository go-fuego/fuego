package controller

import (
	"net/http"

	"simple-crud/store"

	"github.com/go-fuego/fuego"
)

func (rs Ressource) getAllRecipesStandardWithHelpers(w http.ResponseWriter, r *http.Request) {
	recipes, err := rs.Queries.GetRecipes(r.Context())
	if err != nil {
		fuego.SendJSONError(w, err)
		return
	}

	fuego.SendJSON(w, recipes)
}

func (rs Ressource) getAllRecipes(c fuego.Ctx[any]) ([]store.Recipe, error) {
	recipes, err := rs.Queries.GetRecipes(c.Context())
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (rs Ressource) newRecipe(c fuego.Ctx[store.CreateRecipeParams]) (store.Recipe, error) {
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

func (rs Ressource) getRecipeWithIngredients(c fuego.Ctx[any]) ([]store.GetIngredientsOfRecipeRow, error) {
	recipe, err := rs.Queries.GetIngredientsOfRecipe(c.Context(), "uggjghj")
	if err != nil {
		return nil, err
	}

	return recipe, nil
}
