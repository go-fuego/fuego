package controller

import (
	"strings"

	"simple-crud/store/ingredients"

	"github.com/go-fuego/fuego"
)

type ingredientRessource struct {
	Queries ingredients.Queries
}

func (rs ingredientRessource) MountRoutes(s *fuego.Server) {
	fuego.Get(s, "/ingredients", rs.getAllIngredients)
	fuego.Post(s, "/ingredients/new", rs.newIngredient)
}

func (rs ingredientRessource) getAllIngredients(c fuego.Ctx[any]) ([]ingredients.Ingredient, error) {
	ingredients, err := rs.Queries.GetIngredients(c.Context())
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

type CreateIngredient struct {
	Name        string `json:"name" validate:"required,min=3,max=20"`
	Description string `json:"description"`
}

func (ci *CreateIngredient) InTransform() error {
	if ci.Description == "" {
		ci.Description = "No description"
	}
	ci.Name = strings.TrimSpace(ci.Name)
	return nil
}

func (rs ingredientRessource) newIngredient(c fuego.Ctx[CreateIngredient]) (ingredients.Ingredient, error) {
	body, err := c.Body()
	if err != nil {
		return ingredients.Ingredient{}, err
	}

	payload := ingredients.CreateIngredientParams{
		ID:          generateID(),
		Name:        body.Name,
		Description: body.Description,
	}

	ingredient, err := rs.Queries.CreateIngredient(c.Context(), payload)
	if err != nil {
		return ingredients.Ingredient{}, err
	}

	return ingredient, nil
}
