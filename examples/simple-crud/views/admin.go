package views

import "github.com/go-op/op"

func (rs Ressource) pageAdmin(c op.Ctx[any]) (op.HTML, error) {
	ingredients, err := rs.Queries.GetIngredients(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render(ingredients, "pages/admin.page.html")
}
