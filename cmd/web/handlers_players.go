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

	dbPlayers, err := cfg.queries.GetAllPlayers(r.Context())

	if err != nil && errors.Is(err, sql.ErrNoRows) == false {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to get players: %v", err)
		return
	}

	viewPlayers := []view.Player{}
	for _, dbPlayer := range dbPlayers {
		player := player.FromDB(dbPlayer)
		viewPlayers = append(viewPlayers, view.Player{
			ID:           player.ID,
			MyclubId:     player.MyClubId,
			Name:         player.Name,
			IsGoalie:     player.IsGoalie,
			RunPower:     player.RunPower(),
			BallHandling: player.BallHandling(),
			Score:        player.Score(),
		})
	}
	playersPage := components.PlayersPage(viewPlayers)
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
