package player

import (
	"fmt"
	"strconv"

	"github.com/Zigelzi/go-tiimit/db"
	"github.com/Zigelzi/go-tiimit/file"
	"github.com/xuri/excelize/v2"
)

const playerDirectory = "player-files/"

func ImportToClub() error {
	fileName, err := file.Select(playerDirectory)

	if err != nil {
		return err
	}
	file, err := excelize.OpenFile(playerDirectory + fileName)
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
	var addedPlayers []Player

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

		isExisting, err := isExistingPlayer(player.MyClubId)
		if err != nil {
			fmt.Printf("failed to check existing player on row %d: %s\n", i, err)
			return err
		}

		if !isExisting {
			err := Insert(player)
			if err != nil {
				fmt.Printf("failed to insert player on row %d: %s\n", i, err)
				continue
			}
			addedPlayers = append(addedPlayers, player)
		}
	}

	fmt.Printf("Loaded %d players from file %s\n", len(addedPlayers), fileName)
	return nil
}

func closeFile(openFile *excelize.File) {
	if err := openFile.Close(); err != nil {
		fmt.Println(err)
		return
	}
}

func isExistingPlayer(myClubId int64) (isExisting bool, err error) {
	query := "SELECT EXISTS (SELECT 1 FROM players WHERE myclub_id=?)"
	err = db.DB.QueryRow(query, myClubId).Scan(&isExisting)
	if err != nil {
		return isExisting, err
	}
	return isExisting, nil

}
