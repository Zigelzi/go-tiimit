package main

import (
	"database/sql"
	"time"

	"github.com/Zigelzi/go-tiimit/internal/db"
)

type webConfig struct {
	queries  *db.Queries
	db       *sql.DB
	address  string
	env      string
	phApiKey string
	// now returns the current time. Tests replace it to fix the clock.
	now func() time.Time
}
