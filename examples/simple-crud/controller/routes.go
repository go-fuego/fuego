package controller

import (
	"time"

	"simple-crud/store"

	"github.com/go-op/op"
)

// Ressource is the struct that holds useful sources of informations available for the controllers.
func NewRessource(queries store.Queries) Ressource {
	return Ressource{
		Queries: queries,
	}
}

// Ressource is the struct that holds useful sources of informations available for the controllers.
type Ressource struct {
	Queries     store.Queries          // Database queries
	UserQueries store.Queries          // Database queries from another store
	ExternalAPI interface{}            // External API
	Cache       map[string]interface{} // Some cache
	Now         func() time.Time       // Function to get the current time
}

func (rs Ressource) Routes(s *op.Server) {
	op.GetStd(s, "/recipes-standard-with-helpers", rs.getAllRecipesStandardWithHelpers).
		AddTags("Recipe")

	op.Get(s, "/recipes", rs.getAllRecipes).
		WithQueryParam("limit", "number of recipes to return").
		WithSummary("Get all recipes").
		WithDescription("Get all recipes").
		AddTags("custom")

	op.Post(s, "/recipes/new", rs.newRecipe)

	op.Get(s, "/recipes/{id}", rs.getRecipeWithIngredients)

	op.Get(s, "/ingredients", rs.getAllIngredients)
	op.Post(s, "/ingredients/new", rs.newIngredient)

	op.Post(s, "/dosings/new", rs.newDosing)

	op.Group(s, "/api", func(s *op.Server) {
		op.Get(s, "/mounted-route", func(c op.Ctx[any]) (string, error) {
			return "hello", nil
		})

		op.Post(s, "/mounted-route-post", func(c op.Ctx[any]) (string, error) {
			return "hello", nil
		})

		op.Group(s, "/mounted-group", func(groupedS *op.Server) {
			op.Get(groupedS, "/mounted-route", func(c op.Ctx[any]) (string, error) {
				return "hello", nil
			})
		})
	})
}
