---
sidebar_position: 1
---

# html/template

Fuego supports rendering HTML templates with the
[html/template](https://pkg.go.dev/html/template) package.

Just use the `fuego.HTML` type as a return type for your handler, and return
`c.Render()` with the template name and data.

```go
// highlight-next-line
func (rs Ressource) unitPreselected(c fuego.ContextNoBody) (fuego.HTML, error) {
	id := c.QueryParam("IngredientID")

	ingredient, err := rs.IngredientsQueries.GetIngredient(c.Context(), id)
	if err != nil {
		return "", err
	}

// highlight-start
	return c.Render("preselected-unit.partial.html", fuego.H{
		"Units":        types.UnitValues,
		"SelectedUnit": ingredient.DefaultUnit,
	})
// highlight-end
}
```
