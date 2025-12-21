package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Zigelzi/go-tiimit/cmd/web/components"
	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/Zigelzi/go-tiimit/internal/file"
	"github.com/Zigelzi/go-tiimit/internal/player"
	"github.com/Zigelzi/go-tiimit/internal/practice"
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

	confirmedRows, err := file.GetAttendanceRowsByStatus(attendanceRows, file.AttendanceIn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("unable to get the confirmed rows in handler: %v", err)
		return
	}

	dbConfirmedPlayers := []db.Player{}
	for _, row := range confirmedRows {
		confirmedDbPlayer, err := cfg.db.GetPlayerByMyclubID(r.Context(), int64(row.PlayerRow.MyclubID))
		if err != nil {
			log.Println(err)
			continue
		}
		dbConfirmedPlayers = append(dbConfirmedPlayers, confirmedDbPlayer)
	}
	confirmedPlayers := []player.Player{}
	for _, dbConfirmedPlayer := range dbConfirmedPlayers {
		confirmedPlayers = append(confirmedPlayers, player.New(
			dbConfirmedPlayer.MyclubID,
			dbConfirmedPlayer.Name,
			dbConfirmedPlayer.RunPower,
			dbConfirmedPlayer.BallHandling,
			dbConfirmedPlayer.IsGoalie,
		))
	}

	unknownRows, err := file.GetAttendanceRowsByStatus(attendanceRows, file.AttendanceUnknown)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("unable to get the possibly attending player rows in handler: %v", err)
		return
	}

	dbUnknownPlayers := []db.Player{}
	for _, row := range unknownRows {
		unknownDbPlayer, err := cfg.db.GetPlayerByMyclubID(r.Context(), int64(row.PlayerRow.MyclubID))
		if err != nil {
			log.Println(err)
			continue
		}
		dbUnknownPlayers = append(dbUnknownPlayers, unknownDbPlayer)
	}

	unknownPlayers := []player.Player{}
	for _, dbUnknownPlayer := range dbUnknownPlayers {
		unknownPlayers = append(unknownPlayers, player.New(
			dbUnknownPlayer.MyclubID,
			dbUnknownPlayer.Name,
			dbUnknownPlayer.RunPower,
			dbUnknownPlayer.BallHandling,
			dbUnknownPlayer.IsGoalie,
		))
	}

	player.SortByScore(confirmedPlayers)
	player.SortByScore(unknownPlayers)
	goalies, fieldPlayers := player.GetPreferences(confirmedPlayers)
	team1, team2, err := practice.Distribute(fieldPlayers, goalies)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("unable to distribute the players in handler: %v", err)
		return
	}

	newPractice := practice.Practice{
		TeamOnePlayers: team1,
		TeamTwoPlayers: team2,
		UnknownPlayers: unknownPlayers,
		Date:           time.Now(),
	}

	if err != nil {
		fmt.Println(err)
	}

	component := components.DistributedTeams(newPractice)
	component.Render(r.Context(), w)
}
