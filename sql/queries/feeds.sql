-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES ($1, $2, $3, $4, $5, $6, null)
RETURNING *;

-- name: ReadFeeds :many
SELECT *
FROM feeds;