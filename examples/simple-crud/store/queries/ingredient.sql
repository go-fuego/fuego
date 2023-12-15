-- name: GetIngredient :one
SELECT * FROM ingredient WHERE id = ?;

-- name: GetIngredients :many
SELECT * FROM ingredient;

-- name: GetIngredientsOfRecipe :many
SELECT quantity, unit, sqlc.embed(ingredient) FROM ingredient
JOIN dosing ON ingredient.id = dosing.ingredient_id
WHERE dosing.recipe_id = ?;

-- name: CreateIngredient :one
INSERT INTO ingredient (id, name, description) VALUES (?, ?, ?) RETURNING *;

-- name: UpdateIngredient :one
UPDATE ingredient SET 
  name=COALESCE(sqlc.arg(name), name),
  category=COALESCE(sqlc.narg(category), category),
  default_unit=COALESCE(sqlc.narg(default_unit), default_unit)
WHERE id = @id
RETURNING *;
