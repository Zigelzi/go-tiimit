package player

import (
	"errors"
	"fmt"

	"github.com/Zigelzi/go-tiimit/internal/file"
)

const playerDirectory = "player-files/"

func ImportToClub() error {
	fileName, err := file.Select(playerDirectory)

	if err != nil {
		return err
	}

	playerRows, err := file.ImportClubPlayerRows(playerDirectory + fileName)
	var errColumnCount *file.ErrorIncorrectColumnCount

	if errors.As(err, &errColumnCount) {
		fmt.Println(errColumnCount.Msg)
		fmt.Println(errColumnCount.ErrorList)
	} else {
		return err
	}

	var addedPlayers []Player
	var importedPlayers []Player

	for i, clubPlayerRow := range playerRows {
		player := New(int64(clubPlayerRow.PlayerRow.MyClubId), clubPlayerRow.PlayerRow.Name, clubPlayerRow.RunPower, clubPlayerRow.BallHandling, false)
		isExisting, err := exists(player.MyClubId)

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
		} else {
			// Update run power and ball handling for all existing players
			// This probably should be done only when there's changes.
			err := player.UpdateRunPower(player.runPower)
			if err != nil {
				fmt.Printf("failed to update player run power on row %d: %s\n", i, err)
				continue
			}
		}
		importedPlayers = append(importedPlayers, player)
	}

	fmt.Printf("Imported %d players to club which of %d were added to database from a file %s\n\n", len(importedPlayers), len(addedPlayers), fileName)
	return nil
}
