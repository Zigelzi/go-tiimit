package team

import (
	"errors"
	"sort"

	"github.com/Zigelzi/go-tiimit/player"
)

func Distribute(players []player.Player) (team1, team2 Team, err error) {
	team1 = New("Team 1")
	team2 = New("Team 2")
	if len(players) == 0 {
		return Team{}, Team{}, errors.New("no attending players to distribute")
	}

	sort.Sort(player.ByScore(players))
	goalies, fieldPlayers := getGoalies(players)
	distributePlayers(goalies, &team1, &team2)
	distributePlayers(fieldPlayers, &team1, &team2)

	return team1, team2, nil
}

func getGoalies(players []player.Player) (goalies, fieldPlayers []player.Player) {
	// TODO: Move this to players?
	for _, player := range players {
		if player.IsGoalie {
			goalies = append(goalies, player)
		} else {
			fieldPlayers = append(fieldPlayers, player)
		}
	}

	return goalies, fieldPlayers
}

func distributePlayers(players []player.Player, team1, team2 *Team) {
	for i, distributedPlayer := range players {
		if (i+1)&2 == 0 {
			team1.players = append(team1.players, distributedPlayer)
		} else {
			team2.players = append(team2.players, distributedPlayer)
		}
	}
}
