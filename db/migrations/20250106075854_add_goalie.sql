-- +goose Up
-- +goose StatementBegin
ALTER TABLE players
ADD COLUMN is_goalie BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE players
DROP COLUMN is_goalie;
-- +goose StatementEnd
