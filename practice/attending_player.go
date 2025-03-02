package practice

import (
	"fmt"

	"github.com/Zigelzi/go-tiimit/player"
)

type AttendingPlayer struct {
	Player     player.Player
	Attendance AttendanceStatus
}

func NewAttendingPlayer(p player.Player, attendanceStatus string) (AttendingPlayer, error) {
	parsedAttendanceStatus := determineStatus(attendanceStatus)
	if parsedAttendanceStatus == AttendanceInvalid {
		return AttendingPlayer{}, fmt.Errorf("invalid attendance status: %s", attendanceStatus)
	}
	return AttendingPlayer{
			Player:     p,
			Attendance: parsedAttendanceStatus,
		},
		nil
}
