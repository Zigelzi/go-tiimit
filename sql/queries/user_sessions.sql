-- name: StartUserSession :one
INSERT INTO
    user_sessions (session_id, user_id, expires_at)
VALUES
    (?, ?, ?)
RETURNING
    *;

-- name: GetActiveSession :one
SELECT
    *
FROM
    user_sessions
WHERE
    session_id = ?
    AND expires_at > CURRENT_TIMESTAMP
    AND deleted_at IS NULL;

-- name: EndUserSession :exec
UPDATE user_sessions
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    session_id = ?
    AND deleted_at is null;