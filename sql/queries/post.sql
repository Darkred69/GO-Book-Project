-- name: CreatePost :one
INSERT INTO posts (id, title, description, published_at, url, feed_id)
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING *;

-- name: GetPosts :many
SELECT posts.* FROM posts
JOIN feed_follow ON posts.feed_id = feed_follow.feed_id
WHERE feed_follow.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;