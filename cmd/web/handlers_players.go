package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/Zigelzi/go-tiimit/cmd/web/components"
	"github.com/Zigelzi/go-tiimit/cmd/web/view"
)

func (cfg *webConfig) handleViewPlayersPage(w http.ResponseWriter, r *http.Request) {

	dbPlayers, err := cfg.queries.GetAllPlayers(r.Context())

	if err != nil && errors.Is(err, sql.ErrNoRows) == false {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to get players: %v", err)
		return
	}

	players := []view.Player{}
	for _, dbPlayer := range dbPlayers {
		players = append(players, view.Player{
			ID:       dbPlayer.ID,
			Name:     dbPlayer.Name,
			IsGoalie: dbPlayer.IsGoalie,
		})
	}
	playersPage := components.PlayersPage(players)
	playersPage.Render(r.Context(), w)
}
