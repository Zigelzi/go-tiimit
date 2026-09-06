package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Zigelzi/go-tiimit/cmd/web/components"
	"github.com/Zigelzi/go-tiimit/cmd/web/view"
	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/Zigelzi/go-tiimit/internal/file"
	"github.com/Zigelzi/go-tiimit/internal/player"
	"github.com/Zigelzi/go-tiimit/internal/practice"
)

func (cfg *webConfig) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	dbPracticeSummaries, err := cfg.queries.GetNewestPractices(r.Context(), 5)
	if err != nil {
		log.Printf("failed to get practices: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dbPlayers, err := cfg.queries.GetActivePlayers(r.Context())
	if err != nil {
		log.Printf("failed to get all players: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	practiceSummaries := []view.PracticeSummary{}
	for _, dbPracticeSummary := range dbPracticeSummaries {
		practiceSummaries = append(practiceSummaries, view.PracticeSummary{
			ID:          dbPracticeSummary.ID,
			Date:        dbPracticeSummary.Date,
			PlayerCount: dbPracticeSummary.PlayerCount,
		})
	}

	components.IndexPage(practiceSummaries, len(dbPlayers)).Render(r.Context(), w)
}

func (cfg *webConfig) handleSetupPracticePage(w http.ResponseWriter, r *http.Request) {
	components.CreatePracticePage().Render(r.Context(), w)
}

func (cfg *webConfig) handleCreatePractice(w http.ResponseWriter, r *http.Request) {
	formFile, header, err := r.FormFile("attendance-list")
	if err != nil {
		log.Printf("Error parsing file from form: %v", err)
		renderError(w, r, "Couldn't set up new practice. Contact administrator.")
		return
	}
	defer formFile.Close()

	attendanceRows, err := file.ImportAttendancePlayerRowsFromReader(formFile)
	if err != nil {
		log.Printf("failed to parse the attendance rows in handler: %v", err)
		renderError(w, r, "Couldn't process file to set up a new practice right now. Check that the file is formatted correctly and try again.")
		return
	}
	log.Printf("parsed %d rows from attendance excel\n", len(attendanceRows))

	confirmedRows, err := file.GetAttendanceRowsByStatus(attendanceRows, file.AttendanceIn)
	if err != nil {
		log.Printf("failed to get the confirmed rows in handler: %v", err)
		renderError(w, r, "Couldn't get the confirmed players from the file to set up a new practice right now. Check that the file is formatted correctly and try again.")
		return
	}

	dbConfirmedPlayers := []db.Player{}
	for _, row := range confirmedRows {
		confirmedDbPlayer, err := cfg.queries.GetPlayerByMyclubID(r.Context(), int64(row.PlayerRow.MyclubID))
		if err != nil && errors.Is(err, sql.ErrNoRows) == false {
			log.Printf("faile to get confirmed player with myclub_id [%d]: %v", row.PlayerRow.MyclubID, err)
			continue
		} else if errors.Is(err, sql.ErrNoRows) {
			log.Printf("tried to add player that doesn't exist with myclub_id [%d]", row.PlayerRow.MyclubID)
			continue
		}
		dbConfirmedPlayers = append(dbConfirmedPlayers, confirmedDbPlayer)
	}

	confirmedPlayers := []player.Player{}
	for _, dbConfirmedPlayer := range dbConfirmedPlayers {
		confirmedPlayers = append(confirmedPlayers, player.FromDB(dbConfirmedPlayer))
	}

	/*
		TODO: Add unknown players to practice and display them when viewing a practice.
		Commented out now as nothing is done for these at the moment.


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

	*/

	goalies, fieldPlayers := player.GetPreferences(confirmedPlayers)
	team1, team2, err := practice.Distribute(fieldPlayers, goalies)

	if err != nil {
		log.Printf("unable to distribute the players in handler. (team 1: %d players, team 2: %d players): %v", len(team1), len(team2), err)
		renderError(w, r, "Couldn't distribute players to the practice right now. Try again later.")
		return
	}

	practiceDate, err := file.ParseDate(header.Filename)
	if err != nil {
		log.Printf("failed to parse date from the file name [%s]: %v", header.Filename, err)
		renderError(w, r, "File name needs to have date in yyyy-mm-dd format. Check the file name and try again.")
		return
	}

	newPractice := practice.Practice{
		TeamOnePlayers: practice.FromPlayer(team1),
		TeamTwoPlayers: practice.FromPlayer(team2),
		// UnknownPlayers: unknownPlayers,
		Date: practiceDate,
	}

	// REPO START - Storing the data.
	tx, err := cfg.db.Begin()
	if err != nil {
		log.Printf("failed to begin a transaction on practice with date [%v]: %v", newPractice.Date, err)
		renderError(w, r, "Couldn't set up new practice right now. Try again later.")
		return
	}
	defer tx.Rollback()

	queryTx := cfg.queries.WithTx(tx)
	dbPracticeId, err := queryTx.CreatePractice(r.Context(), newPractice.Date)
	if err != nil {
		log.Printf("failed to create the practice with date [%v]: %v", newPractice.Date, err)
		renderError(w, r, "Couldn't set up new practice right now. Try again later.")
		return
	}
	for _, teamOnePlayer := range newPractice.TeamOnePlayers {
		err = queryTx.AddPlayerToPractice(r.Context(), db.AddPlayerToPracticeParams{
			PracticeID: dbPracticeId,
			PlayerID:   teamOnePlayer.Player.ID,
			TeamNumber: 1,
		})
		if err != nil {
			log.Printf("failed to add player [%d] %s to practice %d in team 1", teamOnePlayer.Player.ID, teamOnePlayer.Player.Name, dbPracticeId)
			renderError(w, r, "Couldn't set up new practice right now. Try again later.")
			return
		}
	}

	for _, teamTwoPlayer := range newPractice.TeamTwoPlayers {
		err = queryTx.AddPlayerToPractice(r.Context(), db.AddPlayerToPracticeParams{
			PracticeID: dbPracticeId,
			PlayerID:   teamTwoPlayer.Player.ID,
			TeamNumber: 2,
		})
		if err != nil {
			log.Printf("failed to add player [%d] %s to practice %d in team 2", teamTwoPlayer.Player.ID, teamTwoPlayer.Player.Name, dbPracticeId)
			renderError(w, r, "Couldn't set up new practice right now. Try again later.")
			return
		}
	}
	tx.Commit()
	// REPO END
	w.Header().Add("HX-Redirect", fmt.Sprintf("/practices/%d", dbPracticeId))
}

func (cfg *webConfig) handleViewPracticePage(w http.ResponseWriter, r *http.Request) {

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
	practice.SortByScore(currentPractice.TeamOnePlayers)
	practice.SortByScore(currentPractice.TeamTwoPlayers)

	practiceView := view.Practice{
		ID:    currentPractice.ID,
		Date:  currentPractice.Date,
		Teams: make([]view.Team, 2),
	}
	practiceView.Teams[0] = view.FromPractice(currentPractice.TeamOnePlayers, 1)
	practiceView.Teams[1] = view.FromPractice(currentPractice.TeamTwoPlayers, 2)

	practiceView.Teams[0].GeneratePlayerURLs(currentPractice.ID)
	practiceView.Teams[1].GeneratePlayerURLs(currentPractice.ID)

	component := components.PracticePage(practiceView)
	component.Render(r.Context(), w)
}

func (cfg *webConfig) handleMovePlayer(w http.ResponseWriter, r *http.Request) {
	practiceId, err := strconv.Atoi(r.PathValue("practice_id"))
	if err != nil {
		log.Printf("unable to parse practice id [%s] from path: %v", r.PathValue("practice_id"), err)
		renderError(w, r, "Practice id needs to be a number. Contact administrator.")
		return
	}

	playerId, err := strconv.Atoi(r.PathValue("player_id"))
	if err != nil {
		log.Printf("unable to parse player id [%s] from path: %v", r.PathValue("player_id"), err)
		renderError(w, r, "Player id needs to be a number. Contact administrator.")
		return
	}

	dbPracticePlayer, err := cfg.queries.GetPracticePlayer(r.Context(), db.GetPracticePlayerParams{
		PracticeID: int64(practiceId),
		PlayerID:   int64(playerId),
	})
	if err != nil {
		log.Printf("failed to get player with ID [%d] from practice [%d]: %v", playerId, practiceId, err)
		renderError(w, r, "Couldn't move the player to new team right now. Try again later.")
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
			log.Printf("failed to move player with ID [%d] from practice [%d] from team 1 to 2: %v", playerId, practiceId, err)
			renderError(w, r, "Couldn't move the player to new team right now. Try again later.")
			return
		}
	case 2:
		err = cfg.queries.SetPlayerTeam(r.Context(), db.SetPlayerTeamParams{
			PracticeID: dbPracticePlayer.PracticeID,
			PlayerID:   dbPracticePlayer.PlayerID,
			TeamNumber: 1,
		})
		if err != nil {
			log.Printf("failed to move player with ID [%d] from practice [%d] from team 2 to 1: %v", playerId, practiceId, err)
			renderError(w, r, "Couldn't move the player to new team right now. Try again later.")
			return
		}
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("/practices/%d", practiceId))
}

func (cfg *webConfig) handleTogglePlayerVest(w http.ResponseWriter, r *http.Request) {
	practiceId, err := strconv.Atoi(r.PathValue("practice_id"))
	if err != nil {
		log.Printf("unable to parse practice id [%s] from path: %v", r.PathValue("practice_id"), err)
		renderError(w, r, "Couldn't toggle the vest of a player. Contact administrator.")
		return
	}

	playerId, err := strconv.Atoi(r.PathValue("player_id"))
	if err != nil {
		log.Printf("unable to parse player [%s] id from path: %v", r.PathValue("player_id"), err)
		renderError(w, r, "Couldn't toggle the vest of a player. Contact administrator.")
		return
	}

	updatedPlayerTeamId, err := cfg.queries.TogglePracticePlayerVest(r.Context(), db.TogglePracticePlayerVestParams{
		PracticeID: int64(practiceId),
		PlayerID:   int64(playerId),
	})
	if err != nil {
		log.Printf("failed to toggle the player [%d] vest for practice [%d]: %v", playerId, practiceId, err)
		renderError(w, r, "Couldn't toggle the vest of a player right now. Try again later.")
		return
	}

	dbTeamPlayers, err := cfg.queries.GetPracticeTeamPlayers(r.Context(), db.GetPracticeTeamPlayersParams{
		PracticeID: int64(practiceId),
		TeamNumber: updatedPlayerTeamId,
	})

	if err != nil {
		log.Printf("failed toget the players of practice [%d]: %v", practiceId, err)
		renderError(w, r, "Couldn't toggle the vest of a player right now. Try again later.")
		return
	}

	teamPracticePlayers := []practice.PracticePlayer{}
	for _, dbTeamPlayer := range dbTeamPlayers {
		practicePlayer := practice.PracticePlayerFromDB(dbTeamPlayer)
		teamPracticePlayers = append(teamPracticePlayers, practicePlayer)
	}

	practice.SortByScore(teamPracticePlayers)
	teamView := view.FromPractice(teamPracticePlayers, int(updatedPlayerTeamId))
	teamView.GeneratePlayerURLs(int64(practiceId))

	renderOK(w, r, components.TeamHeader(teamView))
}

func (cfg *webConfig) handleViewTeam(w http.ResponseWriter, r *http.Request) {
	practiceId, err := strconv.Atoi(r.PathValue("practice_id"))
	if err != nil {
		log.Printf("unable to parse practice id [%s] from path: %v", r.PathValue("practice_id"), err)
		renderError(w, r, "Couldn't display team. Contact administrator.")
		return
	}

	teamNumber, err := strconv.Atoi(r.PathValue("team_number"))
	if err != nil {
		log.Printf("unable to parse team number [%s] from path: %v", r.PathValue("team_number"), err)
		renderError(w, r, "Couldn't display team. Contact administrator.")
		return
	}

	if teamNumber < 1 || teamNumber > 2 {
		log.Println("team number must be 1 or 2")
		renderError(w, r, "Couldn't display team. Contact administrator.")
		return
	}
	dbTeamPlayers, err := cfg.queries.GetPracticeTeamPlayers(r.Context(), db.GetPracticeTeamPlayersParams{
		PracticeID: int64(practiceId),
		TeamNumber: int64(teamNumber),
	})

	if err != nil {
		log.Printf("failed toget the players of practice [%d]: %v", practiceId, err)
		renderError(w, r, "Couldn't display team right now. Try again later.")
		return
	}

	teamPracticePlayers := []practice.PracticePlayer{}
	for _, dbTeamPlayer := range dbTeamPlayers {
		practicePlayer := practice.PracticePlayerFromDB(dbTeamPlayer)
		teamPracticePlayers = append(teamPracticePlayers, practicePlayer)
	}

	practice.SortByScore(teamPracticePlayers)
	teamView := view.FromPractice(teamPracticePlayers, int(teamNumber))
	teamView.GeneratePlayerURLs(int64(practiceId))

	renderOK(w, r, components.TeamPanel(teamView))
}
