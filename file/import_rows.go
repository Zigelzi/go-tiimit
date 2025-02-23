package file

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

const startIndex = 4

func ImportPlayerRows(path string) (playerRows [][]string, err error) {
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
	rawPlayerRows := rows[startIndex:]

	return rawPlayerRows, nil
}

func closeFile(openFile *excelize.File) {
	if err := openFile.Close(); err != nil {
		fmt.Println(err)
		return
	}
}
