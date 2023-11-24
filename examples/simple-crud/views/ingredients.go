package views

import "github.com/go-fuego/fuego"

func (rs Ressource) showIngredients(c fuego.Ctx[any]) (fuego.HTML, error) {
	ingredients, err := rs.IngredientsQueries.GetIngredients(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/ingredients.page.html", ingredients)
}
