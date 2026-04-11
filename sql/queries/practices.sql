-- name: CreatePractice :execlastid
INSERT INTO
    practices (date)
VALUES
    (?) RETURNING id;

-- name: AddPlayerToPractice :exec
INSERT INTO
    practice_players (practice_id, player_id, team_number)
VALUES
    (?, ?, ?);

-- name: GetPracticeWithPlayers :many
SELECT
    pr.id as practice_id,
    pr.date,
    pp.team_number,
    pp.has_vest,
    pl.id as player_id,
    pl.myclub_id,
    pl.name,
    pl.run_power,
    pl.ball_handling,
    pl.is_goalie
FROM
    practice_players pp
    LEFT JOIN practices pr ON pr.id = pp.practice_id
    LEFT JOIN players pl ON pl.id = pp.player_id
WHERE
    pr.id = ?;

-- name: GetNewestPractices :many
SELECT
    *
FROM
    practices
ORDER BY
    date DESC
LIMIT
    ?;

-- name: SetPlayerTeam :exec
UPDATE practice_players
SET
    team_number = ?
WHERE
    practice_id = ?
    AND player_id = ?;

-- name: GetPracticePlayer :one
SELECT
    pp.practice_id,
    pp.player_id,
    pp.team_number,
    pp.has_vest,
    pl.myclub_id,
    pl.name,
    pl.run_power,
    pl.ball_handling,
    pl.is_goalie
FROM
    practice_players pp
    INNER JOIN players pl ON pl.id = pp.player_id
WHERE
    pp.practice_id = ?
    AND pp.player_id = ?;

-- name: TogglePracticePlayerVest :exec
UPDATE practice_players
SET
    has_vest = ?
WHERE
    practice_id = ?
    AND player_id = ?;

-- name: GetPracticeTeamPlayers :many
SELECT
    pp.practice_id,
    pp.player_id,
    pp.team_number,
    pp.has_vest,
    pl.myclub_id,
    pl.name,
    pl.run_power,
    pl.ball_handling,
    pl.is_goalie
FROM
    practice_players pp
    INNER JOIN players pl ON pl.id = pp.player_id
WHERE
    pp.practice_id = ?
    AND pp.team_number = ?;