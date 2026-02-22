package db

import (
	"database/sql"
	"fmt"
	"os"

	migrations "github.com/Zigelzi/go-tiimit/sql"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

const maxOpenConnections = 10
const maxIdleConnections = 2

func InitDB() (*sql.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		return nil, fmt.Errorf("DB_PATH environment variable isn't set")
	}

	db, err := sql.Open("sqlite", dbPath)

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

func RunMigrations(db *sql.DB) error {
	migrationsFs := migrations.GetMigrationFS()
	goose.SetBaseFS(migrationsFs)
	err := goose.SetDialect("sqlite3")
	if err != nil {
		return fmt.Errorf("failed to set the Goose dialect: %v", err)
	}

	// Directory set and copied in the migrationFs
	// See https://github.com/pressly/goose#embedded-sql-migrations for details
	err = goose.Up(db, "schema")
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	return nil
}
