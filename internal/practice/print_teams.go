package practice

import (
	"fmt"
	"strings"

	"github.com/Zigelzi/go-tiimit/internal/player"
)

func (p *Practice) PrintTeams() {
	teamDetails(p.TeamOnePlayers, 1)
	teamDetails(p.TeamTwoPlayers, 2)
}

func teamDetails(players []player.Player, teamNumber int) {
	const separatorCount = 20
	fmt.Printf("Team %d players\n", teamNumber)
	for i, player := range players {
		fmt.Printf("%d. %s [%.1f]\n", i+1, player.Details(), player.Score())
	}
	fmt.Println()
	teamTwoDetails := fmt.Sprintf("Team %d has %d players with total score of %.1f", teamNumber, len(players), totalScore(players))
	fmt.Println(teamTwoDetails)
	fmt.Println()
	fmt.Println(strings.Repeat("=", separatorCount))
	fmt.Println()
}

func totalScore(players []player.Player) float64 {
	totalScore := 0.0
	for _, player := range players {
		totalScore += player.Score()
	}
	return totalScore
}
