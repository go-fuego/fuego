package views_test

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-fuego/fuego/examples/full-app-gourmet/server"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/store"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/views"
	"github.com/stretchr/testify/require"
)

func TestShowIndex(t *testing.T) {
	// rs := views.Ressource{
	// 	RecipesQueries:     nil,
	// 	IngredientsQueries: nil,
	// 	DosingQueries:      nil,
	// }

	// c := &fuego.Context[any]{}

	// rs.showIndex(c)
}

type RecipeRepositoryMock struct {
	views.RecipeRepository
}

func (r RecipeRepositoryMock) GetRecipes(ctx context.Context) ([]store.Recipe, error) {
	time.Sleep(1 * time.Millisecond)
	return []store.Recipe{}, nil
}

func (r RecipeRepositoryMock) SearchRecipes(ctx context.Context, params store.SearchRecipesParams) ([]store.Recipe, error) {
	return []store.Recipe{}, nil
}

func (r RecipeRepositoryMock) GetRandomRecipes(ctx context.Context) ([]store.Recipe, error) {
	return []store.Recipe{}, nil
}

type IngredientRepositoryMock struct {
	views.IngredientRepository
}

func (r IngredientRepositoryMock) GetIngredients(ctx context.Context) ([]store.Ingredient, error) {
	time.Sleep(1 * time.Millisecond)
	return []store.Ingredient{}, nil
}

func TestShowIndexExt(t *testing.T) {
	viewsRessources := views.Ressource{
		RecipesQueries:     RecipeRepositoryMock{},
		IngredientsQueries: IngredientRepositoryMock{},
	}

	serverRessources := server.Ressources{
		Views: viewsRessources,
	}

	os.Setenv("ADMIN_USER", "admin")
	os.Setenv("ADMIN_PASSWORD", "admin")

	app := serverRessources.Setup()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	app.Mux.ServeHTTP(w, r)

	require.Equal(t, 200, w.Code)
}

func BenchmarkShowIndexExt(b *testing.B) {
	viewsRessources := views.Ressource{
		RecipesQueries:     RecipeRepositoryMock{},
		IngredientsQueries: IngredientRepositoryMock{},
	}

	serverRessources := server.Ressources{
		Views: viewsRessources,
	}

	os.Setenv("ADMIN_USER", "admin")
	os.Setenv("ADMIN_PASSWORD", "admin")

	app := serverRessources.Setup()

	for range b.N {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		app.Mux.ServeHTTP(w, r)

		if w.Code != 200 {
			b.Fail()
		}
	}
}
