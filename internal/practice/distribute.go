package practice

import (
	"errors"
	"slices"

	"github.com/Zigelzi/go-tiimit/internal/player"
)

var ErrNoPlayers = errors.New("no players to distribute")

func Distribute(fieldPlayers, goalies []player.Player) (teamOnePlayers, teamTwoPlayers []player.Player, err error) {
	teamOnePlayers = []player.Player{}
	teamTwoPlayers = []player.Player{}

	if len(fieldPlayers) == 0 && len(goalies) == 0 {
		return teamOnePlayers, teamTwoPlayers, ErrNoPlayers
	}

	teamOnePlayers, teamTwoPlayers = distributePlayers(fieldPlayers, teamOnePlayers, teamTwoPlayers)

	teamOnePlayers, teamTwoPlayers = distributePlayers(goalies, teamOnePlayers, teamTwoPlayers)

	return teamOnePlayers, teamTwoPlayers, nil
}

func distributePlayers(players []player.Player, initialTeam1, initialTeam2 []player.Player) (distributedTeam1, distributedTeam2 []player.Player) {
	// Avoid sorting the original list to avoid unexpected side effects
	sortedPlayers := slices.Clone(players)
	player.SortByScore(sortedPlayers)

	team1Score := getTeamScore(initialTeam1)
	distributedTeam1 = initialTeam1

	team2Score := getTeamScore(initialTeam2)
	distributedTeam2 = initialTeam2

	for _, distributedPlayer := range sortedPlayers {
		assignToTeam1 := shouldAssignToTeam1(distributedTeam1, distributedTeam2, team1Score, team2Score)
		if assignToTeam1 {
			distributedTeam1 = append(distributedTeam1, distributedPlayer)
			team1Score += distributedPlayer.Score()
		} else {
			distributedTeam2 = append(distributedTeam2, distributedPlayer)
			team2Score += distributedPlayer.Score()
		}

	}
	return distributedTeam1, distributedTeam2
}

func getTeamScore(teamPlayers []player.Player) float64 {
	teamScore := 0.0
	for _, p := range teamPlayers {
		teamScore += p.Score()
	}
	return teamScore
}

func shouldAssignToTeam1(team1, team2 []player.Player, team1Score, team2Score float64) bool {
	// Teams can have at most one player difference.
	playerDifference := len(team1) - len(team2)
	if playerDifference > 0 {
		return false
	}
	if playerDifference < 0 {
		return true
	}

	// Assign to team 1 if even score (arbitrary)
	if team1Score == team2Score {
		return true
	}
	// Optimize based on score
	return team1Score < team2Score
}
