-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id) 
VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds 
ORDER BY last_fetch ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedAsFetched :one
UPDATE feeds 
SET last_fetch = NOW(), updated_at = NOW() 
WHERE id = $1
RETURNING *;

-- name: UpdateFeed :one
UPDATE feeds SET name = $2, url = $3 WHERE user_id = $1 AND id = $4 RETURNING *;

-- name: DeleteFeed :exec
DELETE FROM feeds WHERE id = $1 AND user_id = $2;

-- name: GetFeed :one
SELECT * FROM feeds WHERE id = $1;