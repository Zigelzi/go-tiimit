package main

import (
	"github.com/Zigelzi/go-tiimit/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

type cliConfig struct {
	db *db.Queries
}
