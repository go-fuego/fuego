package controller

import (
	"simple-crud/store"

	"github.com/go-op/op"
)

func NewRessource(queries store.Queries) Ressource {
	return Ressource{
		Queries: queries,
	}
}

type Ressource struct {
	Queries store.Queries
}

func (rs Ressource) Routes(s *op.Server) {
	op.Get(s, "/recipes", rs.getAllRecipes).
		WithQueryParam("limit", "number of recipes to return").
		WithSummary("Get all recipes").
		WithDescription("Get all recipes").
		WithTags("custom")

	op.Post(s, "/recipes/new", rs.newRecipe)

	op.Get(s, "/recipes/{id}", rs.getRecipeWithIngredients)

	op.Get(s, "/ingredients", rs.getAllIngredients)
	op.Post(s, "/ingredients/new", rs.newIngredient)

	op.Post(s, "/dosings/new", rs.newDosing)
}
