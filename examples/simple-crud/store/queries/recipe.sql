-- name: GetRecipe :one
SELECT * FROM recipe WHERE id = ?;

-- name: GetRecipeWithIngredients :one
SELECT * FROM recipe
JOIN dosing ON recipe.id = dosing.recipe_id
JOIN ingredient ON dosing.ingredient_id = ingredient.id
WHERE recipe.id = ?;

-- name: GetRecipes :many
SELECT * FROM recipe;

-- name: SearchRecipes :many
-- Saerch anything that contains the given string
SELECT * FROM recipe WHERE name LIKE ?;

-- name: CreateRecipe :one
INSERT INTO recipe (id, name, description) VALUES (?, ?, ?) RETURNING *;

-- name: DeleteRecipe :exec
DELETE FROM recipe WHERE id = ?;
