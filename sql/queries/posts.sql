-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, Now(), Now(), $2, $3, $4, $5, $6)
RETURNING *;

-- name: ReadPostsByUser :many
SELECT *
FROM posts
WHERE feed_id IN (
  SELECT feed_id
  FROM users_feeds
  WHERE user_id = $1
)
ORDER BY published_at DESC 
LIMIT $2;