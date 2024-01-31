-- name: GetIngredient :one
SELECT * FROM ingredient WHERE id = ?;

-- name: GetIngredients :many
SELECT * FROM ingredient;

-- name: SearchIngredients :many
SELECT * FROM ingredient
WHERE name LIKE ?
ORDER BY name ASC
LIMIT ?
OFFSET ?;

-- name: GetIngredientsOfRecipe :many
SELECT quantity, unit, sqlc.embed(ingredient) FROM ingredient
JOIN dosing ON ingredient.id = dosing.ingredient_id
WHERE dosing.recipe_id = ?;

-- name: CreateIngredient :one
INSERT INTO ingredient 
(id, name, description, available_all_year, available_jan, available_feb, available_mar, available_apr, available_may, available_jun, available_jul, available_aug, available_sep, available_oct, available_nov, available_dec, category, default_unit)
VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
RETURNING *;

-- name: UpdateIngredient :one
UPDATE ingredient SET 
  name=COALESCE(sqlc.arg(name), name),
  description=COALESCE(sqlc.arg(description), description),
  category=COALESCE(sqlc.narg(category), category),
  default_unit=COALESCE(sqlc.narg(default_unit), default_unit),
  available_all_year=COALESCE(sqlc.arg(available_all_year), available_all_year),
  available_jan=COALESCE(sqlc.arg(available_jan), available_jan),
  available_feb=COALESCE(sqlc.arg(available_feb), available_feb),
  available_mar=COALESCE(sqlc.arg(available_mar), available_mar),
  available_apr=COALESCE(sqlc.arg(available_apr), available_apr),
  available_may=COALESCE(sqlc.arg(available_may), available_may),
  available_jun=COALESCE(sqlc.arg(available_jun), available_jun),
  available_jul=COALESCE(sqlc.arg(available_jul), available_jul),
  available_aug=COALESCE(sqlc.arg(available_aug), available_aug),
  available_sep=COALESCE(sqlc.arg(available_sep), available_sep),
  available_oct=COALESCE(sqlc.arg(available_oct), available_oct),
  available_nov=COALESCE(sqlc.arg(available_nov), available_nov),
  available_dec=COALESCE(sqlc.arg(available_dec), available_dec)
WHERE id = @id
RETURNING *;
