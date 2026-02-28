package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Zigelzi/go-tiimit/cmd/web/components"
	"github.com/Zigelzi/go-tiimit/internal/auth"
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
		queryTx.AddPlayerToPractice(r.Context(), db.AddPlayerToPracticeParams{
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
		queryTx.AddPlayerToPractice(r.Context(), db.AddPlayerToPracticeParams{
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
	w.Header().Add("HX-Redirect", fmt.Sprintf("/practice/%d", dbPracticeId))
}

func (cfg webConfig) handleViewPractice(w http.ResponseWriter, r *http.Request) {

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

func (cfg *webConfig) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	component := components.LoginPage()
	component.Render(r.Context(), w)
}

func (cfg *webConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	dbUser, err := cfg.queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.Write([]byte("<p class='text-red-500'>Incorrect username or password </p>"))
			return
		} else {
			log.Printf("failed to get user [%s]: %v", username, err)
			w.Write([]byte("<p class='text-red-500'>Logging in failed, try again later.</p>"))
			return
		}
	}

	password := r.FormValue("password")
	err = auth.CheckPassword(password, dbUser.HashedPassword)
	if err != nil {
		w.Write([]byte("<p class='text-red-500'>Incorrect username or password </p>"))
		return
	}

	sessionToken, err := auth.MakeToken()
	if err != nil {
		log.Printf("failed to generate session token: %v", err)
		w.Write([]byte("<p class='text-red-500'>Logging in failed, try again later.</p>"))
		return
	}
	session, err := cfg.queries.StartUserSession(r.Context(), db.StartUserSessionParams{
		SessionID: sessionToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 14),
	})
	if err != nil {
		log.Printf("failed to start new session: %v", err)
		w.Write([]byte("<p class='text-red-500'>Logging in failed, try again later.</p>"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_session_id",
		Value:    session.SessionID,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  session.ExpiresAt,
	})
	w.Header().Add("HX-Redirect", "/")
}

func (cfg *webConfig) handleLogout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("user_session_id")
	if err != nil {
		return
	}

	err = cfg.queries.EndUserSession(r.Context(), sessionCookie.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) == false {
			log.Printf("failed to log out session [%v]: %v", sessionCookie.Value, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	sessionCookie.Expires = time.Now().Add(-time.Hour)
	http.SetCookie(w, sessionCookie)
	w.Header().Add("HX-Redirect", "/")
}
