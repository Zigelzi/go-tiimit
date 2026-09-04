-- +goose Up
ALTER TABLE players
ADD COLUMN inactivated_at DATETIME;

-- +goose Down
ALTER TABLE players
DROP COLUMN inactivated_at;