package file

import (
	"fmt"
	"strconv"
)

type PlayerRow struct {
	MyClubId   int
	Name       string
	Attendance string
}

type AttendanceStatus int

const (
	AttendanceIn AttendanceStatus = iota
	AttendanceOut
	AttendanceUnknown
)

var attendanceName = map[string]AttendanceStatus{
	"Osallistuu":   AttendanceIn,
	"Ei osallistu": AttendanceOut,
	"Ei vastausta": AttendanceUnknown,
}

func NewPlayerRow(newMyClubId, name, attendanceStatus string) (PlayerRow, error) {
	myClubId, err := strconv.Atoi(newMyClubId)
	if err != nil {
		return PlayerRow{}, fmt.Errorf("unable to convert MyClubId to integer: %w", err)
	}

	if name == "" {
		return PlayerRow{}, fmt.Errorf("player name can't be empty")
	}

	_, exists := attendanceName[attendanceStatus]
	if !exists {
		return PlayerRow{}, fmt.Errorf("unknown attendance status %s: %w", attendanceStatus, err)
	}

	return PlayerRow{
			MyClubId:   myClubId,
			Name:       name,
			Attendance: attendanceStatus,
		},
		nil
}
