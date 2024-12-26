package practice

import (
	"fmt"
	"slices"

	"github.com/Zigelzi/go-tiimit/player"
)

func (p *Practice) Remove(playerToRemove player.Player) error {
	index := slices.IndexFunc(p.Players, func(searchedPlayer player.Player) bool {
		return playerToRemove.MyClubId == searchedPlayer.MyClubId
	})

	if index == -1 {
		return fmt.Errorf("player %s wasn't found from the attending players", playerToRemove.Name)
	}

	// Remove the player by combining slices before and after the player to remove
	p.Players = append(p.Players[:index], p.Players[index+1:]...)

	return nil
}
