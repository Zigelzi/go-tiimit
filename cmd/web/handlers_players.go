package main

import (
	"net/http"

	"github.com/Zigelzi/go-tiimit/cmd/web/components"
	"github.com/Zigelzi/go-tiimit/cmd/web/view"
)

func (cfg *webConfig) handleViewPlayersPage(w http.ResponseWriter, r *http.Request) {
	players := []view.Player{}
	playersPage := components.PlayersPage(players)
	playersPage.Render(r.Context(), w)
}
