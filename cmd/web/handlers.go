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

	newPractice := practice.New()
	for _, row := range attendanceRows {
		err := newPractice.AddPlayer(row.PlayerRow.MyClubId, row.Attendance)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	dbConfirmedPlayers, unknownPlayers, err := newPractice.GetPlayersByStatus(practice.AttendanceIn, cfg.db.Get)
	if err != nil {
		fmt.Println(err)
	}
	confirmedPlayers := []player.Player{}
	for _, dbConfirmedPlayer := range dbConfirmedPlayers {
		confirmedPlayers = append(confirmedPlayers, player.New(dbConfirmedPlayer.MyClubId, dbConfirmedPlayer.Name, dbConfirmedPlayer.RunPower, dbConfirmedPlayer.BallHandling, dbConfirmedPlayer.IsGoalie))
	}

	dbPossiblyAttendingPlayers, unknownPossiblePlayers, err := newPractice.GetPlayersByStatus(practice.AttendanceUnknown, cfg.db.Get)
	if err != nil {
		fmt.Println(err)
	}

	possiblyAttendingPlayers := []player.Player{}
	for _, dbPossiblyAttendingPlayer := range dbPossiblyAttendingPlayers {
		possiblyAttendingPlayers = append(possiblyAttendingPlayers, player.New(dbPossiblyAttendingPlayer.MyClubId, dbPossiblyAttendingPlayer.Name, dbPossiblyAttendingPlayer.RunPower, dbPossiblyAttendingPlayer.BallHandling, dbPossiblyAttendingPlayer.IsGoalie))
	}

	for myClubId, attendanceStatus := range unknownPossiblePlayers {
		unknownPlayers[myClubId] = attendanceStatus
	}

	player.SortByScore(confirmedPlayers)
	player.SortByScore(possiblyAttendingPlayers)
	goalies, fieldPlayers := player.GetPreferences(confirmedPlayers)
	team1, team2, err := team.Distribute(goalies, fieldPlayers)

	if err != nil {
		fmt.Println(err)

	}

	err = newPractice.AddTeams(team1, team2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(unknownPlayers)

	component := components.DistributedTeams(newPractice, possiblyAttendingPlayers, unknownPlayers)
	component.Render(r.Context(), w)
}
