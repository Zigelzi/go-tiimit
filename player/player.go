package player

import (
	"fmt"
)

type Player struct {
	MyClubId     int64
	Name         string
	runPower     float64
	ballHandling float64
}

type ByScore []Player

func (p ByScore) Len() int           { return len(p) }
func (p ByScore) Less(i, j int) bool { return p[i].GetScore() > p[j].GetScore() }
func (p ByScore) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func New(id int64, name string, runPower float64, ballHandling float64) Player {
	return Player{
		MyClubId:     id,
		Name:         name,
		runPower:     runPower,
		ballHandling: ballHandling,
	}
}

func (player Player) GetScore() float64 {
	const runPowerWeight float64 = 1.2
	const ballHandlingWeight float64 = 1
	return player.runPower*runPowerWeight + player.ballHandling*ballHandlingWeight
}

func (player Player) PrintDetails() {
	fmt.Printf("[%d] %s score: %.1f\n", player.MyClubId, player.Name, player.GetScore())
}
