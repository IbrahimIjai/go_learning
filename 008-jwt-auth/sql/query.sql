-- name: CreateUser :one
INSERT INTO users (
    user_id, first_name, last_name, password, email, phone, role, token, refresh_token
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: CountByEmailOrPhone :one
SELECT count(*) FROM users
WHERE email = $1 OR phone = $2;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE user_id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at;

-- name: UpdateTokens :exec
UPDATE users
SET token = $2, refresh_token = $3, updated_at = now()
WHERE user_id = $1;
