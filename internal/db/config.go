package db

import (
	"database/sql"
	"fmt"

	migrations "github.com/Zigelzi/go-tiimit/sql"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

const maxOpenConnections = 10
const maxIdleConnections = 2

// InitDB opens the SQLite database at dbPath, creating the file if it does
// not exist. The caller decides where the path comes from, so tests can pass
// a temp file and never touch the real database.
func InitDB(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("no database path given")
	}

	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		return nil, fmt.Errorf("could not connect do database: %w", err)
	}

	// Verify that connection is opened successfully.
	if err := db.Ping(); err != nil {
		db.Close()
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
