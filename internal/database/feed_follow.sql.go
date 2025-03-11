// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: feed_follow.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4
    )
    RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT 
    inserted_feed_follow.id, inserted_feed_follow.created_at, inserted_feed_follow.updated_at, inserted_feed_follow.user_id, inserted_feed_follow.feed_id,
    f.name AS feed_name,
    u.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds f ON f.id = inserted_feed_follow.feed_id
INNER JOIN users u ON u.id = inserted_feed_follow.user_id
`

type CreateFeedFollowParams struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    int32
}

type CreateFeedFollowRow struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    int32
	FeedName  string
	UserName  string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
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
		&i.FeedName,
		&i.UserName,
	)
	return i, err
}

const getFeetFollowsForUser = `-- name: GetFeetFollowsForUser :many
SELECT f_f.id, f_f.created_at, f_f.updated_at, f_f.user_id, f_f.feed_id, f.name AS feed_name, u.id, u.created_at, u.updated_at, u.name 
FROM feed_follows f_f
INNER JOIN feeds f ON f_f.feed_id = f.id
INNER JOIN users u ON f_f.user_id = u.id
WHERE f_f.user_id = $1
`

type GetFeetFollowsForUserRow struct {
	ID          int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      uuid.UUID
	FeedID      int32
	FeedName    string
	ID_2        uuid.UUID
	CreatedAt_2 time.Time
	UpdatedAt_2 time.Time
	Name        string
}

func (q *Queries) GetFeetFollowsForUser(ctx context.Context, userID uuid.UUID) ([]GetFeetFollowsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeetFollowsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeetFollowsForUserRow
	for rows.Next() {
		var i GetFeetFollowsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.FeedName,
			&i.ID_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.Name,
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
