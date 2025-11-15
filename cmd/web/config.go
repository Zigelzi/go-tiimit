package main

import "github.com/Zigelzi/go-tiimit/internal/db"

type webConfig struct {
	db      *db.Queries
	address string
	env     string
}
