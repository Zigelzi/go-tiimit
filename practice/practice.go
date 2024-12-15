package practice

import (
	"slices"

	"example.com/go-tiimit/player"
	"example.com/go-tiimit/team"
)

type Practice struct {
	Players []player.Player
	Teams   [2]team.Team
}

func (practice *Practice) Attend(attendingPlayer player.Player) {
	index := slices.IndexFunc(practice.Players, func(searchedPlayer player.Player) bool {
		return attendingPlayer.MyClubId == searchedPlayer.MyClubId
	})

	if index == -1 {
		practice.Players = append(practice.Players, attendingPlayer)
	}
}
