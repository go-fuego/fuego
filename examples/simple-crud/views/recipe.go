package views

import (
	"html/template"
	"net/http"
	"path"

	"simple-crud/store"

	"github.com/go-op/op"
	"github.com/google/uuid"
)

func NewRessource(queries store.Queries) Ressource {
	return Ressource{
		Queries: queries,
	}
}

// Ressource is the struct that holds useful sources of informations available for the controllers.
type Ressource struct {
	Queries store.Queries // Database queries
}

func (rs Ressource) showRecipesStd(w http.ResponseWriter, r *http.Request) {
	recipes, err := rs.Queries.GetRecipes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fp := path.Join("templates", "recipes.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, recipes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rs Ressource) showRecipes(c op.Ctx[any]) (op.HTML, error) {
	recipes, err := rs.Queries.GetRecipes(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render(recipes, "pages/recipes.page.html", "layouts/main.layout.html", "partials/**.html")
}

func (rs Ressource) addRecipe(c op.Ctx[store.CreateRecipeParams]) (op.HTML, error) {
	body, err := c.Body()
	if err != nil {
		return "", err
	}

	body.ID = uuid.NewString()

	_, err = rs.Queries.CreateRecipe(c.Context(), body)
	if err != nil {
		return "", err
	}

	recipes, err := rs.Queries.GetRecipes(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render(recipes, "pages/admin.page.html", "layouts/main.layout.html", "partials/**.html")
}

func (rs Ressource) showIngredients(c op.Ctx[any]) (op.HTML, error) {
	ingredients, err := rs.Queries.GetIngredients(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render(ingredients, "pages/ingredients.page.html", "layouts/main.layout.html", "partials/**.html")
}

func (rs Ressource) showHTML(c op.Ctx[any]) (op.HTML, error) {
	return `<h1>test</h1>`, nil
}

func (rs Ressource) showString(c op.Ctx[any]) (string, error) {
	return `<h1>test</h1>`, nil
}
