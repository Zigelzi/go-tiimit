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
	goalies, fieldPlayers := getGoalies(players)
	distributeGoalies(goalies, &team1, &team2)
	distributeFieldPlayers(fieldPlayers, &team1, &team2)
	return team1, team2, nil
}

func getGoalies(players []player.Player) (goalies, fieldPlayers []player.Player) {
	for _, player := range players {
		if player.IsGoalie {
			goalies = append(goalies, player)
		} else {
			fieldPlayers = append(fieldPlayers, player)
		}
	}

	return goalies, fieldPlayers
}

func distributeGoalies(goalies []player.Player, team1, team2 *Team) {
	for i, goalie := range goalies {
		if (i+1)&2 == 0 {
			team1.players = append(team1.players, goalie)
		} else {
			team2.players = append(team2.players, goalie)
		}
	}
}
func distributeFieldPlayers(fieldPlayers []player.Player, team1, team2 *Team) {
	for i, fieldPlayer := range fieldPlayers {
		if (i+1)&2 == 0 {
			team1.players = append(team1.players, fieldPlayer)
		} else {
			team2.players = append(team2.players, fieldPlayer)
		}
	}
}
