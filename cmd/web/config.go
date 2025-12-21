package main

import (
	"database/sql"

	"github.com/Zigelzi/go-tiimit/internal/db"
)

type webConfig struct {
	queries *db.Queries
	db      *sql.DB
	address string
	env     string
}
