package views

import (
	"simple-crud/store/ingredients"
	"simple-crud/store/recipes"

	"github.com/go-fuego/fuego"
)

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

	return c.Render("pages/admin/single-recipe.page.html", fuego.H{
		"Name":         recipe.Name,
		"Description":  recipe.Description,
		"Ingredients":  ingredients,
		"Instructions": nil,
	})
}

func (rs Ressource) adminAddRecipes(c fuego.Ctx[recipes.CreateRecipeParams]) (any, error) {
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

func (rs Ressource) adminIngredients(c fuego.Ctx[any]) (fuego.HTML, error) {
	ingredients, err := rs.IngredientsQueries.GetIngredients(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/admin/ingredients.page.html", fuego.H{
		"Ingredients": ingredients,
	})
}

func (rs Ressource) adminAddIngredient(c fuego.Ctx[ingredients.CreateIngredientParams]) (any, error) {
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
