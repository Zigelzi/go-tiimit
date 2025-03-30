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
	startFromOdd := true
	if len(team1.players) != len(team2.players) && !isEven(len(players)) {
		// Need to switch order of distribution when number of players in teams isn't even and number of distributed players is odd.
		// This is to ensure teams are distributed evenly.
		startFromOdd = false
	}
	for i, distributedPlayer := range players {
		if startFromOdd {
			if isEven(i + 1) {
				team1.players = append(team1.players, distributedPlayer)
			} else {
				team2.players = append(team2.players, distributedPlayer)
			}

		} else {
			if isEven(i) {
				team1.players = append(team1.players, distributedPlayer)
			} else {
				team2.players = append(team2.players, distributedPlayer)
			}
		}
	}
}

func isEven(number int) bool {
	return number%2 == 0
}
