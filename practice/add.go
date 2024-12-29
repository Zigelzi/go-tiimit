package practice

import (
	"slices"

	"github.com/Zigelzi/go-tiimit/player"
)

func (practice *Practice) Add(attendingPlayer player.Player) bool {
	index := slices.IndexFunc(practice.Players, func(searchedPlayer player.Player) bool {
		return attendingPlayer.MyClubId == searchedPlayer.MyClubId
	})

	if index == -1 {
		practice.Players = append(practice.Players, attendingPlayer)
		return true
	}
	return false
}
