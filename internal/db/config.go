package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const maxOpenConnections = 10
const maxIdleConnections = 2

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./internal/db/tiimit.db")

	if err != nil {
		return nil, fmt.Errorf("could not connect do database: %w", err)
	}

	// Verify that connection is opened successfully.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	return db, nil
}
