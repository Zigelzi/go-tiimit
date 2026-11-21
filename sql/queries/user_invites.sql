-- name: CreateUserInvite :one
INSERT INTO
    user_invites (
        token,
        created_by_user_id,
        created_at,
        expires_at
    )
VALUES
    (?, ?, ?, ?)
RETURNING
    id,
    expires_at;

-- name: GetUserInviteStatus :one
SELECT
    id,
    CAST(used_at IS NOT NULL AS INTEGER) AS is_used,
    CAST(expires_at <= sqlc.arg(now) AS INTEGER) AS is_expired
FROM
    user_invites
WHERE
    token = sqlc.arg(token);

-- name: UseUserInvite :execrows
UPDATE user_invites
SET
    used_at = sqlc.arg(now)
WHERE
    id = sqlc.arg(id)
    AND used_at IS NULL
    AND expires_at > sqlc.arg(now);

-- name: SetUserInviteUsedBy :exec
UPDATE user_invites
SET
    used_by_user_id = ?
WHERE
    id = ?;
