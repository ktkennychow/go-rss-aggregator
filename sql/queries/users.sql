-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, apikey)
VALUES ($1, NOW(), NOW(), $2, encode(sha256(random()::text::bytea), 'hex'))
RETURNING *;

-- name: ReadUser :one
SELECT *
FROM users
WHERE apikey = $1;
