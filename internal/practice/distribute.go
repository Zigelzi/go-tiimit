package practice

import (
	"errors"

	"github.com/Zigelzi/go-tiimit/internal/player"
)

func Distribute(fieldPlayers, goalies []player.Player) (teamOnePlayers, teamTwoPlayers []player.Player, err error) {
	teamOnePlayers = []player.Player{}
	teamTwoPlayers = []player.Player{}

	if len(fieldPlayers) == 0 && len(goalies) == 0 {
		return teamOnePlayers, teamTwoPlayers, errors.New("no players to distribute")
	}
	distributePlayers(fieldPlayers, &teamOnePlayers, &teamTwoPlayers)
	distributePlayers(goalies, &teamOnePlayers, &teamTwoPlayers)

	return teamOnePlayers, teamTwoPlayers, nil
}

func distributePlayers(players []player.Player, team1, team2 *[]player.Player) {
	startFromOdd := true
	if len(*team1) != len(*team2) && isEven(len(players)) == false {
		// Need to switch order of distribution when number of players in teams isn't even and number of distributed players is odd.
		// This is to ensure teams are distributed evenly.
		startFromOdd = false
	}
	for i, distributedPlayer := range players {
		if startFromOdd {
			if isEven(i + 1) {
				*team1 = append(*team1, distributedPlayer)
			} else {
				*team2 = append(*team2, distributedPlayer)
			}

		} else {
			if isEven(i) {
				*team1 = append(*team1, distributedPlayer)
			} else {
				*team2 = append(*team2, distributedPlayer)
			}
		}
	}
}

func isEven(number int) bool {
	return number%2 == 0
}
