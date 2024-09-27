-- name: CreateUserFeed :one
INSERT INTO users_feeds (id, created_at, updated_at, feed_id, user_id)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING *;

-- name: ReadUserFeeds :many
SELECT *
FROM users_feeds
WHERE user_id = $1;

-- name: DeleteUserFeed :exec
DELETE
FROM users_feeds
WHERE id = $1
AND user_id = $2;