package practice

import (
	"slices"

	"github.com/Zigelzi/go-tiimit/player"
)

func (p *Practice) Remove(playerToRemove player.Player) bool {
	index := slices.IndexFunc(p.Players, func(searchedPlayer player.Player) bool {
		return playerToRemove.MyClubId == searchedPlayer.MyClubId
	})

	if index == -1 {
		return false
	}

	// Remove the player by combining slices before and after the player to remove
	p.Players = append(p.Players[:index], p.Players[index+1:]...)

	return true
}
