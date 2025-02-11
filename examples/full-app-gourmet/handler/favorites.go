package handler

import (
	"context"
	"log/slog"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/store"
)

type FavoriteRepository interface {
	AddFavorite(ctx context.Context, arg store.AddFavoriteParams) (store.UsersRecipesFavorite, error)
	RemoveFavorite(ctx context.Context, arg store.RemoveFavoriteParams) error
	GetFavoritesByUser(ctx context.Context, username string) ([]store.GetFavoritesByUserRow, error)
}

func (rs Resource) addFavorite(c fuego.ContextNoBody) (store.UsersRecipesFavorite, error) {
	payload := store.AddFavoriteParams{
		Username: c.PathParam("username"),
		RecipeID: c.QueryParam("recipeID"),
	}
	slog.Info("adding favorite", "payload", payload)
	fav, err := rs.FavoritesQueries.AddFavorite(c, payload)
	if err != nil {
		slog.Error("error adding favorite", "err", err)
		return store.UsersRecipesFavorite{}, err
	}
	return fav, nil
}

func (rs Resource) removeFavorite(c fuego.ContextNoBody) (any, error) {
	err := rs.FavoritesQueries.RemoveFavorite(c, store.RemoveFavoriteParams{
		Username: c.PathParam("username"),
		RecipeID: c.QueryParam("recipeID"),
	})
	return nil, err
}

func (rs Resource) getMyFavorites(c fuego.ContextNoBody) ([]store.GetFavoritesByUserRow, error) {
	username := "string"
	return rs.FavoritesQueries.GetFavoritesByUser(c, username)
}

func (rs Resource) getFavoritesByUser(c fuego.ContextNoBody) ([]store.GetFavoritesByUserRow, error) {
	return rs.FavoritesQueries.GetFavoritesByUser(c, c.PathParam("username"))
}

type UserFavorite struct {
	Username string `json:"username"`
	RecipeID string `json:"recipe_id"`
}

func (rs Resource) getFavoritesByUserUnsecureSql(c fuego.ContextNoBody) ([]UserFavorite, error) {
	dbConn := store.InitDB("recipe.bad.db")
	defer dbConn.Close()

	badDBRequest := "SELECT username, recipe_id FROM users_recipes_favorites JOIN recipe ON users_recipes_favorites.recipe_id = recipe.id WHERE username = '" + c.PathParam("username") + "'" // nolint:gosec

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
