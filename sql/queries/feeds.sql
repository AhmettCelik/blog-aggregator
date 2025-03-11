-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetFeeds :many
SELECT f.name, f.url, u.name AS user_name
FROM feeds f
INNER JOIN users u
ON f.user_id = u.id;

-- name: GetFeedForUrl :one
SELECT id, name, url, user_id FROM feeds WHERE url = $1;
