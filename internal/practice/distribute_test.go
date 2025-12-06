package practice

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Zigelzi/go-tiimit/internal/player"
)

func TestDistributingPlayers(t *testing.T) {
	type testCase struct {
		name                   string
		fieldPlayers           []player.Player
		goalies                []player.Player
		expectedErr            error
		expectedTeamOnePlayers []player.Player
		expectedTeamTwoPlayers []player.Player
	}

	fieldPlayers := []player.Player{
		{MyClubId: 1004, Name: "Lauri Laavu", IsGoalie: false},
		{MyClubId: 1005, Name: "Marja Myyry", IsGoalie: false},
		{MyClubId: 1006, Name: "Tarja Tyyry", IsGoalie: false},
		{MyClubId: 1007, Name: "Kalevi Kuurula", IsGoalie: false},
	}
	goalies := []player.Player{
		{MyClubId: 1000, Name: "Matti Meikäläinen", IsGoalie: true},
		{MyClubId: 1001, Name: "Teppo Teikäläinen", IsGoalie: true},
		{MyClubId: 1002, Name: "Kaija Karppi", IsGoalie: true},
		{MyClubId: 1003, Name: "Saija Siirappi", IsGoalie: true},
	}

	testCases := []testCase{
		{
			name:         "creates even teams with even number of field players and goalies",
			fieldPlayers: fieldPlayers,
			goalies:      goalies,
			expectedErr:  nil,
			expectedTeamOnePlayers: []player.Player{
				{MyClubId: 1004, Name: "Lauri Laavu", IsGoalie: false},
				{MyClubId: 1006, Name: "Tarja Tyyry", IsGoalie: false},
				{MyClubId: 1000, Name: "Matti Meikäläinen", IsGoalie: true},
				{MyClubId: 1002, Name: "Kaija Karppi", IsGoalie: true},
			},
			expectedTeamTwoPlayers: []player.Player{
				{MyClubId: 1005, Name: "Marja Myyry", IsGoalie: false},
				{MyClubId: 1007, Name: "Kalevi Kuurula", IsGoalie: false},
				{MyClubId: 1001, Name: "Teppo Teikäläinen", IsGoalie: true},
				{MyClubId: 1003, Name: "Saija Siirappi", IsGoalie: true},
			},
		},
		{
			name:         "creates even teams with even number of field players and odd number of goalies",
			fieldPlayers: fieldPlayers,
			goalies:      goalies[:3],
			expectedErr:  nil,
			expectedTeamOnePlayers: []player.Player{
				{MyClubId: 1004, Name: "Lauri Laavu", IsGoalie: false},
				{MyClubId: 1006, Name: "Tarja Tyyry", IsGoalie: false},
				{MyClubId: 1000, Name: "Matti Meikäläinen", IsGoalie: true},
				{MyClubId: 1002, Name: "Kaija Karppi", IsGoalie: true},
			},
			expectedTeamTwoPlayers: []player.Player{
				{MyClubId: 1005, Name: "Marja Myyry", IsGoalie: false},
				{MyClubId: 1007, Name: "Kalevi Kuurula", IsGoalie: false},
				{MyClubId: 1001, Name: "Teppo Teikäläinen", IsGoalie: true},
			},
		},
		{
			name:         "creates even teams with even number of field players and no goalies",
			fieldPlayers: fieldPlayers,
			goalies:      []player.Player{},
			expectedErr:  nil,
			expectedTeamOnePlayers: []player.Player{
				{MyClubId: 1004, Name: "Lauri Laavu", IsGoalie: false},
				{MyClubId: 1006, Name: "Tarja Tyyry", IsGoalie: false},
			},
			expectedTeamTwoPlayers: []player.Player{
				{MyClubId: 1005, Name: "Marja Myyry", IsGoalie: false},
				{MyClubId: 1007, Name: "Kalevi Kuurula", IsGoalie: false},
			},
		},
		{
			name:         "creates even teams with no field players and even number of goalies",
			fieldPlayers: []player.Player{},
			goalies:      goalies,
			expectedErr:  nil,
			expectedTeamOnePlayers: []player.Player{
				{MyClubId: 1000, Name: "Matti Meikäläinen", IsGoalie: true},
				{MyClubId: 1002, Name: "Kaija Karppi", IsGoalie: true},
			},
			expectedTeamTwoPlayers: []player.Player{
				{MyClubId: 1001, Name: "Teppo Teikäläinen", IsGoalie: true},
				{MyClubId: 1003, Name: "Saija Siirappi", IsGoalie: true},
			},
		},
		{
			name:         "creates even teams with odd number of field players and goalies",
			fieldPlayers: fieldPlayers[:3],
			goalies:      goalies[:3],
			expectedErr:  nil,
			expectedTeamOnePlayers: []player.Player{
				{MyClubId: 1004, Name: "Lauri Laavu", IsGoalie: false},
				{MyClubId: 1006, Name: "Tarja Tyyry", IsGoalie: false},
				{MyClubId: 1001, Name: "Teppo Teikäläinen", IsGoalie: true},
			},
			expectedTeamTwoPlayers: []player.Player{
				{MyClubId: 1005, Name: "Marja Myyry", IsGoalie: false},
				{MyClubId: 1000, Name: "Matti Meikäläinen", IsGoalie: true},
				{MyClubId: 1002, Name: "Kaija Karppi", IsGoalie: true},
			},
		},
		{
			name:                   "errors when there's no players to distribute",
			fieldPlayers:           []player.Player{},
			goalies:                []player.Player{},
			expectedErr:            ErrNoPlayers,
			expectedTeamOnePlayers: []player.Player{},
			expectedTeamTwoPlayers: []player.Player{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			team1, team2, err := Distribute(testCase.fieldPlayers, testCase.goalies)

			if testCase.expectedErr == nil {
				if err != nil {
					t.Errorf("unexpected error: got [%v] want [nil]", err)
					return
				}
			} else {
				if err == nil {
					t.Errorf("expected error: got [nil] want [%v]", testCase.expectedErr)
				}

				if errors.Is(err, testCase.expectedErr) == false {
					t.Errorf("errors don't match: got [%v] want [%v]", err, testCase.expectedErr)
					return
				}
			}
			if len(team1) != len(testCase.expectedTeamOnePlayers) {
				t.Errorf("number of players in team 1 don't match: got [%d] want [%d]", len(team1), len(testCase.expectedTeamOnePlayers))
			}
			if len(team2) != len(testCase.expectedTeamTwoPlayers) {
				t.Errorf("number of players in team 2 don't match: got [%d] want [%d]", len(team2), len(testCase.expectedTeamTwoPlayers))
			}

			if reflect.DeepEqual(team1, testCase.expectedTeamOnePlayers) == false {
				t.Errorf("team 1 players don't match: got [%v] want [%v]", team1, testCase.expectedTeamOnePlayers)
			}

			if reflect.DeepEqual(team2, testCase.expectedTeamTwoPlayers) == false {
				t.Errorf("team 2 players don't match: got [%v] want [%v]", team2, testCase.expectedTeamTwoPlayers)
			}

		})

	}

}
