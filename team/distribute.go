package team

import (
	"errors"

	"github.com/Zigelzi/go-tiimit/player"
)

func Distribute(goalies, fieldPlayers []player.Player) (team1, team2 Team, err error) {
	team1 = New("Team 1")
	team2 = New("Team 2")
	if len(goalies) == 0 && len(fieldPlayers) == 0 {
		return Team{}, Team{}, errors.New("no attending players to distribute")
	}

	distributePlayers(goalies, &team1, &team2)
	distributePlayers(fieldPlayers, &team1, &team2)

	return team1, team2, nil
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
