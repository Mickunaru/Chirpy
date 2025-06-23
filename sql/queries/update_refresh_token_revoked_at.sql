-- name: UpdateRefreshTokenRevokedAt :exec
UPDATE refresh_tokens
SET revoked_at = $1,
    updated_at = NOW()
WHERE token = $2;