package team

import (
	"errors"
	"sort"

	"github.com/Zigelzi/go-tiimit/player"
)

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
