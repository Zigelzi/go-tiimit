-- name: StartUserSession :one
INSERT INTO
    user_sessions (session_id, user_id, expires_at)
VALUES
    (?, ?, ?) RETURNING *;

-- name: GetActiveSession :one
SELECT
    us.user_id,
    u.username
FROM
    user_sessions us
    JOIN users u ON u.id = us.user_id
WHERE
    us.session_id = ?
    AND us.expires_at > CURRENT_TIMESTAMP
    AND us.deleted_at IS NULL;

-- name: EndUserSession :exec
UPDATE user_sessions
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    session_id = ?
    AND deleted_at is null;