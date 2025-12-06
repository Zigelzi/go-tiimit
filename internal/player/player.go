package player

import (
	"context"
	"fmt"

	"github.com/Zigelzi/go-tiimit/internal/db"
)

type Player struct {
	ID           int64
	MyClubId     int64
	Name         string
	runPower     float64
	ballHandling float64
	IsGoalie     bool
}

func New(myclub_id int64, name string, runPower float64, ballHandling float64, isGoalie bool) Player {
	if myclub_id < 0 {
		myclub_id = 0
	}
	if runPower < 0 {
		runPower = 0
	}
	if runPower > 10 {
		runPower = 10
	}
	if ballHandling < 0 {
		ballHandling = 0
	}
	if ballHandling > 10 {
		ballHandling = 10
	}
	return Player{
		MyClubId:     myclub_id,
		Name:         name,
		runPower:     runPower,
		ballHandling: ballHandling,
		IsGoalie:     isGoalie,
	}
}

func FromDB(dbPlayer db.Player) Player {
	newPlayer := New(
		dbPlayer.MyclubID,
		dbPlayer.Name,
		dbPlayer.RunPower,
		dbPlayer.BallHandling,
		dbPlayer.IsGoalie)

	newPlayer.ID = dbPlayer.ID
	return newPlayer
}

func (player Player) Details() (details string) {
	if player.IsGoalie {
		goalieSymbol := "[G]"
		return fmt.Sprintf("%s %s", player.Name, goalieSymbol)
	} else {
		return player.Name
	}
}

func (player Player) UpdateRunPower(dbQuery *db.Queries, newRunPower float64) error {
	if newRunPower < 0 {
		newRunPower = 0
	}
	if newRunPower > 10 {
		newRunPower = 10
	}

	previousRunPower := player.runPower
	err := dbQuery.UpdatePlayerRunPower(context.Background(), db.UpdatePlayerRunPowerParams{
		MyclubID: player.MyClubId,
		RunPower: player.runPower,
	})
	if err != nil {
		return fmt.Errorf("unable to update player ID [%d] run power to [%.2f]: %w", player.MyClubId, newRunPower, err)
	}
	fmt.Printf("Updated player [%d] %s run power from %.1f to %.1f\n", player.MyClubId, player.Name, previousRunPower, player.runPower)
	return nil
}
