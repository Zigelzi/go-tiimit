package file

import "fmt"

type AttendancePlayerRow struct {
	PlayerRow  PlayerRow
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

func newAttendancePlayerRow(myClubId, name, attendanceStatus string) (AttendancePlayerRow, error) {
	playerRow, err := newPlayerRow(myClubId, name)
	if err != nil {
		return AttendancePlayerRow{}, fmt.Errorf("failed to create base player row: %w", err)
	}
	_, exists := attendanceName[attendanceStatus]
	if !exists {
		return AttendancePlayerRow{}, fmt.Errorf("unknown attendance status %s: %w", attendanceStatus, err)
	}
	return AttendancePlayerRow{
			PlayerRow:  playerRow,
			Attendance: attendanceStatus,
		},
		nil
}
