-- name: CreatePractice :execlastid
INSERT INTO practices (date)
VALUES (?)
RETURNING id;

-- name: AddPlayerToPractice :exec
INSERT INTO practice_players (practice_id, player_id, team_number)
VALUES (?,?,?);

-- name: GetPracticeWithPlayers :many
SELECT 
    pr.id as practice_id,
    pr.date,
    pp.team_number,
    pl.myclub_id,
    pl.name,
    pl.run_power,
    pl.ball_handling,
    pl.is_goalie
FROM practice_players pp
LEFT JOIN practices pr ON pr.id=pp.practice_id
LEFT JOIN players pl ON pl.id=pp.player_id
WHERE pr.id=?;