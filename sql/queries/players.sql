-- name: GetAllPlayers :many
SELECT
    id,
    name,
    myclub_id,
    run_power,
    ball_handling,
    is_goalie
FROM
    players;

-- name: GetPlayerById :one
SELECT
    id,
    name,
    myclub_id,
    run_power,
    ball_handling,
    is_goalie
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
    is_goalie
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
    is_goalie;