package file

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ImportAttendancePlayerRows(path string) (attendancePlayerRows []AttendancePlayerRow, err error) {
	// List of players in MyClub start on row 5 (index 4). Rows before that are other details or empty.
	const startIndex = 4

	var columnType = map[string]int{
		"myClubId":   0,
		"name":       1,
		"attendance": 3,
	}

	file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file to import attendees from a file %s: %w", path, err)
	}
	defer closeFile(file)

	rows, err := file.GetRows("Tapahtuma")
	if err != nil {
		return nil, fmt.Errorf("unable to read rows to import attendees from a file: %w", err)
	}
	attendancePlayerRows = []AttendancePlayerRow{}
	for _, row := range rows[startIndex:] {
		attendancePlayerRow, err := newAttendancePlayerRow(row[columnType["myClubId"]], row[columnType["name"]], row[columnType["attendance"]])
		if err != nil {
			return nil, fmt.Errorf("unable to create new player row: %w", err)
		}
		attendancePlayerRows = append(attendancePlayerRows, attendancePlayerRow)
	}
	return attendancePlayerRows, nil
}

func closeFile(openFile *excelize.File) {
	if err := openFile.Close(); err != nil {
		fmt.Println(err)
		return
	}
}
