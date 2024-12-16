package team

import "fmt"

func (team *Team) Details() {
	fmt.Printf("%s players\n", team.name)
	for i, player := range team.players {
		fmt.Printf("%d. %s\n", i+1, player.Name)
	}
	fmt.Printf("\n%s has %d players with total score of %.1f\n\n", team.name, len(team.players), team.score())
}

func (team *Team) score() float64 {
	totalScore := 0.0
	for _, player := range team.players {
		totalScore += player.GetScore()
	}
	return totalScore
}
