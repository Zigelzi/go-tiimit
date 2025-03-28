package file

import (
	"fmt"
	"strconv"
)

type PlayerRow struct {
	MyClubId int
	Name     string
}

func newPlayerRow(newMyClubId, name string) (PlayerRow, error) {
	myClubId, err := strconv.Atoi(newMyClubId)
	if err != nil {
		return PlayerRow{}, fmt.Errorf("unable to convert MyClubId to integer: %w", err)
	}

	// TODO: Move to service layer from model.
	if name == "" {
		return PlayerRow{}, fmt.Errorf("player name can't be empty")
	}

	return PlayerRow{
			MyClubId: myClubId,
			Name:     name,
		},
		nil
}

type ClubPlayerRow struct {
	PlayerRow    PlayerRow
	BallHandling float64
	RunPower     float64
}

func newClubPlayerRow(myClubId, name, runPower, ballHandling string) (ClubPlayerRow, error) {
	playerRow, err := newPlayerRow(myClubId, name)
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
	Attendance string
}

func newAttendancePlayerRow(myClubId, name, attendanceStatus string) (AttendancePlayerRow, error) {
	playerRow, err := newPlayerRow(myClubId, name)
	if err != nil {
		return AttendancePlayerRow{}, fmt.Errorf("failed to create base player row: %w", err)
	}
	if !isValidStatus(attendanceStatus) {
		return AttendancePlayerRow{}, fmt.Errorf("unknown attendance status %s: %w", attendanceStatus, err)
	}
	return AttendancePlayerRow{
			PlayerRow:  playerRow,
			Attendance: attendanceStatus,
		},
		nil
}

func isValidStatus(attendanceStatus string) bool {
	validStatuses := [3]string{"Osallistuu", "Ei osallistu", "Ei vastausta"}
	for _, status := range validStatuses {
		if status == attendanceStatus {
			return true
		}
	}
	return false
}
