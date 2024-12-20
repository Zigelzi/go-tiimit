package practice

import (
	"fmt"
	"slices"

	"github.com/Zigelzi/go-tiimit/player"
)

func (practice *Practice) Add(attendingPlayer player.Player) error {
	index := slices.IndexFunc(practice.Players, func(searchedPlayer player.Player) bool {
		return attendingPlayer.MyClubId == searchedPlayer.MyClubId
	})

	if index == -1 {
		practice.Players = append(practice.Players, attendingPlayer)
		return nil
	} else {
		return fmt.Errorf("player %s already exists", attendingPlayer.Name)
	}
}
