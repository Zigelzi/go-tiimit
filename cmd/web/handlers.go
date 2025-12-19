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

func (cfg *webConfig) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	component := components.PraticePage()
	component.Render(r.Context(), w)
}

func (cfg *webConfig) handleSubmitAttendanceList(w http.ResponseWriter, r *http.Request) {
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

	// Get the attending player MyClubIds from excel
	// Distribute the attending players to two teams

	// Get the possibly attending players MyClubId from excel
	// Add the players possibly attending to the practice

	// Store practice to DB
	// Show the practice to user

	confirmedPlayers := []player.Player{}
	possiblyAttendingPlayers := []player.Player{}
	unknownPlayers := []int{}
	player.SortByScore(confirmedPlayers)
	player.SortByScore(possiblyAttendingPlayers)
	goalies, fieldPlayers := player.GetPreferences(confirmedPlayers)
	team1, team2, err := team.Distribute(goalies, fieldPlayers)

	if err != nil {
		fmt.Println(err)
	}

	newPractice := practice.Practice{}
	err = newPractice.AddTeams(team1, team2)
	if err != nil {
		fmt.Println(err)
	}

	component := components.DistributedTeams(newPractice, possiblyAttendingPlayers, unknownPlayers)
	component.Render(r.Context(), w)
}
