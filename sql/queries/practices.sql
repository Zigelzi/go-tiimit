-- name: CreatePractice :execlastid
INSERT INTO practices (date)
VALUES (?)
RETURNING id;

-- name: AddPlayerToPractice :exec
INSERT INTO practice_players (practice_id, player_id, team_number)
VALUES (?,?,?);