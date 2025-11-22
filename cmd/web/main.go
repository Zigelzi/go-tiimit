package main

import (
	"log"
	"net/http"

	"github.com/Zigelzi/go-tiimit/internal/db"
)

func main() {
	newDb, err := db.InitDB()
	if err != nil {
		log.Fatalf("initializing database failed: %v", err)
		return
	}
	cfg := webConfig{
		db:      db.New(newDb),
		address: ":8080",
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", cfg.handleIndexPage)

	mux.HandleFunc("POST /api/attendees", cfg.handleSubmitAttendanceList)
	server := http.Server{
		Handler: mux,
		Addr:    cfg.address,
	}
	log.Printf("Starting server on address %s", cfg.address)
	server.ListenAndServe()
}
