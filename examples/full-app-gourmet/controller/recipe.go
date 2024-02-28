package controller

import (
	"context"
	"net/http"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/store"
)

type recipeRessource struct {
	RecipeRepository     RecipeRepository
	IngredientRepository IngredientRepository
}

func (rs recipeRessource) MountRoutes(s *fuego.Server) {
	fuego.GetStd(s, "/recipes-standard-with-helpers", rs.getAllRecipesStandardWithHelpers).
		AddTags("Recipe")

	fuego.Get(s, "/recipes", rs.getAllRecipes).
		Summary("Get all recipes").Description("Get all recipes").
		QueryParam("limit", "number of recipes to return").
		AddTags("custom")

	fuego.Post(s, "/recipes/new", rs.newRecipe)
	fuego.Get(s, "/recipes/{id}", rs.getRecipeWithIngredients)
}

func (rs recipeRessource) getAllRecipesStandardWithHelpers(w http.ResponseWriter, r *http.Request) {
	recipes, err := rs.RecipeRepository.GetRecipes(r.Context())
	if err != nil {
		fuego.SendJSONError(w, err)
		return
	}

	fuego.SendJSON(w, recipes)
}

func (rs recipeRessource) getAllRecipes(c fuego.ContextNoBody) ([]store.Recipe, error) {
	recipes, err := rs.RecipeRepository.GetRecipes(c.Context())
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (rs recipeRessource) newRecipe(c *fuego.ContextWithBody[store.CreateRecipeParams]) (store.Recipe, error) {
	body, err := c.Body()
	if err != nil {
		return store.Recipe{}, err
	}

	recipe, err := rs.RecipeRepository.CreateRecipe(c.Context(), body)
	if err != nil {
		return store.Recipe{}, err
	}

	return recipe, nil
}

func (rs recipeRessource) getRecipeWithIngredients(c fuego.ContextNoBody) (store.Recipe, error) {
	recipe, err := rs.RecipeRepository.GetRecipe(c.Context(), c.PathParam("id"))
	if err != nil {
		return store.Recipe{}, err
	}

	return recipe, nil
}

type RecipeRepository interface {
	CreateRecipe(ctx context.Context, arg store.CreateRecipeParams) (store.Recipe, error)
	DeleteRecipe(ctx context.Context, id string) error
	GetRecipe(ctx context.Context, id string) (store.Recipe, error)
	GetRecipes(ctx context.Context) ([]store.Recipe, error)
	SearchRecipes(ctx context.Context, args store.SearchRecipesParams) ([]store.Recipe, error)
}

var _ RecipeRepository = (*store.Queries)(nil)
