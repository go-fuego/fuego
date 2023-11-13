package views

import (
	"github.com/go-op/op"
)

func (rs Ressource) Routes(s *op.Server) {

	op.GetStd(s, "/recipes-std", rs.showRecipesStd)
	op.Get(s, "/recipes", rs.showRecipes)
	op.Get(s, "/html", rs.showHTML)
	op.Get(s, "/h1string", rs.showString)

}
