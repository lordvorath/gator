// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, created_at, updated_at, user_id, feed_id,
(SELECT name FROM users WHERE id = $4) AS username,
(SELECT name FROM feeds WHERE id = $5) AS feedname
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	Username  string
	Feedname  string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i CreateFeedFollowRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.Username,
		&i.Feedname,
	)
	return i, err
}

const deleteFeedFollow = `-- name: DeleteFeedFollow :exec
WITH feedid AS (SELECT id FROM feeds WHERE feeds.url = $1),
    userid AS (SELECT id FROM users WHERE users.name = $2)
DELETE FROM feed_follows
WHERE feed_follows.user_id = (SELECT id FROM userid)
AND feed_follows.feed_id = (SELECT id FROM feedid)
`

type DeleteFeedFollowParams struct {
	Url  string
	Name string
}

func (q *Queries) DeleteFeedFollow(ctx context.Context, arg DeleteFeedFollowParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollow, arg.Url, arg.Name)
	return err
}

const getFeedFollows = `-- name: GetFeedFollows :many
SELECT id, created_at, updated_at, user_id, feed_id FROM feed_follows
`

func (q *Queries) GetFeedFollows(ctx context.Context) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollow
	for rows.Next() {
		var i FeedFollow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
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

const getFeedFollowsForUser = `-- name: GetFeedFollowsForUser :many
SELECT ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id,
        users.name AS username,
        feeds.name AS feedname
    FROM feed_follows ff
    INNER JOIN users ON ff.user_id = users.id
    INNER JOIN feeds ON ff.feed_id = feeds.id
    WHERE ff.user_id = $1
`

type GetFeedFollowsForUserRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	Username  string
	Feedname  string
}

func (q *Queries) GetFeedFollowsForUser(ctx context.Context, userID uuid.UUID) ([]GetFeedFollowsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsForUserRow
	for rows.Next() {
		var i GetFeedFollowsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.Username,
			&i.Feedname,
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

const resetFeedFollows = `-- name: ResetFeedFollows :exec
DELETE FROM feed_follows *
`

func (q *Queries) ResetFeedFollows(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetFeedFollows)
	return err
}
