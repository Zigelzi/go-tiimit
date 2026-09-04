package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Zigelzi/go-tiimit/cmd/web/components"
	"github.com/Zigelzi/go-tiimit/cmd/web/view"
	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/Zigelzi/go-tiimit/internal/player"
)

func (cfg *webConfig) handleViewPlayersPage(w http.ResponseWriter, r *http.Request) {
	showInactive := r.URL.Query().Get("inactive") == "true"
	dbPlayers := []db.Player{}
	if showInactive {
		inactivePlayers, err := cfg.queries.GetInactivePlayers(r.Context())
		if err != nil && errors.Is(err, sql.ErrNoRows) == false {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to get players: %v", err)
			return
		}
		dbPlayers = inactivePlayers
	} else {
		activePlayers, err := cfg.queries.GetActivePlayers(r.Context())
		if err != nil && errors.Is(err, sql.ErrNoRows) == false {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to get players: %v", err)
			return
		}
		dbPlayers = activePlayers
	}

	players := []player.Player{}
	for _, dbPlayer := range dbPlayers {
		players = append(players, player.FromDB(dbPlayer))
	}

	querySortKey := r.URL.Query().Get("sort")
	player.SortPlayersByDescending(players, querySortKey)

	viewPlayers := []view.Player{}
	for _, player := range players {
		viewPlayers = append(viewPlayers, view.FromPlayer(player))
	}

	playersPage := components.PlayersPage(viewPlayers, querySortKey, showInactive)
	playersPage.Render(r.Context(), w)
}

func (cfg *webConfig) handleAddPlayerForm(w http.ResponseWriter, r *http.Request) {
	addPlayerForm := components.AddPlayerForm(view.AddPlayerForm{})
	addPlayerForm.Render(r.Context(), w)
}

func (cfg *webConfig) handleAddPlayerToClub(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("failed to parse form: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	form := view.AddPlayerForm{
		Name:        r.FormValue("name"),
		MyClubId:    r.FormValue("myclub-id"),
		FieldErrors: make(map[string]string),
	}
	if strings.TrimSpace(form.Name) == "" {
		log.Printf("user tried to add player with empty name\n")
		form.FieldErrors["name"] = "Name can't be empty."
	}

	if strings.TrimSpace(form.MyClubId) == "" {
		log.Printf("user tried to add player without MyClub ID\n")
		form.FieldErrors["myclub-id"] = "MyClub ID can't be empty."
	}

	myClubId, err := strconv.ParseInt(form.MyClubId, 10, 64)
	if err != nil {
		log.Printf("MyClubId isn't a number: %\n", err)
		form.FieldErrors["myclub-id"] = "MyClub ID must be a number, like 1337."
	}

	if len(form.FieldErrors) > 0 {
		component := components.AddPlayerForm(form)
		component.Render(r.Context(), w)
		return
	}

	// Start with default run power and ball handling.
	const defaultSkillRating = 5
	player := player.New(myClubId, form.Name, defaultSkillRating, defaultSkillRating, false)
	isExistingPlayer, err := cfg.queries.IsExistingPlayer(r.Context(), player.MyClubId)
	if err != nil {
		log.Printf("failed to check if player exists: %v\n", err)
		form.GeneralError = "Couldn't add new player. Please try again later."
	}
	if isExistingPlayer != 0 {
		fmt.Printf("player with MyClubId [%d] already exists\n", myClubId)
		form.FieldErrors["myclub-id"] = "Player with this MyClub ID already exists. Check the ID or find them in the player list."
	}

	if len(form.FieldErrors) > 0 || form.GeneralError != "" {
		component := components.AddPlayerForm(form)
		component.Render(r.Context(), w)
		return
	}

	err = cfg.queries.AddPlayer(r.Context(), db.AddPlayerParams{
		Name:         player.Name,
		MyclubID:     player.MyClubId,
		RunPower:     player.RunPower(),
		BallHandling: player.BallHandling(),
	})
	if err != nil {
		log.Printf("failed to add new player: %v\n", err)
		form.GeneralError = "Couldn't add new player. Please try again later."
	}

	if form.GeneralError != "" {
		component := components.AddPlayerForm(form)
		component.Render(r.Context(), w)
		return
	}
	w.Header().Add("HX-Redirect", "/players")
}

func (cfg *webConfig) handleEditPlayerRow(w http.ResponseWriter, r *http.Request) {
	playerId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Printf("failed to convert the id [%s] to int64: %v", r.PathValue("id"), err)
		return
	}

	dbPlayer, err := cfg.queries.GetPlayerById(r.Context(), playerId)
	if err != nil && errors.Is(err, sql.ErrNoRows) == false {
		log.Printf("failed to get player with id %d: %v", playerId, err)
		return
	}
	viewPlayer := view.FromPlayer(player.FromDB(dbPlayer))
	editRow := view.EditPlayerRow{
		Player:      viewPlayer,
		FieldErrors: make(map[string]string),
	}
	component := components.PlayerRowEdit(editRow)
	component.Render(r.Context(), w)
}

func (cfg *webConfig) handleSavePlayerRow(w http.ResponseWriter, r *http.Request) {
	form := view.EditPlayerRow{
		FieldErrors: make(map[string]string),
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("failed to parse form: %v", err)
		form.GeneralError = "Something went wrong when parsing the form. Try again later."
	}

	playerId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Printf("failed to convert the id [%s] to int64: %v", r.PathValue("id"), err)
		form.GeneralError = "Player ID needs to be number."
	}

	form.ID = playerId

	// TODO: Fetch the player details to display name and myclub ID correctly.
	dbPlayer, err := cfg.queries.GetPlayerById(r.Context(), form.ID)
	if err != nil && errors.Is(err, sql.ErrNoRows) == false {
		log.Printf("failed to query player with ID [%d]: %v", form.ID, err)
		form.GeneralError = "Something went wrong when updating the player. Try again later."
	} else if errors.Is(err, sql.ErrNoRows) {
		log.Printf("user tried to query player with ID [%d] when updating their details", form.ID)
		form.GeneralError = fmt.Sprintf("Player with ID %d doesn't exist", form.ID)
	}

	form.Player = view.FromPlayer(player.FromDB(dbPlayer))

	runPower, err := strconv.ParseFloat(r.PostForm.Get("run-power"), 64)
	if err != nil {
		log.Printf("failed to convert run power [%s] to float: %v", r.PostForm.Get("run-power"), err)
		form.FieldErrors["run-power"] = "Run power needs to be number between 0-10."
	}
	ballHandling, err := strconv.ParseFloat(r.PostForm.Get("ball-handling"), 64)
	if err != nil {
		log.Printf("failed to convert ball handling [%s] to float: %v", r.PostForm.Get("ball-handling"), err)
		form.FieldErrors["ball-handling"] = "Ball handling needs to be number between 0-10."
	}

	form.RunPower = runPower
	form.BallHandling = ballHandling
	form.IsGoalie = r.PostForm.Has("is-goalie")

	err = player.ValidateScore(form.RunPower)
	if errors.Is(err, player.ErrInvalidScore) {
		log.Printf("run power needs to be between 0-10")
		form.FieldErrors["run-power"] = "Run power needs to be number between 0-10."
	}

	err = player.ValidateScore(form.BallHandling)
	if errors.Is(err, player.ErrInvalidScore) {
		log.Printf("ball handling needs to be between 0-10")
		form.FieldErrors["ball-handling"] = "Ball handling needs to be number between 0-10."
	}

	if len(form.FieldErrors) > 0 || form.GeneralError != "" {
		component := components.PlayerRowEdit(form)
		component.Render(r.Context(), w)
		return
	}

	updatedPlayer, err := cfg.queries.UpdatePlayerAttributes(r.Context(), db.UpdatePlayerAttributesParams{
		RunPower:     form.RunPower,
		BallHandling: form.BallHandling,
		IsGoalie:     form.IsGoalie,
		ID:           form.ID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("player with ID [%d] doesn't exist", playerId)
			form.GeneralError = fmt.Sprintf("Player with ID %d doesn't exist.", playerId)
		} else {
			log.Printf("failed to update player ID [%d] attributes: %v", playerId, err)
			form.GeneralError = "Something went wrong when updating a player. Try again later."
		}
	}

	viewPlayer := view.FromPlayer(player.FromDB(updatedPlayer))
	component := components.PlayerRow(viewPlayer)
	component.Render(r.Context(), w)
}

func (cfg *webConfig) handleViewPlayerRow(w http.ResponseWriter, r *http.Request) {
	playerId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Printf("failed to convert the id [%s] to int64: %v", r.PathValue("id"), err)
		return
	}

	dbPlayer, err := cfg.queries.GetPlayerById(r.Context(), playerId)
	if err != nil && errors.Is(err, sql.ErrNoRows) == false {
		log.Printf("failed to get player with id %d: %v", playerId, err)
		return
	}
	viewPlayer := view.FromPlayer(player.FromDB(dbPlayer))

	component := components.PlayerRow(viewPlayer)
	component.Render(r.Context(), w)
}

func (cfg *webConfig) handleSetPlayerInactive(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("failed to parse form: %v", err)
		renderError(w, r, "Couldn't set player to inactive. Try again later")
		return
	}

	playerId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Printf("failed to convert the id [%s] to int64: %v", r.PathValue("id"), err)
		renderError(w, r, "Couldn't set player to inactive. Try again later")
		return
	}

	dbPlayers := []db.Player{}
	isInactivated := r.Form.Get("inactive") == "true"
	if isInactivated {
		err = cfg.queries.SetPlayerInactive(r.Context(), playerId)
		if err != nil {
			log.Printf("failed to set player [%d] to inactive: %v", playerId, err)
			renderError(w, r, "Couldn't set player to inactive. Try again later")
			return
		}
		activePlayers, err := cfg.queries.GetActivePlayers(r.Context())

		if err != nil && errors.Is(err, sql.ErrNoRows) == false {
			log.Printf("failed to get players: %v", err)
			renderError(w, r, "Couldn't set player to inactive. Try again later")
			return
		}
		dbPlayers = activePlayers
	} else {
		err = cfg.queries.SetPlayerActive(r.Context(), playerId)
		if err != nil {
			log.Printf("failed to set player [%d] to active: %v", playerId, err)
			renderError(w, r, "Couldn't set player to inactive. Try again later")
			return
		}
		inactivePlayers, err := cfg.queries.GetInactivePlayers(r.Context())
		if err != nil && errors.Is(err, sql.ErrNoRows) == false {
			log.Printf("failed to get players: %v", err)
			renderError(w, r, "Couldn't set player to inactive. Try again later")
			return
		}
		dbPlayers = inactivePlayers
	}

	players := []player.Player{}
	for _, dbPlayer := range dbPlayers {
		players = append(players, player.FromDB(dbPlayer))
	}

	viewPlayers := []view.Player{}
	for _, player := range players {
		viewPlayers = append(viewPlayers, view.FromPlayer(player))
	}

	renderOK(w, r, components.PlayerTableContent(viewPlayers))
}
