package controller

import (
	"net/http"

	"simple-crud/store/ingredients"
	"simple-crud/store/recipes"

	"github.com/go-fuego/fuego"
)

type recipeRessource struct {
	recipeQueries      recipes.Queries
	ingredientsQueries ingredients.Queries
}

func (rs recipeRessource) MountRoutes(s *fuego.Server) {
	fuego.GetStd(s, "/recipes-standard-with-helpers", rs.getAllRecipesStandardWithHelpers).
		AddTags("Recipe")

	fuego.Get(s, "/recipes", rs.getAllRecipes).
		WithSummary("Get all recipes").WithDescription("Get all recipes").
		WithQueryParam("limit", "number of recipes to return").
		AddTags("custom")

	fuego.Post(s, "/recipes/new", rs.newRecipe)
	fuego.Get(s, "/recipes/{id}", rs.getRecipeWithIngredients)
}

func (rs recipeRessource) getAllRecipesStandardWithHelpers(w http.ResponseWriter, r *http.Request) {
	recipes, err := rs.recipeQueries.GetRecipes(r.Context())
	if err != nil {
		fuego.SendJSONError(w, err)
		return
	}

	fuego.SendJSON(w, recipes)
}

func (rs recipeRessource) getAllRecipes(c fuego.Ctx[any]) ([]recipes.Recipe, error) {
	recipes, err := rs.recipeQueries.GetRecipes(c.Context())
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (rs recipeRessource) newRecipe(c fuego.Ctx[recipes.CreateRecipeParams]) (recipes.Recipe, error) {
	body, err := c.Body()
	if err != nil {
		return recipes.Recipe{}, err
	}

	recipe, err := rs.recipeQueries.CreateRecipe(c.Context(), body)
	if err != nil {
		return recipes.Recipe{}, err
	}

	return recipe, nil
}

func (rs recipeRessource) getRecipeWithIngredients(c fuego.Ctx[any]) ([]ingredients.GetIngredientsOfRecipeRow, error) {
	recipe, err := rs.ingredientsQueries.GetIngredientsOfRecipe(c.Context(), c.QueryParam("id"))
	if err != nil {
		return nil, err
	}

	return recipe, nil
}
