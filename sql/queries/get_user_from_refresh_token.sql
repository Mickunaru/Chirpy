-- name: GetUserFromRefreshToken :one
select * from users
WHERE id = (
    SELECT user_id from refresh_tokens
    WHERE token = $1
);