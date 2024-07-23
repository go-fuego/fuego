package views

import (
	"net/http"
	"os"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/middleware/basicauth"
	"github.com/go-fuego/fuego/middleware/cache"
)

func (rs Ressource) Routes(s *fuego.Server) {
	// Public Pages
	fuego.GetStd(s, "/recipes-std", rs.showRecipesStd)
	fuego.All(s, "/", rs.showIndex, cache.New())
	fuego.GetStd(s, "/robots.txt", rs.robots, cache.New())
	fuego.Get(s, "/recipes", rs.showRecipes)
	fuego.Get(s, "/planner", rs.planner)
	fuego.Get(s, "/recipes/{id}", rs.showSingleRecipes2)
	fuego.Get(s, "/recipes/{id}/related", rs.relatedRecipes)
	fuego.Post(s, "/recipes-new", rs.addRecipe)
	fuego.Get(s, "/ingredients", rs.showIngredients)
	fuego.Get(s, "/fast", rs.fastRecipes)
	fuego.Get(s, "/healthy", rs.healthyRecipes)

	// Public Chunks
	fuego.Get(s, "/recipes-list", rs.showRecipesList)
	fuego.Get(s, "/search", rs.searchRecipes).AddError(http.StatusUnauthorized, "Authorization Error").AddError(500, "My Server Error")
	fuego.Get(s, "/ingredients/preselect-unit", rs.unitPreselected).QueryParam("id", "")

	// Admin Pages
	adminRoutes := fuego.Group(s, "/admin")
	fuego.UseStd(adminRoutes, basicauth.New(basicauth.Config{Username: os.Getenv("ADMIN_USER"), Password: os.Getenv("ADMIN_PASSWORD")}))
	fuego.Get(adminRoutes, "", rs.pageAdmin)
	fuego.Get(adminRoutes, "/recipes", rs.adminRecipes)
	fuego.Get(adminRoutes, "/recipes/{id}", rs.adminOneRecipe)
	fuego.Put(adminRoutes, "/recipes/{id}", rs.adminOneRecipe)
	fuego.Delete(adminRoutes, "/recipes/{id}", rs.adminDeleteRecipe)
	fuego.Get(adminRoutes, "/recipes/create", rs.adminCreateRecipePage)
	fuego.Put(adminRoutes, "/recipes/edit", rs.editRecipe)
	fuego.Post(adminRoutes, "/recipes-new", rs.adminAddRecipes)
	fuego.Post(adminRoutes, "/dosings-new", rs.adminAddDosing)
	fuego.Get(adminRoutes, "/ingredients", rs.adminIngredients)
	fuego.Get(adminRoutes, "/ingredients/create", rs.adminIngredientCreationPage)
	fuego.All(adminRoutes, "/ingredients/{id}", rs.adminOneIngredient)

	fuego.Post(adminRoutes, "/ingredients/new", rs.adminCreateIngredient)
	fuego.Get(adminRoutes, "/users", rs.adminRecipes)
}
