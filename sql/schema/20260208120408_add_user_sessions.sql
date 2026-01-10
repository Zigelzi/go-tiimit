-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_sessions (
    session_id TEXT PRIMARY KEY UNIQUE,
    user_id INTEGER NOT NULL,
    expires_at DATE NOT NULL
)
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE user_sessions;

-- +goose StatementEnd
