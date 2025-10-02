-- name: CreateFeedFollows :one 
WITH inserted_feed_follows AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
) SELECT 
    inserted_feed_follows.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follows
INNER JOIN feeds ON inserted_feed_follows.feed_id = feeds.id 
INNER JOIN users ON inserted_feed_follows.user_id = users.id;


-- name: GetFeedFollowsForUser :many
WITH feeds_followed_by_user AS (
    SELECT * FROM feed_follows
    WHERE feed_follows.user_id = $1 
) SELECT 
    feeds_followed_by_user.*,
    feeds.name AS feed_name,
    users.name AS user_name 
FROM feeds_followed_by_user
INNER JOIN feeds ON feeds_followed_by_user.feed_id = feeds.id
INNER JOIN users ON feeds_followed_by_user.user_id = users.id;

-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows
WHERE feed_id = $1
AND user_id = $2;