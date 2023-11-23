package controller

import (
	"time"

	"simple-crud/store"

	"github.com/go-op/op"
	"github.com/rs/cors"
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
	op.Use(s, cors.Default().Handler)

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

	adminRoutes := op.Group(s, "/admin")
	op.Use(adminRoutes, op.AuthWall("admin", "superadmin"))  // Only admin and superadmin can access the routes in this group
	op.Use(adminRoutes, op.AuthWallRegex(`^(super)?admin$`)) // Same as above, but with a regex

	op.Get(adminRoutes, "/users", placeholderController).
		WithDescription("Get all users").
		WithSummary("Get all users").
		SetTags("Admin")

	testRoutes := op.Group(s, "/tests")
	op.Get(testRoutes, "/slow", slow).WithDescription("This is a slow route").WithSummary("Slow route")
	op.Get(testRoutes, "/mounted-route", placeholderController)
	op.Post(testRoutes, "/mounted-route-post", placeholderController)

	mountedGroup := op.Group(testRoutes, "/mounted-group")
	op.Get(mountedGroup, "/mounted-route", placeholderController)

	apiv2 := op.Group(s, "/v2")
	op.Get(apiv2, "/recipes", rs.getAllRecipes)
}
