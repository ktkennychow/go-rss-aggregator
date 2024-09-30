-- +goose Up
CREATE TABLE posts(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  title TEXT NOT NULL,
  url TEXT NOT NULL,
  description TEXT NOT NULL,
  published_at TIMESTAMP NOT NULL,
  feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  -- check for uniqueness in case of updated posts
  CONSTRAINT unique_url_published_at UNIQUE (url, published_at)
);

-- +goose Down
DROP TABLE posts;