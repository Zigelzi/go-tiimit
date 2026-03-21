-- +goose Up
ALTER TABLE practice_players
ADD COLUMN has_vest BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE practice_players
DROP COLUMN has_vest;