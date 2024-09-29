-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES ($1, NOW(), NOW(), $2, $3, $4, null)
RETURNING *;

-- name: ReadFeeds :many
SELECT *
FROM feeds;

-- name: ReadNFeedsByLastFetchedAt :many
SELECT *
FROM feeds
ORDER BY last_fetched_at IS NULL, last_fetched_at ASC
LIMIT $1;

-- name: UpdateFeedFetched :exec
UPDATE feeds
SET updated_at = NOW(), last_fetched_at = NOW()
WHERE id = $1;