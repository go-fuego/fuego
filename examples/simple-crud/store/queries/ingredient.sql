-- name: GetIngredient :one
SELECT * FROM ingredient WHERE id = ?;

-- name: GetIngredients :many
SELECT * FROM ingredient;

-- name: GetIngredientsOfRecipe :many
SELECT quantity, sqlc.embed(ingredient) FROM ingredient
JOIN dosing ON ingredient.id = dosing.ingredient_id
WHERE dosing.recipe_id = ?;

-- name: CreateIngredient :one
INSERT INTO ingredient (id, name, description) VALUES (?, ?, ?) RETURNING *;
