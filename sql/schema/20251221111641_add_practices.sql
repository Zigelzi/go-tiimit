-- +goose Up
-- +goose StatementBegin
CREATE TABLE practices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE NOT NULL
);

CREATE TABLE practice_players (
    practice_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL,
    team_number INTEGER NOT NULL CHECK(team_number IN (1,2)),
    PRIMARY KEY (practice_id, player_id),
    FOREIGN KEY (practice_id) REFERENCES practices(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE practice_players;
DROP TABLE practices;
-- +goose StatementEnd
