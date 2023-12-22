package views

import (
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"simple-crud/store"
	"simple-crud/store/types"
	"simple-crud/templa/admin"
	"simple-crud/templa/components"

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

func (rs Ressource) adminOneIngredient(c fuego.Ctx[store.UpdateIngredientParams]) (fuego.CtxRenderer, error) {
	id := c.QueryParam("id") // TODO use PathParam

	if c.Request().Method == "PUT" {
		updateIngredientArgs, err := c.Body()
		if err != nil {
			return nil, err
		}

		updateIngredientArgs.ID = c.QueryParam("id") // TODO use PathParam

		_, err = rs.IngredientsQueries.UpdateIngredient(c.Context(), updateIngredientArgs)
		if err != nil {
			return nil, err
		}

		c.Response().Header().Set("HX-Trigger", "ingredient-updated")
	}

	ingredient, err := rs.IngredientsQueries.GetIngredient(c.Context(), id)
	if err != nil {
		return nil, err
	}

	slog.Debug("ingredient", "ingredient", ingredient, "strconv", strconv.FormatBool(ingredient.AvailableAllYear))

	return admin.IngredientPage(ingredient), nil
}

func (rs Ressource) adminIngredientCreationPage(c fuego.Ctx[store.CreateIngredientParams]) (any, error) {
	return admin.IngredientNew(), nil
}

func (rs Ressource) adminCreateIngredient(c fuego.Ctx[store.CreateIngredientParams]) (any, error) {
	createIngredientArgs, err := c.Body()
	if err != nil {
		return nil, err
	}

	_, err = rs.IngredientsQueries.CreateIngredient(c.Context(), createIngredientArgs)
	if err != nil {
		return nil, err
	}

	c.Response().Header().Set("HX-Trigger", "ingredient-created")

	return c.Redirect(301, "/admin/ingredients")
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

func (rs Ressource) adminIngredients(c fuego.Ctx[any]) (fuego.Templ, error) {
	searchParams := components.SearchParams{
		Name:    c.QueryParam("name"),
		PerPage: c.QueryParamInt("perPage", 20),
		Page:    c.QueryParamInt("page", 1),
		URL:     "/admin/ingredients",
		Lang:    c.MainLang(),
	}

	slog.Debug("params", "params", searchParams)

	ingredients, err := rs.IngredientsQueries.SearchIngredients(c.Context(), store.SearchIngredientsParams{
		Name:   "%" + searchParams.Name + "%",
		Limit:  int64(searchParams.PerPage),
		Offset: int64(searchParams.Page-1) * int64(searchParams.PerPage),
	})
	if err != nil {
		return nil, err
	}

	return admin.IngredientList(ingredients, searchParams), nil
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
