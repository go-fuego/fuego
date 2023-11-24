package views

import (
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"
	"path"

	"simple-crud/store/dosings"
	"simple-crud/store/ingredients"
	"simple-crud/store/recipes"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

func NewRessource(db *sql.DB) Ressource {
	return Ressource{
		RecipesQueries:     *recipes.New(db),
		IngredientsQueries: *ingredients.New(db),
		DosingQueries:      *dosings.New(db),
	}
}

// Ressource is the struct that holds useful sources of informations available for the controllers.
type Ressource struct {
	DosingQueries      dosings.Queries
	RecipesQueries     recipes.Queries
	IngredientsQueries ingredients.Queries
}

func (rs Ressource) showRecipesStd(w http.ResponseWriter, r *http.Request) {
	recipes, err := rs.RecipesQueries.GetRecipes(r.Context())
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
	recipes, err := rs.RecipesQueries.GetRecipes(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/recipes.page.html", recipes)
}

func (rs Ressource) searchRecipes(c fuego.Ctx[any]) (fuego.HTML, error) {
	search := c.QueryParam("q")

	recipes, err := rs.RecipesQueries.SearchRecipes(c.Context(), "%"+search+"%")
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
	recipes, err := rs.RecipesQueries.SearchRecipes(c.Context(), "%"+c.QueryParam("search")+"%")
	if err != nil {
		return "", err
	}

	slog.Debug("recipes", "recipes", recipes, "search", c.QueryParam("search"))

	return c.Render("partials/recipes-list.partial.html", recipes)
}

func (rs Ressource) addRecipe(c fuego.Ctx[recipes.CreateRecipeParams]) (fuego.HTML, error) {
	body, err := c.Body()
	if err != nil {
		return "", err
	}

	body.ID = uuid.NewString()

	_, err = rs.RecipesQueries.CreateRecipe(c.Context(), body)
	if err != nil {
		return "", err
	}

	recipes, err := rs.RecipesQueries.GetRecipes(c.Context())
	if err != nil {
		return "", err
	}

	return c.Render("pages/admin.page.html", fuego.H{
		"Recipes": recipes,
	})
}

func (rs Ressource) showHTML(c fuego.Ctx[any]) (fuego.HTML, error) {
	return `<h1>test</h1>`, nil
}

func (rs Ressource) showString(c fuego.Ctx[any]) (string, error) {
	return `<h1>test</h1>`, nil
}
