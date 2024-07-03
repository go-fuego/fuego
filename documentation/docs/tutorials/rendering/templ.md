# Templ

Fuego supports templating with [Templ](https://github.com/a-h/templ).

Simply return a Templ component from your handler,
with the `fuego.Templ` return type.

Example from [a recipe app](https://github.com/go-fuego/fuego/tree/main/examples/full-app-gourmet):

```go
// highlight-next-line
func (rs Ressource) adminIngredients(c fuego.ContextNoBody) (fuego.Templ, error) {
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

Note that the `fuego.Templ` type is a simple alias for `fuego.CtxRenderer`:
any type that implements the `Render(context.Context, io.Writer) error`
method can be used as a return type for a handler.
