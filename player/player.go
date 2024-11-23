package player

import (
	"fmt"
	"math/rand/v2"
)

type Player struct {
	Id           int64
	Name         string
	runPower     float64
	ballHandling float64
}

func New(id int64, name string) Player {
	return Player{
		Id:           id,
		Name:         name,
		runPower:     float64(rand.Int64N(9) + 1),
		ballHandling: float64(rand.Int64N(9) + 1),
	}
}

func (player Player) GetScore() float64 {
	const runPowerWeight float64 = 1.2
	const ballHandlingWeight float64 = 1
	return player.runPower*runPowerWeight + player.ballHandling*ballHandlingWeight
}

func (player Player) PrintDetails() {
	fmt.Printf("[%d] %s score: %.1f\n", player.Id, player.Name, player.GetScore())
}
