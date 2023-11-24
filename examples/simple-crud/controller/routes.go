package controller

import (
	"database/sql"
	"time"

	"simple-crud/store/dosings"
	"simple-crud/store/ingredients"
	"simple-crud/store/recipes"

	"github.com/go-fuego/fuego"
	"github.com/rs/cors"
)

// Ressource is the struct that holds useful sources of informations available for the controllers.
func NewRessource(db *sql.DB) Ressource {
	return Ressource{
		RecipesQueries:     *recipes.New(db),
		IngredientsQueries: *ingredients.New(db),
		DosingQueries:      *dosings.New(db),
	}
}

// Ressource is the global struct that holds useful sources of informations available for the controllers.
// Usually not used directly, but passed to the controllers.
type Ressource struct {
	DosingQueries      dosings.Queries
	RecipesQueries     recipes.Queries
	IngredientsQueries ingredients.Queries

	ExternalAPI interface{}            // External API
	Cache       map[string]interface{} // Some cache
	Now         func() time.Time       // Function to get the current time. Mocked in tests.
	Security    fuego.Security         // Security configuration
}

func (rs Ressource) MountRoutes(s *fuego.Server) {
	fuego.Use(s, cors.Default().Handler)

	recipeRessource{
		recipeQueries: rs.RecipesQueries,
	}.MountRoutes(s)

	ingredientRessource{
		Queries: rs.IngredientsQueries,
	}.MountRoutes(s)

	dosingRessource{
		Queries: rs.DosingQueries,
	}.MountRoutes(s)

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
}
