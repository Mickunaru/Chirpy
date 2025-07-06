-- name: DeleteChirpFromId :exec
DELETE FROM chirps
WHERE id = $1;