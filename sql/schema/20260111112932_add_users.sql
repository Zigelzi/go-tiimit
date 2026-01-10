-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;

-- +goose StatementEnd
