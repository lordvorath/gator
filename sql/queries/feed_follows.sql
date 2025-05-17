-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *,
(SELECT name FROM users WHERE id = $4) AS username,
(SELECT name FROM feeds WHERE id = $5) AS feedname;

-- name: GetFeedFollowsForUser :many
SELECT ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id,
        users.name AS username,
        feeds.name AS feedname
    FROM feed_follows ff
    INNER JOIN users ON ff.user_id = users.id
    INNER JOIN feeds ON ff.feed_id = feeds.id
    WHERE ff.user_id = $1;

-- name: GetFeedFollows :many
SELECT * FROM feed_follows;

-- name: ResetFeedFollows :exec
DELETE FROM feed_follows *;