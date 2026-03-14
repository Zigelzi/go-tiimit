package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Zigelzi/go-tiimit/cmd/web/components"
	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/Zigelzi/go-tiimit/internal/file"
	"github.com/Zigelzi/go-tiimit/internal/player"
	"github.com/Zigelzi/go-tiimit/internal/practice"
)

func (cfg *webConfig) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	dbPractices, err := cfg.queries.GetNewestPractices(r.Context(), 5)
	if err != nil {
		log.Printf("failed to get practices: %v", err)
	}
	practices := []practice.Practice{}
	for _, dbPractice := range dbPractices {
		practices = append(practices, practice.FromDB(dbPractice))
	}
	component := components.IndexPage(practices)
	component.Render(r.Context(), w)
}

func (cfg *webConfig) handleSetupPracticePage(w http.ResponseWriter, r *http.Request) {
	component := components.CreatePracticePage()
	component.Render(r.Context(), w)
}

func (cfg *webConfig) handleCreatePractice(w http.ResponseWriter, r *http.Request) {
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
		confirmedDbPlayer, err := cfg.queries.GetPlayerByMyclubID(r.Context(), int64(row.PlayerRow.MyclubID))
		if err != nil {
			log.Println(err)
			continue
		}
		dbConfirmedPlayers = append(dbConfirmedPlayers, confirmedDbPlayer)
	}

	confirmedPlayers := []player.Player{}
	for _, dbConfirmedPlayer := range dbConfirmedPlayers {
		confirmedPlayers = append(confirmedPlayers, player.FromDB(dbConfirmedPlayer))
	}

	unknownRows, err := file.GetAttendanceRowsByStatus(attendanceRows, file.AttendanceUnknown)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("unable to get the possibly attending player rows in handler: %v", err)
		return
	}

	dbUnknownPlayers := []db.Player{}
	for _, row := range unknownRows {
		unknownDbPlayer, err := cfg.queries.GetPlayerByMyclubID(r.Context(), int64(row.PlayerRow.MyclubID))
		if err != nil {
			log.Println(err)
			continue
		}
		dbUnknownPlayers = append(dbUnknownPlayers, unknownDbPlayer)
	}

	unknownPlayers := []player.Player{}
	for _, dbUnknownPlayer := range dbUnknownPlayers {
		unknownPlayers = append(unknownPlayers, player.FromDB(dbUnknownPlayer))
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

	practiceDate, err := file.ParseDate(header.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to parse date from the file name: %v", err)
		return
	}

	newPractice := practice.Practice{
		TeamOnePlayers: team1,
		TeamTwoPlayers: team2,
		UnknownPlayers: unknownPlayers,
		Date:           practiceDate,
	}

	if err != nil {
		fmt.Println(err)
	}

	// REPO START - Storing the data.
	tx, err := cfg.db.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to begin a transaction: %v", err)
		return
	}
	defer tx.Rollback()

	queryTx := cfg.queries.WithTx(tx)
	dbPracticeId, err := queryTx.CreatePractice(r.Context(), newPractice.Date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to create the practice: %v", err)
		return
	}
	for _, teamOnePlayer := range newPractice.TeamOnePlayers {
		err = queryTx.AddPlayerToPractice(r.Context(), db.AddPlayerToPracticeParams{
			PracticeID: dbPracticeId,
			PlayerID:   teamOnePlayer.ID,
			TeamNumber: 1,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to add player [%d] %s to practice %d in team 1", teamOnePlayer.ID, teamOnePlayer.Name, dbPracticeId)
			return
		}
	}

	for _, teamTwoPlayer := range newPractice.TeamTwoPlayers {
		err = queryTx.AddPlayerToPractice(r.Context(), db.AddPlayerToPracticeParams{
			PracticeID: dbPracticeId,
			PlayerID:   teamTwoPlayer.ID,
			TeamNumber: 2,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to add player [%d] %s to practice %d in team 2", teamTwoPlayer.ID, teamTwoPlayer.Name, dbPracticeId)
			return
		}
	}
	tx.Commit()
	// REPO END
	w.Header().Add("HX-Redirect", fmt.Sprintf("/practices/%d", dbPracticeId))
}

func (cfg *webConfig) handleViewPractice(w http.ResponseWriter, r *http.Request) {

	practiceId, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("unable to parse practice id from path: %v", err)
		return
	}
	dbPracticeRows, err := cfg.queries.GetPracticeWithPlayers(r.Context(), int64(practiceId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to get practice from database: %v", err)
		return
	}
	currentPractice, err := practice.FromDBWithPlayers(dbPracticeRows)
	if err != nil {
		if errors.Is(err, practice.ErrNoPracticeRows) {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("user tried to view practice with ID [%d] which doesn't exist", practiceId)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to convert the practice players: %v", err)
		return
	}
	player.SortByScore(currentPractice.TeamOnePlayers)
	player.SortByScore(currentPractice.TeamTwoPlayers)
	component := components.PracticePage(currentPractice)
	component.Render(r.Context(), w)
}

func (cfg *webConfig) handleMovePlayer(w http.ResponseWriter, r *http.Request) {
	practiceId, err := strconv.Atoi(r.PathValue("practice_id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("unable to parse practice id from path: %v", err)
		return
	}

	playerId, err := strconv.Atoi(r.PathValue("player_id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("unable to parse practice id from path: %v", err)
		return
	}
	dbPracticePlayer, err := cfg.queries.GetPracticePlayer(r.Context(), db.GetPracticePlayerParams{
		PracticeID: int64(practiceId),
		PlayerID:   int64(playerId),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to get player with ID [%d] from practice [%d]: %v", playerId, practiceId, err)
		return
	}

	switch dbPracticePlayer.TeamNumber {
	case 1:
		err = cfg.queries.SetPlayerTeam(r.Context(), db.SetPlayerTeamParams{
			PracticeID: dbPracticePlayer.PracticeID,
			PlayerID:   dbPracticePlayer.PlayerID,
			TeamNumber: 2,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to move player with ID [%d] from practice [%d] from team 1 to 2: %v", playerId, practiceId, err)
			return
		}
	case 2:
		err = cfg.queries.SetPlayerTeam(r.Context(), db.SetPlayerTeamParams{
			PracticeID: dbPracticePlayer.PracticeID,
			PlayerID:   dbPracticePlayer.PlayerID,
			TeamNumber: 1,
		})
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("/practices/%d", practiceId))
}
