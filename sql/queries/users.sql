-- name: CreateUser :one
INSERT INTO
    users (username, hashed_password, created_at, updated_at)
VALUES
    (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING
    id,
    username;

-- name: IsValidPassword :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            users
        WHERE
            username = ?
            AND hashed_password = ?
    ) AS is_valid_password;