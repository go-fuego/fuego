package views

import (
	"os"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/middleware/basicauth"
	"github.com/go-fuego/fuego/middleware/cache"
)

func (rs Ressource) Routes(s *fuego.Server) {
	// Public Pages
	fuego.GetStd(s, "/recipes-std", rs.showRecipesStd)
	fuego.Get(s, "/", rs.showIndex, cache.New())
	fuego.Get(s, "/recipes", rs.showRecipes)
	fuego.Get(s, "/planner", rs.planner)
	fuego.Get(s, "/recipes/one", rs.showSingleRecipes2)
	fuego.Post(s, "/recipes-new", rs.addRecipe)
	fuego.Get(s, "/ingredients", rs.showIngredients)

	// Public Chunks
	fuego.Get(s, "/recipes-list", rs.showRecipesList)
	fuego.Get(s, "/search", rs.searchRecipes)
	fuego.Get(s, "/ingredients/preselect-unit", rs.unitPreselected).WithQueryParam("id", "")

	// Admin Pages
	adminRoutes := fuego.Group(s, "/admin")
	fuego.UseStd(adminRoutes, basicauth.New(basicauth.Config{Username: os.Getenv("ADMIN_USER"), Password: os.Getenv("ADMIN_PASSWORD")}))
	fuego.Get(adminRoutes, "", rs.pageAdmin)
	fuego.Get(adminRoutes, "/recipes", rs.adminRecipes)
	fuego.All(adminRoutes, "/recipes/one", rs.adminOneRecipe)
	fuego.Put(adminRoutes, "/recipes/edit", rs.editRecipe)
	fuego.Post(adminRoutes, "/recipes-new", rs.adminAddRecipes)
	fuego.Post(adminRoutes, "/dosings-new", rs.adminAddDosing)
	fuego.Get(adminRoutes, "/ingredients", rs.adminIngredients)
	fuego.Get(adminRoutes, "/ingredients/create", rs.adminIngredientCreationPage)
	fuego.All(adminRoutes, "/ingredients/one", rs.adminOneIngredient)

	fuego.Post(adminRoutes, "/ingredients/new", rs.adminCreateIngredient)
	fuego.Get(adminRoutes, "/users", rs.adminRecipes)

	// Admin Chunks
	fuego.Delete(adminRoutes, "/recipes/del", rs.deleteRecipe)
}
