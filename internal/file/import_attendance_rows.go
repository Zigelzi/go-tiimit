package file

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ImportAttendancePlayerRows(path string) (attendancePlayerRows []AttendancePlayerRow, err error) {
	file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file to import attendees from a file %s: %w", path, err)
	}
	defer closeFile(file)

	rows, err := file.GetRows("Tapahtuma")
	if err != nil {
		return nil, fmt.Errorf("unable to read rows to import attendees from a file: %w", err)
	}

	// List of players in MyClub start on row 5 (index 4). Rows before that are other details or empty.
	const startIndex = 4
	// TODO: Add validation to ensure correct columns exist and are in expected order.
	attendancePlayerRows, err = parseAttendanceRows(rows[startIndex:])
	return attendancePlayerRows, nil
}

func closeFile(openFile *excelize.File) {
	if err := openFile.Close(); err != nil {
		fmt.Println(err)
		return
	}
}

func parseAttendanceRows(rows [][]string) ([]AttendancePlayerRow, error) {
	// Takes excel row represented as string[row][column] as input.
	// Returns error when there's no rows to parse.
	// Returns error if myclubId, name or attendance is missing.
	// Returns error if there's not exactly 4 columns.
	// Returns equal number of attendees as there is parsed rows.
	// Returns error if attendance status is not known
	// Attendance rows always contain myclubid, name and attendance
	var columnType = map[string]int{
		"myClubId":   0,
		"name":       1,
		"attendance": 3,
	}
	attendancePlayerRows := []AttendancePlayerRow{}
	if len(rows) == 0 {
		return attendancePlayerRows, fmt.Errorf("no players found from the rows")
	}
	for _, row := range rows {
		attendancePlayerRow, err := newAttendancePlayerRow(row[columnType["myClubId"]], row[columnType["name"]], row[columnType["attendance"]])
		if err != nil {
			return attendancePlayerRows, fmt.Errorf("unable to create new player row: %w", err)
		}
		attendancePlayerRows = append(attendancePlayerRows, attendancePlayerRow)
	}
	return attendancePlayerRows, nil
}
