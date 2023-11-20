package views

import (
	"github.com/go-op/op"
)

func (rs Ressource) pageAdmin(c op.Ctx[any]) (op.HTML, error) {
	recipes, err := rs.Queries.GetRecipes(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/admin.page.html", op.H{
		"Recipes": recipes,
	})
}

func (rs Ressource) deleteRecipe(c op.Ctx[any]) (any, error) {
	id := c.QueryParam("id") // TODO use PathParam
	err := rs.Queries.DeleteRecipe(c.Context(), id)
	if err != nil {
		return nil, err
	}

	return c.Redirect(301, "/recipes-list")
}
