package views

import (
	"html/template"
	"log/slog"
	"net/http"
	"path"

	"simple-crud/store"

	"github.com/go-fuego/fuego"
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

func (rs Ressource) showIndex(c fuego.Ctx[any]) (fuego.HTML, error) {
	return c.Render("pages/index.page.html", nil)
}

func (rs Ressource) showRecipes(c fuego.Ctx[any]) (fuego.HTML, error) {
	recipes, err := rs.Queries.GetRecipes(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/recipes.page.html", recipes)
}

func (rs Ressource) searchRecipes(c fuego.Ctx[any]) (fuego.HTML, error) {
	search := c.QueryParam("q")

	recipes, err := rs.Queries.SearchRecipes(c.Context(), "%"+search+"%")
	if err != nil {
		return "", err
	}

	slog.Debug("recipes", "recipes", recipes, "search", search)

	return c.Render("pages/search.page.html", fuego.H{
		"Recipes": recipes,
		"Search":  search,
		"Filters": fuego.H{
			"Types":       []any{"Entrée", "Plat", "Dessert", "Apéritif"},
			"Ingredients": []any{"Poivron", "Tomate", "Oignon", "Ail", "Piment"},
		},
	})
}

func (rs Ressource) showRecipesList(c fuego.Ctx[any]) (fuego.HTML, error) {
	recipes, err := rs.Queries.SearchRecipes(c.Context(), "%"+c.QueryParam("search")+"%")
	if err != nil {
		return "", err
	}

	slog.Debug("recipes", "recipes", recipes, "search", c.QueryParam("search"))

	return c.Render("partials/recipes-list.partial.html", recipes)
}

func (rs Ressource) addRecipe(c fuego.Ctx[store.CreateRecipeParams]) (fuego.HTML, error) {
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

	return c.Render("pages/admin.page.html", fuego.H{
		"Recipes": recipes,
	})
}

func (rs Ressource) showIngredients(c fuego.Ctx[any]) (fuego.HTML, error) {
	ingredients, err := rs.Queries.GetIngredients(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/ingredients.page.html", ingredients)
}

func (rs Ressource) showHTML(c fuego.Ctx[any]) (fuego.HTML, error) {
	return `<h1>test</h1>`, nil
}

func (rs Ressource) showString(c fuego.Ctx[any]) (string, error) {
	return `<h1>test</h1>`, nil
}
