package views

import (
	"simple-crud/store"

	"github.com/go-fuego/fuego"
)

type Unit struct {
	Name string
}

func (rs Ressource) unitPreselected(c fuego.Ctx[Unit]) (fuego.HTML, error) {
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
