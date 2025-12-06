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

type ClubPlayerRow struct {
	PlayerRow    PlayerRow
	BallHandling float64
	RunPower     float64
}

func newClubPlayerRow(myclubID, name, runPower, ballHandling string) (ClubPlayerRow, error) {
	playerRow, err := newPlayerRow(myclubID, name)
	if err != nil {
		return ClubPlayerRow{}, fmt.Errorf("failed to create base player row: %w", err)
	}

	parsedRunPower, err := strconv.ParseFloat(runPower, 64)
	if err != nil {
		return ClubPlayerRow{}, fmt.Errorf("failed to parse run power: %w", err)
	}
	if parsedRunPower < 0.0 {
		return ClubPlayerRow{}, fmt.Errorf("run power too low: needs to be between 0-10")
	}
	if parsedRunPower > 10.0 {
		return ClubPlayerRow{}, fmt.Errorf("run power too high: needs to be between 0-10")
	}

	parsedBallHandling, err := strconv.ParseFloat(ballHandling, 64)
	if err != nil {
		return ClubPlayerRow{}, fmt.Errorf("failed to parse ball handling: %w", err)
	}

	if parsedBallHandling < 0.0 {
		return ClubPlayerRow{}, fmt.Errorf("ball handling too low: needs to be between 0-10")
	}
	if parsedBallHandling > 10.0 {
		return ClubPlayerRow{}, fmt.Errorf("ball handling too high: needs to be between 0-10")
	}

	return ClubPlayerRow{
		PlayerRow:    playerRow,
		RunPower:     parsedRunPower,
		BallHandling: parsedBallHandling,
	}, nil
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
