package views

import (
	"html/template"
	"log/slog"
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

	return c.Render("pages/recipes.page.html", recipes)
}

func (rs Ressource) searchRecipes(c op.Ctx[any]) (op.HTML, error) {
	recipes, err := rs.Queries.SearchRecipes(c.Context(), "%"+c.QueryParam("search")+"%")
	if err != nil {
		return "", err
	}

	slog.Debug("recipes", "recipes", recipes, "search", c.QueryParam("search"))

	return c.Render("partials/search-result.partial.html", op.H{
		"Recipes": recipes,
		"Search":  c.QueryParam("search"),
		"Filters": op.H{
			"Types":       []any{"Entrée", "Plat", "Dessert", "Apéritif"},
			"Ingredients": []any{"Poivron", "Tomate", "Oignon", "Ail", "Piment"},
		},
	})
}

func (rs Ressource) showRecipesList(c op.Ctx[any]) (op.HTML, error) {
	recipes, err := rs.Queries.SearchRecipes(c.Context(), "%"+c.QueryParam("search")+"%")
	if err != nil {
		return "", err
	}

	slog.Debug("recipes", "recipes", recipes, "search", c.QueryParam("search"))

	return c.Render("partials/recipes-list.partial.html", recipes)
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

	return c.Render("pages/admin.page.html", op.H{
		"Recipes": recipes,
	})
}

func (rs Ressource) showIngredients(c op.Ctx[any]) (op.HTML, error) {
	ingredients, err := rs.Queries.GetIngredients(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/ingredients.page.html", ingredients)
}

func (rs Ressource) showHTML(c op.Ctx[any]) (op.HTML, error) {
	return `<h1>test</h1>`, nil
}

func (rs Ressource) showString(c op.Ctx[any]) (string, error) {
	return `<h1>test</h1>`, nil
}
