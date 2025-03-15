package player

import "testing"

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
