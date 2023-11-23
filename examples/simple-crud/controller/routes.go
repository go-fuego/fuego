package controller

import (
	"time"

	"simple-crud/store"

	"github.com/go-fuego/fuego"
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
	Security    fuego.Security         // Security configuration
}

func (rs Ressource) Routes(s *fuego.Server) {
	fuego.Use(s, cors.Default().Handler)

	fuego.GetStd(s, "/recipes-standard-with-helpers", rs.getAllRecipesStandardWithHelpers).
		AddTags("Recipe")

	fuego.Get(s, "/recipes", rs.getAllRecipes).
		WithSummary("Get all recipes").WithDescription("Get all recipes").
		WithQueryParam("limit", "number of recipes to return").
		AddTags("custom")

	fuego.Post(s, "/recipes/new", rs.newRecipe)
	fuego.Get(s, "/recipes/{id}", rs.getRecipeWithIngredients)
	fuego.Get(s, "/ingredients", rs.getAllIngredients)
	fuego.Post(s, "/ingredients/new", rs.newIngredient)
	fuego.Post(s, "/dosings/new", rs.newDosing)

	// Me ! Get the current user information
	fuego.Get(s, "/users/me", func(c fuego.Ctx[any]) (string, error) {
		claims, err := fuego.GetToken[MyCustomToken](c.Context())
		if err != nil {
			return "", err
		}

		return "My name is" + claims.Username, nil
	})

	adminRoutes := fuego.Group(s, "/admin")
	fuego.Use(adminRoutes, fuego.AuthWall("admin", "superadmin"))  // Only admin and superadmin can access the routes in this group
	fuego.Use(adminRoutes, fuego.AuthWallRegex(`^(super)?admin$`)) // Same as above, but with a regex

	fuego.Get(adminRoutes, "/users", placeholderController).
		WithDescription("Get all users").
		WithSummary("Get all users").
		SetTags("Admin")

	testRoutes := fuego.Group(s, "/tests")
	fuego.Get(testRoutes, "/slow", slow).WithDescription("This is a slow route").WithSummary("Slow route")
	fuego.Get(testRoutes, "/mounted-route", placeholderController)
	fuego.Post(testRoutes, "/mounted-route-post", placeholderController)

	mountedGroup := fuego.Group(testRoutes, "/mounted-group")
	fuego.Get(mountedGroup, "/mounted-route", placeholderController)

	apiv2 := fuego.Group(s, "/v2")
	fuego.Get(apiv2, "/recipes", rs.getAllRecipes)
}
