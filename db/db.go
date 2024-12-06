package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const maxOpenConnections = 10
const maxIdleConnections = 2

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "./db/tiimit.db")

	if err != nil {
		panic("Could not connect do database")
	}

	DB.SetMaxOpenConns(maxOpenConnections)
	DB.SetMaxIdleConns(maxIdleConnections)
	createTables()
}

func createTables() {
	createPlayersTable := `
	CREATE TABLE IF NOT EXISTS players (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		myclub_id INTEGER NOT NULL UNIQUE,
		run_power REAL NOT NULL,
		ball_handling REAL NOT NULL
	);
	`

	_, err := DB.Exec(createPlayersTable)
	if err != nil {
		panic("Could not create players table")
	}
}
