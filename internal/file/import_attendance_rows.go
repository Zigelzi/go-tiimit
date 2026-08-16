package file

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

func ImportAttendancePlayerRowsFromReader(reader io.Reader) ([]AttendancePlayerRow, error) {
	excelFile, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to open excel file from reader: %w", err)
	}
	return ImportAttendancePlayerRows(excelFile)
}

func ImportAttendancePlayerRows(file *excelize.File) (attendancePlayerRows []AttendancePlayerRow, err error) {

	rows, err := file.GetRows("Tapahtuma")
	if err != nil {
		return nil, fmt.Errorf("unable to read rows to import attendees from a file: %w", err)
	}

	// List of players in MyClub start on row 5 (index 4). Rows before that are other details or empty.
	const startIndex = 4
	// TODO: Add validation to ensure correct columns exist and are in expected order.
	attendancePlayerRows, err = parseAttendanceRows(rows[startIndex:])
	if err != nil {
		return []AttendancePlayerRow{}, fmt.Errorf("unable to parse attendance rows: %w", err)
	}
	return attendancePlayerRows, nil
}

func parseAttendanceRows(rows [][]string) ([]AttendancePlayerRow, error) {
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
