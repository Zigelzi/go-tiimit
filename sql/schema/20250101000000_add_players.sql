-- +goose Up
-- +goose StatementBegin
CREATE TABLE players (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		myclub_id INTEGER NOT NULL UNIQUE,
		run_power REAL NOT NULL,
		ball_handling REAL NOT NULL
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE players;
-- +goose StatementEnd
