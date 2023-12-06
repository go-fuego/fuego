package views

import (
	"simple-crud/store"

	"github.com/go-fuego/fuego"
)

func (rs Ressource) unitPreselected(c fuego.Ctx[any]) (fuego.HTML, error) {
	id := c.QueryParam("IngredientID")

	ingredient, err := rs.IngredientsQueries.GetIngredient(c.Context(), id)
	if err != nil {
		return "", err
	}

	return c.Render("preselected-unit.partial.html", fuego.H{
		"Units":        store.UnitValues,
		"SelectedUnit": ingredient.DefaultUnit,
	})
}
