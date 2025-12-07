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

-- name: ToggleGoalieStatus :exec
UPDATE players
SET
    is_goalie = ?
WHERE
    id = ?;

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

-- name: UpdatePlayerRunPower :exec
UPDATE players
SET
    run_power = ?
WHERE
    myclub_id = ?;