-- name: GetDosings :many
SELECT * FROM dosing;


-- name: CreateDosing :one
INSERT INTO dosing (recipe_id, ingredient_id, quantity, unit) VALUES (?, ?, ?, ?) RETURNING *;
