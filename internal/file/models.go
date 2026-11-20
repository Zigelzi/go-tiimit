package file

import (
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

type AttendancePlayerRow struct {
	PlayerRow  PlayerRow
	Attendance AttendanceStatus
}
