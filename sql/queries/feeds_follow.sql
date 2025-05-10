-- name: CreateFollow :one
INSERT INTO feed_follow (id, user_id, feed_id) 
VALUES ($1, $2, $3) 
RETURNING *;

-- name: GetFollows :many
SELECT * FROM feed_follow WHERE user_id = $1;

-- name: Unfollow :exec
DELETE FROM feed_follow WHERE user_id = $1 AND feed_id = $2;

-- name: GetFollowsByFeedID :one
SELECT * FROM feed_follow WHERE feed_id = $1 AND user_id = $2;