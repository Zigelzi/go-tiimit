package practice

import (
	"errors"
	"fmt"

	"github.com/Zigelzi/go-tiimit/player"
)

func (p *Practice) AddPlayer(myClubId int, status string) error {
	parsedStatus := determineStatus(status)
	if parsedStatus == AttendanceInvalid {
		return fmt.Errorf("invalid status parsed for for player with MyClub ID: [%d]", myClubId)
	}
	p.AttendingPlayers[myClubId] = parsedStatus
	return nil
}

func (p *Practice) GetPlayersByStatus(status AttendanceStatus, playerGetter func(int64) (player.Player, error)) ([]player.Player, error) {
	players := []player.Player{}
	var errs []error

	for myClubId, playerStatus := range p.AttendingPlayers {
		if status != playerStatus {
			continue
		}
		player, err := playerGetter(int64(myClubId))
		if err != nil {
			errs = append(errs, fmt.Errorf("unable to get player (MyClub ID: %d) by status [%s]: [%w]", myClubId, status, err))
		}
		players = append(players, player)
	}

	if len(errs) > 0 {
		return players, errors.Join(errs...)
	}
	return players, nil
}
