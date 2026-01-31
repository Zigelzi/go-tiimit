-- name: CreateUser :one
INSERT INTO
    users (username, hashed_password, created_at, updated_at)
VALUES
    (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING
    id,
    username;