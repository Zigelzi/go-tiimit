package file

import (
	"fmt"
	"strconv"
	"time"
)

type FileName struct {
	Path string
	Date time.Time
}

type PlayerRow struct {
	MyclubID int
	Name     string
}

func newPlayerRow(newMyclubID, name string) (PlayerRow, error) {
	myclubID, err := strconv.Atoi(newMyclubID)
	if err != nil {
		return PlayerRow{}, fmt.Errorf("unable to convert MyClubId to integer: %w", err)
	}

	// TODO: Move to service layer from model.
	if name == "" {
		return PlayerRow{}, fmt.Errorf("player name can't be empty")
	}

	return PlayerRow{
			MyclubID: myclubID,
			Name:     name,
		},
		nil
}

type AttendancePlayerRow struct {
	PlayerRow  PlayerRow
	Attendance AttendanceStatus
}

func newAttendancePlayerRow(myclubID, name, attendanceStatus string) (AttendancePlayerRow, error) {
	playerRow, err := newPlayerRow(myclubID, name)
	if err != nil {
		return AttendancePlayerRow{}, fmt.Errorf("failed to create base player row: %w", err)
	}

	parsedAttendance := determineStatus(attendanceStatus)
	if parsedAttendance == AttendanceInvalid {
		return AttendancePlayerRow{}, fmt.Errorf("invalid attendance status %s: %w", attendanceStatus, err)
	}
	return AttendancePlayerRow{
			PlayerRow:  playerRow,
			Attendance: parsedAttendance,
		},
		nil
}
