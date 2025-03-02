package file

import "fmt"

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
