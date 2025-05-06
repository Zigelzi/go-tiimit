package player

import "fmt"

type Player struct {
	id           int64
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

func (player Player) Details() (details string) {
	if player.IsGoalie {
		goalieSymbol := "[G]"
		return fmt.Sprintf("%s %s", player.Name, goalieSymbol)
	} else {
		return player.Name
	}
}

func (player Player) UpdateRunPower(newRunPower float64) error {
	if newRunPower < 0 {
		newRunPower = 0
	}
	if newRunPower > 10 {
		newRunPower = 10
	}

	err := updateRunPower(player.MyClubId, newRunPower)
	if err != nil {
		return fmt.Errorf("unable to update player ID [%d] run power to [%.2f]: %w", player.MyClubId, newRunPower, err)
	}
	return nil
}
