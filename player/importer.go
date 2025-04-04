package player

import (
	"fmt"

	"github.com/Zigelzi/go-tiimit/file"
)

const playerDirectory = "player-files/"

func ImportToClub() error {
	fileName, err := file.Select(playerDirectory)

	if err != nil {
		return err
	}

	playerRows, err := file.ImportClubPlayerRows(playerDirectory + fileName)
	if err != nil {
		return err
	}

	var addedPlayers []Player
	var importedPlayers []Player

	for i, clubPlayerRow := range playerRows {
		player := New(int64(clubPlayerRow.PlayerRow.MyClubId), clubPlayerRow.PlayerRow.Name, clubPlayerRow.RunPower, clubPlayerRow.BallHandling, false)

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
			addedPlayers = append(addedPlayers, player)
		} else {
			importedPlayers = append(importedPlayers, player)
		}
	}

	fmt.Printf("Imported %d players to club which of %d were added to database from a file %s\n\n", len(importedPlayers), len(addedPlayers), fileName)
	return nil
}
