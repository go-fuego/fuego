-- name: CreateUser :one
INSERT INTO users (username, full_name, email, encrypted_password) VALUES (?, ?, ?, ?) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;

-- name: GetUsers :many
SELECT * FROM users;
