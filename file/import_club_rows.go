package file

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ImportClubPlayerRows(path string) (clubPlayerRows []ClubPlayerRow, err error) {
	// List of players in MyClub start on row 5 (index 4). Rows before that are other details or empty.
	const startIndex = 4

	var columnType = map[string]int{
		"myClubId":      0,
		"name":          1,
		"run power":     3,
		"ball handling": 4,
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

	clubPlayerRows = []ClubPlayerRow{}
	for _, row := range rows[startIndex:] {
		clubPlayerRow, err := newClubPlayerRow(row[columnType["myClubId"]], row[columnType["name"]], row[columnType["run power"]], row[columnType["ball handling"]])
		if err != nil {
			return nil, fmt.Errorf("unable to create new player row: %w", err)
		}
		clubPlayerRows = append(clubPlayerRows, clubPlayerRow)
	}
	return clubPlayerRows, nil
}
