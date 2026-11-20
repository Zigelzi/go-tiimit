package file

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// The headers of the columns the import needs in a Finnish MyClub attendance
// report. They match the labels the organizer ticks in MyClub's export
// dialog, so error messages can name them as they are.
const (
	ColumnMyClubID   = "Jäsennumero"
	ColumnName       = "Suku- ja etunimi"
	ColumnAttendance = "Osallistuu"
)

// RequiredColumns returns the headers the import needs, in the order
// messages list them.
func RequiredColumns() []string {
	return []string{ColumnMyClubID, ColumnName, ColumnAttendance}
}

// ErrNoPlayers is returned when a report has no attending player that can be
// read, such as when nobody is attending or every attending row was skipped.
var ErrNoPlayers = errors.New("the report has no readable player rows")

// SkipReason says why a player row was skipped.
type SkipReason int

const (
	SkipUnknownAttendance SkipReason = iota
	SkipMissingMemberNumber
	SkipInvalidMemberNumber
)

// String returns the reason as a short phrase for logs.
func (r SkipReason) String() string {
	switch r {
	case SkipMissingMemberNumber:
		return "missing MyClub ID"
	case SkipInvalidMemberNumber:
		return "MyClub ID is not a number"
	default:
		return "unknown attendance answer"
	}
}

// SkippedRow is a player row the import left out because it could belong to
// an attending player but can't be read. Row is the row number a spreadsheet
// shows, starting from 1. MyClubID holds the cell as it was in the report.
type SkippedRow struct {
	Row      int
	MyClubID string
	Name     string
	Reason   SkipReason
}

// MissingColumnsError is returned when a report lacks required columns.
// Columns lists them in the order of RequiredColumns.
type MissingColumnsError struct {
	Columns []string
}

func (e *MissingColumnsError) Error() string {
	return fmt.Sprintf("the report is missing the columns %s", strings.Join(e.Columns, ", "))
}

// ImportAttendancePlayerRowsFromReader reads the attendance rows from an
// Excel file, and returns the rows it skipped alongside them.
func ImportAttendancePlayerRowsFromReader(reader io.Reader) ([]AttendancePlayerRow, []SkippedRow, error) {
	excelFile, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to open excel file from reader: %w", err)
	}
	defer excelFile.Close()
	return ImportAttendancePlayerRows(excelFile)
}

// ImportAttendancePlayerRows reads the attendance rows from the first sheet
// of file, whatever its name, since a MyClub export has only one sheet. It
// returns the rows it skipped alongside them.
func ImportAttendancePlayerRows(file *excelize.File) ([]AttendancePlayerRow, []SkippedRow, error) {
	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return nil, nil, &MissingColumnsError{Columns: RequiredColumns()}
	}
	rows, err := file.GetRows(sheets[0])
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read rows to import attendees from a file: %w", err)
	}
	return parseAttendanceSheet(rows)
}

// columnPositions holds where each required column is in a row.
type columnPositions struct {
	myClubID, name, attendance int
}

// parseAttendanceSheet finds the header row, locates the required columns in
// it by name, and parses every row below it as a player. A row that could
// belong to an attending player but can't be read is skipped and returned
// in the second list; the other rows are read.
func parseAttendanceSheet(rows [][]string) ([]AttendancePlayerRow, []SkippedRow, error) {
	headerIndex := findHeaderRow(rows)
	if headerIndex == -1 {
		return nil, nil, &MissingColumnsError{Columns: RequiredColumns()}
	}
	columns, err := findColumns(rows[headerIndex])
	if err != nil {
		return nil, nil, err
	}
	playerRows := rows[headerIndex+1:]
	if len(playerRows) == 0 {
		return nil, nil, ErrNoPlayers
	}

	attendancePlayerRows := []AttendancePlayerRow{}
	skippedRows := []SkippedRow{}
	for i, row := range playerRows {
		myClubID := cell(row, columns.myClubID)
		name := cell(row, columns.name)
		answer := cell(row, columns.attendance)
		skip := func(reason SkipReason) {
			// Rows are numbered as a spreadsheet shows them, starting from 1.
			skippedRows = append(skippedRows, SkippedRow{Row: headerIndex + i + 2, MyClubID: myClubID, Name: name, Reason: reason})
		}

		// An empty answer means the player isn't attending. Only attending
		// players are distributed, so the row is left out without a report.
		if answer == "" {
			continue
		}
		status := determineStatus(answer)
		if status == AttendanceInvalid {
			skip(SkipUnknownAttendance)
			continue
		}
		id, err := strconv.Atoi(myClubID)
		if err != nil {
			if status == AttendanceIn {
				if myClubID == "" {
					skip(SkipMissingMemberNumber)
				} else {
					skip(SkipInvalidMemberNumber)
				}
			}
			continue
		}
		attendancePlayerRows = append(attendancePlayerRows, AttendancePlayerRow{
			PlayerRow:  PlayerRow{MyclubID: id, Name: name},
			Attendance: status,
		})
	}
	if !slices.ContainsFunc(attendancePlayerRows, func(row AttendancePlayerRow) bool { return row.Attendance == AttendanceIn }) {
		return nil, skippedRows, ErrNoPlayers
	}
	return attendancePlayerRows, skippedRows, nil
}

// findHeaderRow returns the index of the first row in which any cell is one
// of the required headers, or -1 if there is none. Matching any one header,
// not all of them, lets a report with a missing column still have a header
// row, so the missing column can be named. Status values in player rows can
// equal a header ("Osallistuu"), but the header row always comes first.
func findHeaderRow(rows [][]string) int {
	required := RequiredColumns()
	for i, row := range rows {
		for _, value := range row {
			for _, header := range required {
				if strings.TrimSpace(value) == header {
					return i
				}
			}
		}
	}
	return -1
}

// findColumns returns where each required column is in header, or a
// MissingColumnsError naming those it lacks.
func findColumns(header []string) (columnPositions, error) {
	positions := map[string]int{}
	for i, value := range header {
		value = strings.TrimSpace(value)
		if _, seen := positions[value]; !seen {
			positions[value] = i
		}
	}

	var missing []string
	for _, column := range RequiredColumns() {
		if _, found := positions[column]; !found {
			missing = append(missing, column)
		}
	}
	if len(missing) > 0 {
		return columnPositions{}, &MissingColumnsError{Columns: missing}
	}
	return columnPositions{
		myClubID:   positions[ColumnMyClubID],
		name:       positions[ColumnName],
		attendance: positions[ColumnAttendance],
	}, nil
}

// cell returns the trimmed value at index in row, or "" when the row is
// shorter. Excel readers drop empty cells from the end of a row, so a short
// row is normal and means the cells past its end are empty.
func cell(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}
