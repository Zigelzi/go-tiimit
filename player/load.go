package player

import (
	"fmt"
	"strconv"

	"github.com/Zigelzi/go-tiimit/file"
)

const playerDirectory = "player-files/"

func ImportToClub() error {
	fileName, err := file.Select(playerDirectory)

	if err != nil {
		return err
	}

	playerRows, err := file.ImportRows(playerDirectory + fileName)
	if err != nil {
		return err
	}

	var importedPlayers []Player
	var loadedPlayers []Player

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
		player := New(int64(myClubId), name, runPower, ballHandling, false)

		isExisting, err := IsExisting(player.MyClubId)
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
			importedPlayers = append(importedPlayers, player)
		} else {
			loadedPlayers = append(loadedPlayers, player)
		}
	}

	fmt.Printf("Loaded %d players which of %d were imported to database from a file %s\n\n", len(loadedPlayers), len(importedPlayers), fileName)
	return nil
}
