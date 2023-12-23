package views

import (
	"log/slog"
	"strconv"

	"simple-crud/store"
	"simple-crud/templa/admin"
	"simple-crud/templa/components"

	"github.com/go-fuego/fuego"
)

func (rs Ressource) adminOneIngredient(c fuego.Ctx[store.UpdateIngredientParams]) (fuego.CtxRenderer, error) {
	id := c.QueryParam("id") // TODO use PathParam

	if c.Request().Method == "PUT" {
		updateIngredientArgs, err := c.Body()
		if err != nil {
			return nil, err
		}

		updateIngredientArgs.ID = c.QueryParam("id") // TODO use PathParam

		_, err = rs.IngredientsQueries.UpdateIngredient(c.Context(), updateIngredientArgs)
		if err != nil {
			return nil, err
		}

		c.Response().Header().Set("HX-Trigger", "entity-updated")
	}

	ingredient, err := rs.IngredientsQueries.GetIngredient(c.Context(), id)
	if err != nil {
		return nil, err
	}

	slog.Debug("ingredient", "ingredient", ingredient, "strconv", strconv.FormatBool(ingredient.AvailableAllYear))

	return admin.IngredientPage(ingredient), nil
}

func (rs Ressource) adminIngredientCreationPage(c fuego.Ctx[store.CreateIngredientParams]) (any, error) {
	return admin.IngredientNew(), nil
}

func (rs Ressource) adminCreateIngredient(c fuego.Ctx[store.CreateIngredientParams]) (any, error) {
	createIngredientArgs, err := c.Body()
	if err != nil {
		return nil, err
	}

	_, err = rs.IngredientsQueries.CreateIngredient(c.Context(), createIngredientArgs)
	if err != nil {
		return nil, err
	}

	c.Response().Header().Set("HX-Trigger", "ingredient-created")

	return c.Redirect(301, "/admin/ingredients")
}

func (rs Ressource) adminIngredients(c fuego.Ctx[any]) (fuego.Templ, error) {
	searchParams := components.SearchParams{
		Name:    c.QueryParam("name"),
		PerPage: c.QueryParamInt("perPage", 20),
		Page:    c.QueryParamInt("page", 1),
		URL:     "/admin/ingredients",
		Lang:    c.MainLang(),
	}

	slog.Debug("params", "params", searchParams)

	ingredients, err := rs.IngredientsQueries.SearchIngredients(c.Context(), store.SearchIngredientsParams{
		Name:   "%" + searchParams.Name + "%",
		Limit:  int64(searchParams.PerPage),
		Offset: int64(searchParams.Page-1) * int64(searchParams.PerPage),
	})
	if err != nil {
		return nil, err
	}

	return admin.IngredientList(ingredients, searchParams), nil
}
