package practice

import (
	"errors"

	"github.com/Zigelzi/go-tiimit/internal/player"
)

var ErrNoPlayers = errors.New("no players to distribute")

func Distribute(fieldPlayers, goalies []player.Player) (teamOnePlayers, teamTwoPlayers []player.Player, err error) {
	teamOnePlayers = []player.Player{}
	teamTwoPlayers = []player.Player{}

	if len(fieldPlayers) == 0 && len(goalies) == 0 {
		return teamOnePlayers, teamTwoPlayers, ErrNoPlayers
	}
	distributePlayers(fieldPlayers, &teamOnePlayers, &teamTwoPlayers)
	distributePlayers(goalies, &teamOnePlayers, &teamTwoPlayers)

	return teamOnePlayers, teamTwoPlayers, nil
}

func distributePlayers(players []player.Player, team1, team2 *[]player.Player) {
	startFromFirstTeam := true
	if len(*team1) != len(*team2) {
		// First player needs to be distributed to second team when the teams aren't even.
		// This is to ensure teams are distributed evenly and one team doesn't have two more players.
		startFromFirstTeam = false
	}
	for i, distributedPlayer := range players {
		if startFromFirstTeam {
			if isEven(i) {
				*team1 = append(*team1, distributedPlayer)
			} else {
				*team2 = append(*team2, distributedPlayer)
			}

		} else {
			if isEven(i + 1) {
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
