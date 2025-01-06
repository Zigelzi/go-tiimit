package db

func CreateTables() {
	createPlayerTable := `
	CREATE TABLE IF NOT EXISTS players (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		myclub_id INTEGER NOT NULL UNIQUE,
		run_power REAL NOT NULL,
		ball_handling REAL NOT NULL,
		is_goalie BOOLEAN NOT NULL DEFAULT FALSE
	);
	`

	_, err := DB.Exec(createPlayerTable)
	if err != nil {
		panic("Could not create players table")
	}
}
