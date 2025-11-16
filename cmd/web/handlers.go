package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zigelzi/go-tiimit/cmd/web/components"
	"github.com/Zigelzi/go-tiimit/internal/file"
)

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(201)
	component := components.Index()
	component.Render(r.Context(), w)
}

func handleSubmitAttendanceList(w http.ResponseWriter, r *http.Request) {
	formFile, header, err := r.FormFile("attendace-list")
	if err != nil {
		log.Printf("Error parsing file from form: %v", err)
		return
	}
	defer formFile.Close()
	fmt.Println(header.Filename)
	attendanceRows, err := file.ImportAttendancePlayerRowsFromReader(formFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("unable to parse the attendance rows in handler: %v", err)
		return
	}
	for _, row := range attendanceRows {
		fmt.Println(row)
	}
}
