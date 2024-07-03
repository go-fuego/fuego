# Gomponents

Fuego supports rendering HTML templates with [Gomponents](https://github.com/maragudk/gomponents).

Just use the `fuego.Gomponent` type as a return type for your handler,
and return the gomponent.

```go
// highlight-next-line
func (rs Ressource) adminIngredients(c fuego.ContextNoBody) (fuego.Gomponent, error) {
	searchParams := components.SearchParams{
		Name:    c.QueryParam("name"),
		PerPage: c.QueryParamInt("perPage", 20),
		Page:    c.QueryParamInt("page", 1),
		URL:     "/admin/ingredients",
		Lang:    c.MainLang(),
	}

	ingredients, err := rs.IngredientsQueries.SearchIngredients(c.Context(), store.SearchIngredientsParams{
		Name:   "%" + searchParams.Name + "%",
		Limit:  int64(searchParams.PerPage),
		Offset: int64(searchParams.Page-1) * int64(searchParams.PerPage),
	})
	if err != nil {
		return nil, err
	}

// highlight-next-line
	return admin.IngredientList(ingredients, searchParams), nil
}
```

Note that the `fuego.Gomponent` type is a simple alias for `fuego.Renderer`:
any type that implements the `Render(io.Writer) error`
method can be used as a return type for a handler.
