-- name: GetChirpsByUserId :many
SELECT * FROM chirps
WHERE user_id = $1;