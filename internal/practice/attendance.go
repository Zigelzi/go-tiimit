package practice

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Zigelzi/go-tiimit/internal/db"
)

func (p *Practice) AddPlayer(myClubId int, status string) error {
	parsedStatus := determineStatus(status)
	if parsedStatus == AttendanceInvalid {
		return fmt.Errorf("invalid status parsed for for player with MyClub ID: [%d]", myClubId)
	}
	p.AttendingPlayers[myClubId] = parsedStatus
	return nil
}

func (p *Practice) GetPlayersByStatus(status AttendanceStatus, playerGetter func(int64) (db.Player, error)) (players []db.Player, unknownPlayerIds []int, allErrors error) {
	currentErrors := []error{}
	for myClubId, playerStatus := range p.AttendingPlayers {
		if status != playerStatus {
			continue
		}
		player, err := playerGetter(int64(myClubId))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				unknownPlayerIds = append(unknownPlayerIds, myClubId)
				continue
			}
			currentErrors = append(currentErrors, fmt.Errorf("unable to get player (MyClub ID: %d) by status [%s]: [%w]", myClubId, status, err))
			continue
		}
		players = append(players, player)
	}

	if len(currentErrors) > 0 {
		return players, unknownPlayerIds, errors.Join(currentErrors...)
	}
	return players, unknownPlayerIds, nil
}
