-- name: GetRecipe :one
SELECT * FROM recipe WHERE id = ?;

-- name: GetRecipes :many
SELECT * FROM recipe LIMIT @limit;

-- name: SearchRecipes :many
-- Search anything that contains the given string
SELECT * FROM recipe WHERE
  (name LIKE '%' || @search || '%')
  AND published = @published
  AND calories <= @max_calories
  AND prep_time + cook_time <= @max_time
ORDER BY name ASC
LIMIT @limit
OFFSET @offset;

-- name: CreateRecipe :one
INSERT INTO recipe (
  id,
  name,
  description,
  instructions,
  prep_time,
  cook_time,
  category,
  image_url,
  published,
  servings,
  when_to_eat
) 
VALUES (?,?,?,?,?,?,?,?,?,?,?) RETURNING *;

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
  when_to_eat=COALESCE(sqlc.arg(when_to_eat), when_to_eat),
  image_url=COALESCE(sqlc.arg(image_url), image_url),
  cook_time=COALESCE(sqlc.arg(cook_time), cook_time),
  prep_time=COALESCE(sqlc.arg(prep_time), prep_time),
  servings=COALESCE(sqlc.arg(servings), servings),
  published=COALESCE(sqlc.arg(published), published)
WHERE id = @id
RETURNING *;

