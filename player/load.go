package player

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func Load(fileName string) error {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return err
	}
	defer closeFile(file)

	rows, err := file.GetRows("Tapahtuma")
	if err != nil {
		return err
	}
	// List of players in MyClub start on row 5. Rows before that are other details or empty.
	playerRows := rows[4:]
	var players []Player

	for i, playerRow := range playerRows {
		name := playerRow[1]
		myClubId, err := strconv.Atoi(playerRow[0])
		if err != nil {
			fmt.Printf("Unable to parse MyClub ID on row %d.\n", i)
			return err
		}

		runPower, err := strconv.ParseFloat(playerRow[3], 64)
		if err != nil {
			fmt.Printf("Unable to parse run power on row %d.", i)
			return err
		}

		ballHandling, err := strconv.ParseFloat(playerRow[4], 64)
		if err != nil {
			fmt.Printf("Unable to parse ball handling on row %d.", i)
			return err
		}
		player := New(int64(myClubId), name, runPower, ballHandling)
		players = append(players, player)

		err = Insert(player)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Loaded %d players from file %s\n", len(players), fileName)
	return nil
}

func closeFile(openFile *excelize.File) {
	if err := openFile.Close(); err != nil {
		fmt.Println(err)
		return
	}
}
