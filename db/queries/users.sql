-- name: CreateUser :one
INSERT INTO users (email, first_name, last_name, role, password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET first_name = $2,
    last_name  = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
