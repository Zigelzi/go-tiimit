package practice

import (
	"fmt"
)

func (p *Practice) AddPlayer(myClubId int, status string) error {
	parsedStatus := determineStatus(status)
	if parsedStatus == AttendanceInvalid {
		return fmt.Errorf("invalid status parsed for for player with MyClub ID: [%d]", myClubId)
	}
	p.AttendingPlayers[myClubId] = parsedStatus
	return nil
}
