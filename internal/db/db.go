package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const maxOpenConnections = 10
const maxIdleConnections = 2

var DB *sql.DB

type Queries struct {
	db *sql.DB
}

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "./internal/db/tiimit.db")

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

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./internal/db/tiimit.db")

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

func New(db *sql.DB) *Queries {
	return &Queries{
		db: db,
	}
}
