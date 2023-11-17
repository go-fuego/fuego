package views

import "github.com/go-op/op"

func (rs Ressource) pageAdmin(c op.Ctx[any]) (op.HTML, error) {
	recipes, err := rs.Queries.GetRecipes(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/admin.page.html", op.H{
		"Recipes": recipes,
	})
}
