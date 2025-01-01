package team

import (
	"fmt"
	"strings"
)

func (team *Team) Details() {

	teamDetails := fmt.Sprintf("%s has %d players with total score of %.1f", team.name, len(team.players), team.score())

	fmt.Printf("%s players\n", team.name)
	for i, player := range team.players {
		fmt.Printf("%d. %s\n", i+1, player.Name)
	}
	fmt.Println(teamDetails)
	fmt.Println(strings.Repeat("=", len(teamDetails)))
}

func (team *Team) score() float64 {
	totalScore := 0.0
	for _, player := range team.players {
		totalScore += player.Score()
	}
	return totalScore
}
