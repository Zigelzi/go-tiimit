package player

import (
	"reflect"
	"testing"
)

func TestNewPlayer(t *testing.T) {
	t.Run("player myclub id", func(t *testing.T) {
		myClubIdCases := []struct {
			name             string
			myClubId         int64
			expectedMyClubId int64
		}{
			{"sets positive value", 123, 123},
			{"sets negative value", -123, 0},
			{"sets zero value", 0, 0},
		}

		for _, testCase := range myClubIdCases {
			t.Run(testCase.name, func(t *testing.T) {
				name := "Matti"
				runPower := 10.0
				ballHandling := 10.0
				isGoalie := true
				player := New(testCase.myClubId, name, runPower, ballHandling, isGoalie)
				if player.MyClubId != testCase.expectedMyClubId {
					t.Errorf("got %d, want %d", player.MyClubId, testCase.expectedMyClubId)
				}
			})
		}
	})
	t.Run("player name", func(t *testing.T) {
		nameCases := []struct {
			name         string
			playerName   string
			expectedName string
		}{
			{"capitalized name is set", "Markku", "Markku"},
			{"empty name is set", "", ""},
			{"varying case is persisted", "MarKKu", "MarKKu"},
		}

		for _, testCase := range nameCases {
			t.Run(testCase.name, func(t *testing.T) {
				myClubId := int64(123)
				runPower := 10.0
				ballHandling := 10.0
				isGoalie := true
				player := New(myClubId, testCase.playerName, runPower, ballHandling, isGoalie)
				if player.Name != testCase.expectedName {
					t.Errorf("got %s, want %s", player.Name, testCase.expectedName)
				}
			})
		}
	})
	t.Run("run power", func(t *testing.T) {
		runPowerCases := []struct {
			name             string
			runPower         float64
			expectedRunPower float64
		}{
			{"sets minimum run power ", 0, 0},
			{"sets maximum run power ", 10, 10},
			{"sets run power over max to max ", 100, 10},
			{"sets negative run power to zero ", -10, 0},
		}

		for _, testCase := range runPowerCases {
			t.Run(testCase.name, func(t *testing.T) {
				myClubId := int64(123)
				playerName := "Matti"
				ballHandling := 10.0
				isGoalie := true
				player := New(myClubId, playerName, testCase.runPower, ballHandling, isGoalie)
				if player.runPower != testCase.expectedRunPower {
					t.Errorf("got %f, want %f", player.runPower, testCase.expectedRunPower)
				}
			})
		}
	})

	t.Run("ball handling", func(t *testing.T) {
		ballHandlingCases := []struct {
			name                 string
			ballHandling         float64
			expectedBallHandling float64
		}{
			{"sets minimum ball handling ", 0, 0},
			{"sets maximum ball handling ", 10, 10},
			{"sets ball handling over max to max ", 100, 10},
			{"sets negative ball handling to zero ", -10, 0},
		}

		for _, testCase := range ballHandlingCases {
			t.Run(testCase.name, func(t *testing.T) {
				myClubId := int64(123)
				playerName := "Matti"
				runPower := 10.0
				isGoalie := true
				player := New(myClubId, playerName, runPower, testCase.ballHandling, isGoalie)
				if player.ballHandling != testCase.expectedBallHandling {
					t.Errorf("got %f, want %f", player.runPower, testCase.expectedBallHandling)
				}
			})
		}
	})
}

func TestDetails(t *testing.T) {
	var testCases = []struct {
		name     string
		player   Player
		expected string
	}{
		{
			name:     "Player is goalie",
			player:   Player{Name: "Matti Meikäläinen", IsGoalie: true},
			expected: "Matti Meikäläinen [G]",
		},
		{
			name:     "Player is not a goalie",
			player:   Player{Name: "Matti Meikäläinen", IsGoalie: false},
			expected: "Matti Meikäläinen",
		},
		{
			name:     "Player goalie status is not known",
			player:   Player{Name: "Matti Meikäläinen"},
			expected: "Matti Meikäläinen",
		},
		{
			name:     "Player name is empty and is not goalie",
			player:   Player{Name: "", IsGoalie: false},
			expected: "",
		},
		{
			name:     "Player name is empty and is goalie",
			player:   Player{Name: "", IsGoalie: true},
			expected: " [G]",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.player.Details()
			if actual != testCase.expected {
				t.Errorf("Details() = %q: want %q", actual, testCase.expected)
			}
		})
	}
}

func TestGetPreferences(t *testing.T) {
	var goalies = []Player{
		{Name: "Miika", IsGoalie: true},
		{Name: "Pekka", IsGoalie: true},
		{Name: "Jouni", IsGoalie: true},
		{Name: "Pia", IsGoalie: true},
		{Name: "Seppo", IsGoalie: true},
	}

	var fieldPlayers = []Player{
		{Name: "Janne", IsGoalie: false},
		{Name: "Antti", IsGoalie: false},
		{Name: "Anne", IsGoalie: false},
		{Name: "Mauno", IsGoalie: false},
		{Name: "Kirsi", IsGoalie: false},
	}

	var testCases = []struct {
		name                 string
		initialPlayers       []Player
		expectedGoalies      []Player
		expectedFieldPlayers []Player
	}{
		{
			name:                 "No players",
			initialPlayers:       []Player{},
			expectedGoalies:      []Player{},
			expectedFieldPlayers: []Player{},
		},
		{
			name:                 "Only goalies",
			initialPlayers:       goalies,
			expectedGoalies:      goalies,
			expectedFieldPlayers: []Player{},
		},
		{
			name:                 "Only field players",
			initialPlayers:       fieldPlayers,
			expectedGoalies:      []Player{},
			expectedFieldPlayers: fieldPlayers,
		},
		{
			name:                 "5 goalies and 5 field players",
			initialPlayers:       append(append([]Player{}, goalies...), fieldPlayers...),
			expectedGoalies:      goalies,
			expectedFieldPlayers: fieldPlayers,
		},
		{
			name:                 "1 goalie and 5 field players",
			initialPlayers:       append(append([]Player{}, goalies[:1]...), fieldPlayers...),
			expectedGoalies:      goalies[:1],
			expectedFieldPlayers: fieldPlayers,
		},
		{
			name:                 "5 goalies and 1 field players",
			initialPlayers:       append(append([]Player{}, goalies...), fieldPlayers[:1]...),
			expectedGoalies:      goalies,
			expectedFieldPlayers: fieldPlayers[:1],
		},
		{
			name:                 "nil players",
			initialPlayers:       nil,
			expectedGoalies:      []Player{},
			expectedFieldPlayers: []Player{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualGoalies, actualFieldPlayers := GetPreferences(testCase.initialPlayers)

			if actualGoalies == nil {
				t.Fatalf("expected goalies to be non-nil, got nil")
			}
			if actualFieldPlayers == nil {
				t.Fatalf("expected field players to be non-nil, got nil")
			}

			if len(actualGoalies) != len(testCase.expectedGoalies) {
				t.Errorf("number of goalies doesn't match: got %d want %d", len(actualGoalies), len(testCase.expectedGoalies))
			}
			if len(actualFieldPlayers) != len(testCase.expectedFieldPlayers) {
				t.Errorf("number of goalies doesn't match: got %d want %d", len(actualFieldPlayers), len(testCase.expectedFieldPlayers))
			}

			if !reflect.DeepEqual(actualGoalies, testCase.expectedGoalies) {
				t.Errorf("goalies don't match: got %v want %v", actualGoalies, testCase.expectedGoalies)
			}
		})
	}
}

func TestScore(t *testing.T) {
	var testCases = []struct {
		name         string
		actualPlayer Player
		expected     float64
	}{
		{
			name:         "Player with high attributes",
			actualPlayer: New(123, "Matti Meikäläinen", 10, 10, false),
			expected:     22.0,
		},
		{
			name:         "Run power weight is correct",
			actualPlayer: New(123, "Matti Meikäläinen", 10, 0, false),
			expected:     12.0,
		},
		{
			name:         "Ball handling weight is correct",
			actualPlayer: New(123, "Matti Meikäläinen", 0, 10, false),
			expected:     10.0,
		},
		{
			name:         "Player with zero attributes",
			actualPlayer: New(123, "Matti Meikäläinen", 0, 0, false),
			expected:     0.0,
		},
		{
			name:         "Empty player struct",
			actualPlayer: Player{},
			expected:     0,
		},
		{
			name:         "Player with negative ball handling",
			actualPlayer: New(123, "Matti Meikäläinen", -10, 0, false),
			expected:     0,
		},
		{
			name:         "Player with negative run power",
			actualPlayer: New(123, "Matti Meikäläinen", 0, -10, false),
			expected:     0,
		},
		{
			name:         "Player with negative attributes",
			actualPlayer: New(123, "Matti Meikäläinen", -10, -10, false),
			expected:     0,
		},
		{
			name:         "Run power floating point precision",
			actualPlayer: New(123, "Matti Meikäläinen", 1.123123, 0, false),
			expected:     1.35,
		},
		{
			name:         "Ball handling floating point precision",
			actualPlayer: New(123, "Matti Meikäläinen", 0, 1.123123, false),
			expected:     1.12,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.actualPlayer.Score()
			if actual != testCase.expected {
				t.Errorf("Score() = %.8f: want %.8f", actual, testCase.expected)
			}
		})
	}
}
