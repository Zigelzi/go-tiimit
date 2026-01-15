package main

import (
	"database/sql"

	"github.com/Zigelzi/go-tiimit/internal/db"
)

type cliConfig struct {
	queries *db.Queries
	db      *sql.DB
}
