package views

import (
	"github.com/go-op/op"
	"github.com/go-op/op/middleware/basicauth"
)

func (rs Ressource) Routes(s *op.Server) {
	// Public Pages
	op.GetStd(s, "/recipes-std", rs.showRecipesStd)
	op.Get(s, "/", rs.showRecipes)
	op.Get(s, "/recipes", rs.showRecipes)
	op.Post(s, "/recipes-new", rs.addRecipe)
	op.Get(s, "/ingredients", rs.showIngredients)
	op.Get(s, "/html", rs.showHTML)
	op.Get(s, "/h1string", rs.showString)

	// Public Chunks
	op.Get(s, "/recipes-list", rs.showRecipesList)
	op.Get(s, "/search", rs.searchRecipes)

	// Admin Pages
	adminRoutes := op.Group(s, "/admin")
	op.UseStd(adminRoutes, basicauth.New(basicauth.Config{Username: "admin", Password: "admin"}))
	op.Get(adminRoutes, "/", rs.pageAdmin)
	op.Get(adminRoutes, "/recipes", rs.adminRecipes)
	op.Post(adminRoutes, "/recipes-new", rs.adminAddRecipes)
	op.Get(adminRoutes, "/ingredients", rs.adminIngredients)
	op.Get(adminRoutes, "/users", rs.adminRecipes)

	// Admin Chunks
	op.Delete(adminRoutes, "/recipes/del", rs.deleteRecipe)
}
