-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4
    )
    RETURNING *
)
SELECT 
    inserted_feed_follow.*,
    f.name AS feed_name,
    u.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds f ON f.id = inserted_feed_follow.feed_id
INNER JOIN users u ON u.id = inserted_feed_follow.user_id;

-- name: GetFeetFollowsForUser :many
SELECT f_f.*, f.name AS feed_name, u.* 
FROM feed_follows f_f
INNER JOIN feeds f ON f_f.feed_id = f.id
INNER JOIN users u ON f_f.user_id = u.id
WHERE f_f.user_id = $1;
