package main

import (
	"log"
	"net/http"
)

func main() {
	address := ":8080"

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleIndexPage)

	mux.HandleFunc("POST /api/attendees", handleSubmitAttendanceList)
	server := http.Server{
		Handler: mux,
		Addr:    address,
	}
	log.Printf("Starting server on address %s", address)
	server.ListenAndServe()
}
