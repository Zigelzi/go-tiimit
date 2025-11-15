package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleSubmitAttendanceList(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("attendace-list")
	if err != nil {
		log.Printf("Error parsing file from form: %v", err)
		return
	}
	defer file.Close()
	fmt.Println(header.Filename)

}
