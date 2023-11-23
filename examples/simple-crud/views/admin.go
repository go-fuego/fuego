package views

import (
	"simple-crud/store"

	"github.com/go-op/op"
)

func (rs Ressource) pageAdmin(c op.Ctx[any]) (any, error) {
	return c.Redirect(301, "/admin/recipes")
}

func (rs Ressource) deleteRecipe(c op.Ctx[any]) (any, error) {
	id := c.QueryParam("id") // TODO use PathParam
	err := rs.Queries.DeleteRecipe(c.Context(), id)
	if err != nil {
		return nil, err
	}

	return c.Redirect(301, "/admin/recipes")
}

func (rs Ressource) adminRecipes(c op.Ctx[any]) (op.HTML, error) {
	recipes, err := rs.Queries.GetRecipes(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/admin/recipes.page.html", op.H{
		"Recipes": recipes,
	})
}

func (rs Ressource) adminAddRecipes(c op.Ctx[store.CreateRecipeParams]) (any, error) {
	body, err := c.Body()
	if err != nil {
		return "", err
	}

	_, err = rs.Queries.CreateRecipe(c.Context(), body)
	if err != nil {
		return "", err
	}

	return c.Redirect(301, "/admin/recipes")
}

func (rs Ressource) adminIngredients(c op.Ctx[any]) (op.HTML, error) {
	ingredients, err := rs.Queries.GetIngredients(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/admin/ingredients.page.html", op.H{
		"Ingredients": ingredients,
	})
}
