package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/store"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/templa"
)

type FavoriteRepository interface {
	AddFavorite(ctx context.Context, arg store.AddFavoriteParams) (store.UsersRecipesFavorite, error)
	RemoveFavorite(ctx context.Context, arg store.RemoveFavoriteParams) error
	GetFavoritesByUser(ctx context.Context, username string) ([]store.GetFavoritesByUserRow, error)
	IsFavorite(ctx context.Context, arg store.IsFavoriteParams) (int64, error)
	GetNumberOfFavorite(ctx context.Context, recipeID string) (int64, error)
}

func (rs Resource) addFavorite(c fuego.ContextNoBody) (*store.UsersRecipesFavorite, error) {
	username := c.PathParam("username")

	caller, err := usernameFromContext(c.Context())
	if err != nil {
		return nil, err
	}

	if caller != username {
		return nil, fuego.ForbiddenError{Title: "you can only add your own favorites"}
	}

	payload := store.AddFavoriteParams{
		Username: username,
		RecipeID: c.QueryParam("recipeID"),
	}
	slog.Info("adding favorite", "payload", payload)
	fav, err := rs.FavoritesQueries.AddFavorite(c, payload)
	if err != nil {
		return nil, err
	}

	// if from htmx, return the recipe page
	if c.Header("HX-Request") == "true" {
		http.Redirect(c.Response(), c.Request(), "/recipes/"+payload.RecipeID, http.StatusSeeOther)
		return nil, nil
	}

	return &fav, nil
}

func (rs Resource) removeFavorite(c fuego.ContextNoBody) (any, error) {
	username := c.PathParam("username")

	caller, err := usernameFromContext(c.Context())
	if err != nil {
		return nil, err
	}

	if caller != username {
		return nil, fuego.ForbiddenError{Title: "you can only remove your own favorites"}
	}

	recipeID := c.QueryParam("recipeID")
	err = rs.FavoritesQueries.RemoveFavorite(c, store.RemoveFavoriteParams{
		Username: username,
		RecipeID: recipeID,
	})

	// if from htmx, return the recipe page
	if c.Header("HX-Request") == "true" {
		http.Redirect(c.Response(), c.Request(), "/recipes/"+recipeID, http.StatusSeeOther)
		return nil, nil
	}

	return nil, err
}

func (rs Resource) getMyFavorites(c fuego.ContextNoBody) (*fuego.DataOrTemplate[[]store.GetFavoritesByUserRow], error) {
	caller, err := usernameFromContext(c.Context())
	if err != nil {
		return nil, err
	}

	favorites, err := rs.FavoritesQueries.GetFavoritesByUser(c, caller)
	if err != nil {
		return nil, err
	}

	return fuego.DataOrHTML(
		favorites,
		templa.Favorites(templa.FavoritesProps{
			Username:  caller,
			Favorites: favorites,
		}),
	), nil
}

func (rs Resource) getFavoritesByUser(c fuego.ContextNoBody) ([]store.GetFavoritesByUserRow, error) {
	return rs.FavoritesQueries.GetFavoritesByUser(c, c.PathParam("username"))
}

type UserFavorite struct {
	Username string `json:"username"`
	RecipeID string `json:"recipe_id"`
}

func (rs Resource) getFavoritesByUserUnsecureSql(c fuego.ContextNoBody) ([]UserFavorite, error) {
	dbConn := store.InitDB("data/recipe.bad.db")
	defer dbConn.Close()

	badDBRequest := "SELECT username, recipe_id FROM users_recipes_favorites WHERE username = '" + c.PathParam("username") + "'" // nolint:gosec

	slog.Info("sqlinjection", "request", badDBRequest)

	rows, err := dbConn.Query(badDBRequest)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipeFavorites := make([]UserFavorite, 0)
	for rows.Next() {
		var userFavorite UserFavorite
		if err := rows.Scan(&userFavorite.Username, &userFavorite.RecipeID); err != nil {
			return nil, err
		}

		recipeFavorites = append(recipeFavorites, userFavorite)
	}

	return recipeFavorites, nil
}
