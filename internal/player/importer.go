package player

import (
	"context"
	"errors"
	"fmt"

	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/Zigelzi/go-tiimit/internal/file"
)

const playerDirectory = "player-files/"

func ImportToClub(dbQuery *db.Queries) error {
	fileName, err := file.Select(playerDirectory)

	if err != nil {
		return err
	}

	playerRows, err := file.ImportClubPlayerRows(playerDirectory + fileName)

	if err != nil {
		var errColumnCount *file.ErrorIncorrectColumnCount
		if errors.As(err, &errColumnCount) {
			fmt.Println(errColumnCount.Msg)
			fmt.Println(errColumnCount.ErrorList)
		} else {
			return err
		}

	}

	var addedPlayers []Player
	var importedPlayers []Player

	for i, clubPlayerRow := range playerRows {
		player := New(int64(clubPlayerRow.PlayerRow.MyclubID), clubPlayerRow.PlayerRow.Name, clubPlayerRow.RunPower, clubPlayerRow.BallHandling, false)
		isExisting, err := dbQuery.IsExistingPlayer(context.Background(), player.MyClubId)

		if err != nil {
			fmt.Printf("failed to check existing player on row %d: %s\n", i, err)
			return err
		}

		if isExisting == 0 {
			err := dbQuery.AddPlayer(context.Background(), db.AddPlayerParams{
				Name:         player.Name,
				MyclubID:     player.MyClubId,
				RunPower:     player.runPower,
				BallHandling: player.ballHandling,
			})
			if err != nil {
				fmt.Printf("failed to insert player on row %d: %s\n", i, err)
				continue
			}
			addedPlayers = append(addedPlayers, player)
		} else {
			// Update run power and ball handling for all existing players
			// This probably should be done only when there's changes.
			err := player.UpdateRunPower(dbQuery, player.runPower)
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
