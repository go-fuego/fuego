-- name: GetRecipe :one
SELECT * FROM recipe WHERE id = ?;

-- name: GetRecipes :many
SELECT * FROM recipe;

-- name: SearchRecipes :many
-- Saerch anything that contains the given string
SELECT * FROM recipe WHERE name LIKE ?;

-- name: CreateRecipe :one
INSERT INTO recipe (id, name, description, instructions) VALUES (?, ?, ?, ?) RETURNING *;

-- name: DeleteRecipe :exec
DELETE FROM recipe WHERE id = ?;

-- name: GetRandomRecipes :many
SELECT * FROM recipe ORDER BY RANDOM() DESC LIMIT 10;

-- name: UpdateRecipe :one
UPDATE recipe SET 
  name=COALESCE(sqlc.arg(name), name),
  description=COALESCE(sqlc.narg(description), description),
  instructions=COALESCE(sqlc.narg(instructions), instructions),
  category=COALESCE(sqlc.arg(category), category),
  class=COALESCE(sqlc.arg(class), class),
  image_url=COALESCE(sqlc.arg(image_url), image_url),
  cook_time=COALESCE(sqlc.arg(cook_time), cook_time),
  prep_time=COALESCE(sqlc.arg(prep_time), prep_time),
  servings=COALESCE(sqlc.arg(servings), servings),
  published=COALESCE(sqlc.arg(published), published)
WHERE id = @id
RETURNING *;

