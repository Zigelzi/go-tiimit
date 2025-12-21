package main

import (
	"database/sql"

	"github.com/Zigelzi/go-tiimit/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

type cliConfig struct {
	queries *db.Queries
	db      *sql.DB
}
