package team

import (
	"errors"
	"fmt"
	"sort"

	"example.com/go-tiimit/player"
)

type Team struct {
	name    string
	players []player.Player
}

func New(name string) Team {
	return Team{
		name: name,
	}
}

func Distribute(players []player.Player) (team1 Team, team2 Team, err error) {
	team1 = New("Team 1")
	team2 = New("Team 2")

	if len(players) == 0 {
		return Team{}, Team{}, errors.New("no attending players to distribute")
	}

	sort.Sort(player.ByScore(players))

	for i, player := range players {
		if (i+1)&2 == 0 {
			team1.players = append(team1.players, player)
		} else {
			team2.players = append(team2.players, player)
		}
	}

	return team1, team2, nil
}

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
