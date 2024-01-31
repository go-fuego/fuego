package views

import (
	"github.com/go-fuego/fuego"
)

type AdminRessource struct {
	DosingQueries      DosingRepository
	RecipesQueries     RecipeRepository
	IngredientsQueries IngredientRepository
}

func (rs Ressource) pageAdmin(c fuego.ContextNoBody) (fuego.Templ, error) {
	return rs.adminRecipes(c)
}
