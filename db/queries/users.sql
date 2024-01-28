-- name: CreateUser :one
INSERT INTO users (
  id, username, email, password
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 and suspended = false
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 and suspended = false
LIMIT 1;