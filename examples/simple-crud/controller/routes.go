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
	Now         func() time.Time       // Function to get the current time. Mocked in tests.
	Security    op.Security            // Security configuration
}

func (rs Ressource) Routes(s *op.Server) {
	op.GetStd(s, "/recipes-standard-with-helpers", rs.getAllRecipesStandardWithHelpers).
		AddTags("Recipe")

	op.Get(s, "/recipes", rs.getAllRecipes).
		WithSummary("Get all recipes").WithDescription("Get all recipes").
		WithQueryParam("limit", "number of recipes to return").
		AddTags("custom")

	op.Post(s, "/recipes/new", rs.newRecipe)
	op.Get(s, "/recipes/{id}", rs.getRecipeWithIngredients)
	op.Get(s, "/ingredients", rs.getAllIngredients)
	op.Post(s, "/ingredients/new", rs.newIngredient)
	op.Post(s, "/dosings/new", rs.newDosing)

	// Me ! Get the current user information
	op.Get(s, "/users/me", func(c op.Ctx[any]) (string, error) {
		claims, err := op.GetToken[MyCustomToken](c.Context())
		if err != nil {
			return "", err
		}

		return "My name is" + claims.Username, nil
	})

	op.Group(s, "/admin", func(s *op.Server) {
		op.UseStd(s, op.AuthWall("admin", "superadmin"))  // Only admin and superadmin can access the routes in this group
		op.UseStd(s, op.AuthWallRegex(`^(super)?admin$`)) // Same as above, but with a regex

		op.Get(s, "/users", placeholderController).
			WithDescription("Get all users").
			WithSummary("Get all users").
			SetTags("Admin")
	})

	op.Group(s, "/tests", func(s *op.Server) {
		op.Get(s, "/slow", slow).WithDescription("This is a slow route").WithSummary("Slow route")
		op.Get(s, "/mounted-route", placeholderController)
		op.Post(s, "/mounted-route-post", placeholderController)

		op.Group(s, "/mounted-group", func(groupedS *op.Server) {
			op.Get(groupedS, "/mounted-route", placeholderController)
		})
	})
}
