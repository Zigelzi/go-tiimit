-- +goose Up
-- +goose StatementBegin
-- Force logout and have accurate timestamps
DELETE FROM user_sessions;

ALTER TABLE user_sessions
ADD COLUMN created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE user_sessions
ADD COLUMN deleted_at DATETIME;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_sessions
DROP COLUMN created_at;

ALTER TABLE user_sessions
DROP COLUMN deleted_at;

-- +goose StatementEnd
