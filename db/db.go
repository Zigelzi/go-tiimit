package db

import (
	"database/sql"
	"fmt"

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

	// Verify that connection is opened successfully.
	if err := DB.Ping(); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
	}

	DB.SetMaxOpenConns(maxOpenConnections)
	DB.SetMaxIdleConns(maxIdleConnections)
}
