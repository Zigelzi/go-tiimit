package practice

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/Zigelzi/go-tiimit/internal/player"
)

func TestDistributingPlayers(t *testing.T) {
	type testCase struct {
		scenario     string
		fieldPlayers []player.Player
		goalies      []player.Player
	}

	fieldPlayers := []player.Player{
		player.New(1004, "Lauri Laavu", 10, 8, false),
		player.New(1005, "Marja Myyry", 3, 3, false),
		player.New(1006, "Tarja Tyyry", 7, 10, false),
		player.New(1007, "Tiina Myrtty", 4, 6, false),
		player.New(1008, "Kalevi Kuurula", 2, 3, false),
		player.New(1009, "Perttu Hyppölä", 6, 4, false),
		player.New(1010, "Marjo Seppälä", 2, 4, false),
		player.New(1011, "Mauno Marttila", 10, 10, false),
	}
	goalies := []player.Player{
		player.New(1000, "Matti Meikäläinen", 10, 10, true),
		player.New(1000, "Niilo Villamo", 2, 2, true),
		player.New(1000, "Eerika Läppä", 2, 2, true),
		player.New(1002, "Kaija Karppi", 6, 5, true),
		player.New(1001, "Teppo Teikäläinen", 7, 7, true),
		player.New(1003, "Saija Siirappi", 4, 4, true),
	}

	testCases := []testCase{
		{
			scenario:     "even field players even goalies",
			fieldPlayers: fieldPlayers,
			goalies:      goalies,
		},
		{
			scenario:     "even field players odd goalies",
			fieldPlayers: fieldPlayers,
			goalies:      goalies[:3],
		},
		{
			scenario:     "only field players",
			fieldPlayers: fieldPlayers,
			goalies:      []player.Player{},
		},
		{
			scenario:     "only goalies",
			fieldPlayers: []player.Player{},
			goalies:      goalies,
		},
		{
			scenario:     "odd field players odd goalies",
			fieldPlayers: fieldPlayers[:3],
			goalies:      goalies[:3],
		},
		{
			scenario:     "single field player",
			fieldPlayers: fieldPlayers[:1],
			goalies:      []player.Player{},
		},
		{
			scenario:     "two field players",
			fieldPlayers: fieldPlayers[:2],
			goalies:      []player.Player{},
		},
		{
			scenario:     "more goalies than field players",
			fieldPlayers: fieldPlayers[:2],
			goalies:      goalies,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.scenario, func(t *testing.T) {
			team1, team2, err := Distribute(testCase.fieldPlayers, testCase.goalies)
			logTeamsOnFailure(t, team1, team2)

			if err != nil {
				t.Fatalf("distributing players failed: %v", err)
			}

			assertTeamsDifferAtMostOnePlayer(t, team1, team2)
			assertTeamsDifferAtMostOneGoalie(t, team1, team2)
			assertEveryPlayerIsDistributedOnce(t, testCase.fieldPlayers, testCase.goalies, team1, team2)
			assertTeamScoresAreBalanced(t, team1, team2)
		})

	}
}

func assertTeamsDifferAtMostOnePlayer(t *testing.T, team1, team2 []player.Player) {
	t.Helper()

	playerDifference := len(team1) - len(team2)
	if playerDifference < 0 {
		playerDifference = -playerDifference
	}

	if playerDifference > 1 {
		t.Errorf("team sizes differ by %d, can be at most 1 (team 1: %d, team 2: %d)", playerDifference, len(team1), len(team2))
	}
}
func assertTeamsDifferAtMostOneGoalie(t *testing.T, team1, team2 []player.Player) {
	t.Helper()
	team1Goalies := 0
	for _, team1Player := range team1 {
		if team1Player.IsGoalie {
			team1Goalies++
		}
	}
	team2Goalies := 0
	for _, team2Player := range team2 {
		if team2Player.IsGoalie {
			team2Goalies++
		}
	}
	goalieDifference := team1Goalies - team2Goalies
	if goalieDifference < 0 {
		goalieDifference = -goalieDifference
	}

	if goalieDifference > 1 {
		t.Errorf("number of goalies differ by %d, can be at most 1 (team 1: %d, team 2: %d)",
			goalieDifference, team1Goalies, team2Goalies)
	}
}
func assertEveryPlayerIsDistributedOnce(t *testing.T, fieldPlayers, goalies, team1, team2 []player.Player) {
	t.Helper()

	playerMap := make(map[int64]int)
	for _, fieldPlayer := range fieldPlayers {
		playerMap[fieldPlayer.MyClubId] += 1
	}
	for _, goalie := range goalies {
		playerMap[goalie.MyClubId] += 1
	}

	for _, team1Player := range team1 {
		playerMap[team1Player.MyClubId] -= 1

	}
	for _, team2Player := range team2 {
		playerMap[team2Player.MyClubId] -= 1
	}

	for myClubId, playerSeen := range playerMap {
		if playerSeen > 0 {
			t.Errorf("player %d was not distributed to any team (%d missing)", myClubId, playerSeen)
		}
		if playerSeen < 0 {
			t.Errorf("player %d has been distributed %d additional times", myClubId, -playerSeen)
		}
	}
}

func assertTeamScoresAreBalanced(t *testing.T, team1, team2 []player.Player) {
	t.Helper()

	highestPlayerScore := getHighestPlayerScore(team1, team2)

	team1Score := 0.0
	for _, team1Player := range team1 {
		team1Score += team1Player.Score()
	}

	team2Score := 0.0
	for _, team2Player := range team2 {
		team2Score += team2Player.Score()
	}

	scoreDifference := math.Abs(team1Score - team2Score)
	if scoreDifference > highestPlayerScore {
		t.Errorf("total score of teams differ by %.1f, can be at most %.1f (team 1: %.1f, team 2: %.1f)",
			scoreDifference, highestPlayerScore, team1Score, team2Score)
	}
}

func getHighestPlayerScore(team1, team2 []player.Player) float64 {
	if len(team1) == 0 && len(team2) == 0 {
		return 0
	}

	highestScore := 0.0

	for _, team1Player := range team1 {
		if team1Player.Score() > highestScore {
			highestScore = team1Player.Score()
		}
	}

	for _, team2Player := range team2 {
		if team2Player.Score() > highestScore {
			highestScore = team2Player.Score()
		}
	}

	return highestScore
}

func logTeamsOnFailure(t *testing.T, team1, team2 []player.Player) {
	t.Helper()
	t.Cleanup(func() {
		if t.Failed() {
			t.Logf("team 1 (%d players)\n", len(team1))
			for _, team1Player := range team1 {
				t.Logf("  %s", formatPlayer(team1Player))
			}
			t.Log("")
			t.Logf("team 2 (%d players)\n", len(team2))
			for _, team2Player := range team2 {
				t.Logf("  %s", formatPlayer(team2Player))
			}

		}
	})
}

func formatPlayer(p player.Player) string {
	return fmt.Sprintf("MyClubId: %d, Name: %s, Score: %.1f, Goalie: %t",
		p.MyClubId, p.Name, p.Score(), p.IsGoalie)
}

func TestDistributingWithoutPlayersReturnsError(t *testing.T) {
	noFieldPlayers := []player.Player{}
	noGoalies := []player.Player{}
	team1, team2, err := Distribute(noFieldPlayers, noGoalies)

	if err == nil {
		t.Fatalf("expected error: got nil")
	}
	if errors.Is(err, ErrNoPlayers) == false {
		t.Errorf("expected error [%s], got [%v]", ErrNoPlayers, err)
	}

	if len(team1) != 0 || len(team2) != 0 {
		t.Errorf("teams shouldn't have any players (team 1: %d, team 2: %d)", len(team1), len(team2))
	}

}
