package controller

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jub0bs/cors"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/middleware/basicauth"
	"github.com/go-fuego/fuego/option"
)

// Resource is the global struct that holds useful sources of information available for the controllers.
// Usually not used directly, but passed to the controllers.
type Resource struct {
	DosingQueries      DosingRepository
	RecipesQueries     RecipeRepository
	IngredientsQueries IngredientRepository

	ExternalAPI interface{}            // External API
	Cache       map[string]interface{} // Some cache
	Now         func() time.Time       // Function to get the current time. Mocked in tests.
	Security    fuego.Security         // Security configuration
}

func (rs Resource) MountRoutes(s *fuego.Server) {
	cors, err := cors.NewMiddleware(cors.Config{
		Origins:        []string{"*"},
		Methods:        []string{http.MethodGet, http.MethodHead, http.MethodPost},
		RequestHeaders: []string{"Accept", "Content-Type", "X-Requested-With"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fuego.Use(s, cors.Wrap)
	fuego.UseStd(s, basicauth.New(basicauth.Config{
		Username: os.Getenv("ADMIN_USER"),
		Password: os.Getenv("ADMIN_PASSWORD"),
		AllowGet: true,
	}))

	recipeResource{
		RecipeRepository:     rs.RecipesQueries,
		IngredientRepository: rs.IngredientsQueries,
	}.MountRoutes(s)

	ingredientResource{
		IngredientRepository: rs.IngredientsQueries,
	}.MountRoutes(s)

	dosingResource{
		Queries: rs.DosingQueries,
	}.MountRoutes(s)

	// Me ! Get the current user information
	fuego.Get(s, "/users/me", func(c fuego.ContextNoBody) (string, error) {
		claims, err := fuego.GetToken[MyCustomToken](c.Context())
		if err != nil {
			return "", err
		}

		return "My name is" + claims.Username, nil
	})

	adminRoutes := fuego.Group(s, "/admin")
	fuego.Use(adminRoutes, fuego.AuthWall("admin", "superadmin"))  // Only admin and superadmin can access the routes in this group
	fuego.Use(adminRoutes, fuego.AuthWallRegex(`^(super)?admin$`)) // Same as above, but with a regex

	testRoutes := fuego.Group(s, "/tests")
	fuego.Get(testRoutes, "/slow", slow,
		option.Description("This is a slow route"),
		option.Summary("Slow route"),
	)
}
