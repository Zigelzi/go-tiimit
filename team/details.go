package team

import (
	"fmt"
	"strings"
)

func (team *Team) Details() {
	const separatorCount = 20

	teamDetails := fmt.Sprintf("%s has %d players with total score of %.1f", team.name, len(team.players), team.score())

	fmt.Printf("%s players\n", team.name)
	for i, player := range team.players {
		fmt.Printf("%d. %s [%.1f]\n", i+1, player.Details(), player.Score())
	}
	fmt.Println()
	fmt.Println(teamDetails)
	fmt.Println()
	fmt.Println(strings.Repeat("=", separatorCount))
	fmt.Println()
}

func (team *Team) score() float64 {
	totalScore := 0.0
	for _, player := range team.players {
		totalScore += player.Score()
	}
	return totalScore
}
