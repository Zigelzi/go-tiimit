package practice

import (
	"fmt"
	"slices"

	"github.com/Zigelzi/go-tiimit/player"
)

func (p *Practice) Remove(attendingPlayer player.Player) error {
	if len(p.Players) == 0 {
		return fmt.Errorf("no attending players")
	}
	index := slices.IndexFunc(p.Players, func(searchedPlayer player.Player) bool {
		return attendingPlayer.MyClubId == searchedPlayer.MyClubId
	})

	if index == -1 {
		return fmt.Errorf("player %s wasn't found from the attending players", attendingPlayer.Name)
	}
	p.Players = append(p.Players[:index], p.Players[index+1:]...)

	return nil
}
