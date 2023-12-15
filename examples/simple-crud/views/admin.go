package views

import (
	"slices"
	"strings"

	"simple-crud/store"
	"simple-crud/store/types"
	"simple-crud/templa"

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

func (rs Ressource) deleteRecipe(c fuego.Ctx[any]) (any, error) {
	id := c.QueryParam("id") // TODO use PathParam
	err := rs.RecipesQueries.DeleteRecipe(c.Context(), id)
	if err != nil {
		return nil, err
	}

	return c.Redirect(301, "/admin/recipes")
}

func (rs Ressource) adminRecipes(c fuego.Ctx[any]) (fuego.HTML, error) {
	recipes, err := rs.RecipesQueries.GetRecipes(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/admin/recipes.page.html", fuego.H{
		"Recipes": recipes,
	})
}

func (rs Ressource) adminOneRecipe(c fuego.Ctx[any]) (fuego.HTML, error) {
	id := c.QueryParam("id") // TODO use PathParam

	recipe, err := rs.RecipesQueries.GetRecipe(c.Context(), id)
	if err != nil {
		return "", err
	}

	ingredients, err := rs.IngredientsQueries.GetIngredientsOfRecipe(c.Context(), id)
	if err != nil {
		return "", err
	}

	allIngredients, err := rs.IngredientsQueries.GetIngredients(c.Context())
	if err != nil {
		return "", err
	}

	slices.SortFunc(allIngredients, func(a, b store.Ingredient) int {
		return strings.Compare(a.Name, b.Name)
	})

	return c.Render("pages/admin/single-recipe.page.html", fuego.H{
		"Recipe":         recipe,
		"Ingredients":    ingredients,
		"Instructions":   nil,
		"AllIngredients": allIngredients,
		"Units":          types.UnitValues,
	})
}

func (rs Ressource) adminOneIngredient(c fuego.Ctx[any]) (any, error) {
	id := c.QueryParam("id") // TODO use PathParam

	ingredient, err := rs.IngredientsQueries.GetIngredient(c.Context(), id)
	if err != nil {
		return nil, err
	}

	err = templa.IngredientEdit(ingredient).Render(c.Context(), c.Response())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rs Ressource) editRecipe(c fuego.Ctx[store.UpdateRecipeParams]) (any, error) {
	updateRecipeArgs, err := c.Body()
	if err != nil {
		return "", err
	}

	updateRecipeArgs.ID = c.QueryParam("id") // TODO use PathParam

	recipe, err := rs.RecipesQueries.UpdateRecipe(c.Context(), updateRecipeArgs)
	if err != nil {
		return "", err
	}

	return c.Redirect(301, "/admin/recipes/one?id="+recipe.ID)
}

func (rs Ressource) adminAddRecipes(c fuego.Ctx[store.CreateRecipeParams]) (any, error) {
	body, err := c.Body()
	if err != nil {
		return "", err
	}

	_, err = rs.RecipesQueries.CreateRecipe(c.Context(), body)
	if err != nil {
		return "", err
	}

	return c.Redirect(301, "/admin/recipes")
}

func (rs Ressource) adminAddDosing(c fuego.Ctx[store.CreateDosingParams]) (any, error) {
	body, err := c.Body()
	if err != nil {
		return "", err
	}

	_, err = rs.DosingQueries.CreateDosing(c.Context(), body)
	if err != nil {
		return "", err
	}

	return c.Redirect(301, "/admin/recipes/one?id="+body.RecipeID)
}

func (rs Ressource) adminIngredients(c fuego.Ctx[any]) (any, error) {
	ingredients, err := rs.IngredientsQueries.GetIngredients(c.Context())
	if err != nil {
		return "", err
	}

	err = templa.IngredientList(ingredients).Render(c.Context(), c.Response())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (rs Ressource) adminAddIngredient(c fuego.Ctx[store.CreateIngredientParams]) (any, error) {
	body, err := c.Body()
	if err != nil {
		return "", err
	}

	_, err = rs.IngredientsQueries.CreateIngredient(c.Context(), body)
	if err != nil {
		return "", err
	}

	return c.Redirect(301, "/admin/ingredients")
}
