package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zigelzi/go-tiimit/cmd/web/components"
	"github.com/Zigelzi/go-tiimit/internal/file"
	"github.com/Zigelzi/go-tiimit/internal/player"
	"github.com/Zigelzi/go-tiimit/internal/practice"
	"github.com/Zigelzi/go-tiimit/internal/team"
)

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
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
	fmt.Printf("parsed %d rows from attendance excel\n", len(attendanceRows))

	newPractice := practice.New()

	team1 := team.Team{
		Name: "Team 1",
		Players: []player.Player{
			player.New(1234, "Matti Meikäläinen", 5.0, 5.0, false),
			player.New(2345, "Teppo Teikäläinen", 5.0, 5.0, false),
		},
	}
	team2 := team.Team{
		Name: "Team 1",
		Players: []player.Player{
			player.New(3456, "Heikki Heikäläinen", 5.0, 5.0, false),
			player.New(4567, "Seppo Seikäläinen", 5.0, 5.0, false),
		},
	}
	newPractice.AddTeams(team1, team2)

	component := components.DistributedTeams(newPractice)
	component.Render(r.Context(), w)
}
