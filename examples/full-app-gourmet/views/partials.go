package views

import (
	"simple-crud/store/types"

	"github.com/go-fuego/fuego"
)

func (rs Ressource) unitPreselected(c fuego.ContextNoBody) (fuego.HTML, error) {
	id := c.QueryParam("IngredientID")

	ingredient, err := rs.IngredientsQueries.GetIngredient(c.Context(), id)
	if err != nil {
		return "", err
	}

	return c.Render("preselected-unit.partial.html", fuego.H{
		"Units":        types.UnitValues,
		"SelectedUnit": ingredient.DefaultUnit,
	})
}
