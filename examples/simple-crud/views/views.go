package views

import (
	"github.com/go-op/op"
	"github.com/go-op/op/middleware/basicauth"
)

func (rs Ressource) Routes(s *op.Server) {
	// Pages
	op.GetStd(s, "/recipes-std", rs.showRecipesStd)
	op.Get(s, "/", rs.showRecipes)
	op.Get(s, "/recipes", rs.showRecipes)
	op.Get(s, "/search", rs.searchRecipes)
	op.Post(s, "/recipes-new", rs.addRecipe)
	op.Get(s, "/ingredients", rs.showIngredients)
	op.Get(s, "/html", rs.showHTML)
	op.Get(s, "/h1string", rs.showString)

	op.UseStd(s, basicauth.New(basicauth.Config{Username: "admin", Password: "admin"}))
	op.Get(s, "/admin", rs.pageAdmin)

}
