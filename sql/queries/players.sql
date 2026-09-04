-- name: GetActivePlayers :many
SELECT
    id,
    name,
    myclub_id,
    run_power,
    ball_handling,
    is_goalie,
    inactivated_at
FROM
    players
WHERE
    inactivated_at is null;

-- name: GetInactivePlayers :many
SELECT
    id,
    name,
    myclub_id,
    run_power,
    ball_handling,
    is_goalie,
    inactivated_at
FROM
    players
WHERE
    inactivated_at is not null;

-- name: GetPlayerById :one
SELECT
    id,
    name,
    myclub_id,
    run_power,
    ball_handling,
    is_goalie,
    inactivated_at
FROM
    players
WHERE
    id = ?;

-- name: GetPlayerByMyclubID :one
SELECT
    id,
    name,
    myclub_id,
    run_power,
    ball_handling,
    is_goalie,
    inactivated_at
FROM
    players
WHERE
    myclub_id = ?;

-- name: AddPlayer :exec
INSERT INTO
    players (name, myclub_id, run_power, ball_handling)
VALUES
    (?, ?, ?, ?);

-- name: IsExistingPlayer :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            players
        WHERE
            myclub_id = ?
    );

-- name: UpdatePlayerAttributes :one
UPDATE players
SET
    run_power = ?,
    ball_handling = ?,
    is_goalie = ?
WHERE
    id = ? RETURNING id,
    name,
    myclub_id,
    run_power,
    ball_handling,
    is_goalie,
    inactivated_at;

-- name: SetPlayerInactive :exec
UPDATE players
SET
    inactivated_at = CURRENT_TIMESTAMP
WHERE
    id = ?;

-- name: SetPlayerActive :exec
UPDATE players
SET
    inactivated_at = NULL
WHERE
    id = ?;