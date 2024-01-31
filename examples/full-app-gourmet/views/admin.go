package views

import (
	"github.com/go-fuego/fuego"
)

type AdminRessource struct {
	DosingQueries      DosingRepository
	RecipesQueries     RecipeRepository
	IngredientsQueries IngredientRepository
}

func (rs Ressource) pageAdmin(c fuego.Ctx[any]) (any, error) {
	return c.Redirect(301, "/admin/recipes")
}
