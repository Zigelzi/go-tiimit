package view

import "testing"

func TestGeneratingPlayerURLs(t *testing.T) {
	t.Run("returns players with valid move URL", func(t *testing.T) {
		team := Team{
			Number: 1,
			Players: []Player{
				{ID: 0},
				{ID: 1},
				{ID: 2},
			}}
		practiceId := int64(1)
		expectedPlayers := []Player{
			{ID: 0, MoveURL: "/practices/1/players/0"},
			{ID: 1, MoveURL: "/practices/1/players/1"},
			{ID: 2, MoveURL: "/practices/1/players/2"},
		}

		team.GeneratePlayerURLs(practiceId)

		for i, actualPlayer := range team.Players {
			if actualPlayer.MoveURL != expectedPlayers[i].MoveURL {
				t.Errorf("move URLs don't match: got [%v] want [%v]", actualPlayer.MoveURL, expectedPlayers[i].MoveURL)
			}
		}

	})
	t.Run("returns players with valid toggle vest URL", func(t *testing.T) {
		team := Team{
			Number: 1,
			Players: []Player{
				{ID: 0},
				{ID: 1},
				{ID: 2},
			},
		}
		practiceId := int64(1)
		expectedPlayers := []Player{
			{ID: 0, ToggleVestURL: "/practices/1/players/0/vest"},
			{ID: 1, ToggleVestURL: "/practices/1/players/1/vest"},
			{ID: 2, ToggleVestURL: "/practices/1/players/2/vest"},
		}

		team.GeneratePlayerURLs(practiceId)

		for i, actualPlayer := range team.Players {
			if actualPlayer.ToggleVestURL != expectedPlayers[i].ToggleVestURL {
				t.Errorf("toggle vest URLs don't match: got [%v] want [%v]", actualPlayer.ToggleVestURL, expectedPlayers[i].ToggleVestURL)
			}
		}

	})
}
