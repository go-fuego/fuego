package views

import (
	"context"

	"simple-crud/store"

	"github.com/go-fuego/fuego"
)

func (rs Ressource) showIngredients(c fuego.Ctx[any]) (fuego.HTML, error) {
	ingredients, err := rs.IngredientsQueries.GetIngredients(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/ingredients.page.html", ingredients)
}

type IngredientRepository interface {
	CreateIngredient(ctx context.Context, arg store.CreateIngredientParams) (store.Ingredient, error)
	GetIngredient(ctx context.Context, id string) (store.Ingredient, error)
	GetIngredients(ctx context.Context) ([]store.Ingredient, error)
	GetIngredientsOfRecipe(ctx context.Context, recipeID string) ([]store.GetIngredientsOfRecipeRow, error)
	UpdateIngredient(ctx context.Context, arg store.UpdateIngredientParams) (store.Ingredient, error)
	SearchIngredients(ctx context.Context, arg store.SearchIngredientsParams) ([]store.Ingredient, error)
}

var _ IngredientRepository = (*store.Queries)(nil)
