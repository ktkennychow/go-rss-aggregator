// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feeds.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES ($1, NOW(), NOW(), $2, $3, $4, null)
RETURNING id, created_at, updated_at, name, url, user_id, last_fetched_at
`

type CreateFeedParams struct {
	ID     uuid.UUID
	Name   string
	Url    string
	UserID uuid.UUID
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const readFeeds = `-- name: ReadFeeds :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at
FROM feeds
`

func (q *Queries) ReadFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, readFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.LastFetchedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readNFeedsByLastFetchedAt = `-- name: ReadNFeedsByLastFetchedAt :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at
FROM feeds
ORDER BY last_fetched_at IS NULL, last_fetched_at ASC
LIMIT $1
`

func (q *Queries) ReadNFeedsByLastFetchedAt(ctx context.Context, limit int32) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, readNFeedsByLastFetchedAt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.LastFetchedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateFeedFetched = `-- name: UpdateFeedFetched :exec
UPDATE feeds
SET updated_at = NOW(), last_fetched_at = NOW()
WHERE id = $1
`

func (q *Queries) UpdateFeedFetched(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, updateFeedFetched, id)
	return err
}
