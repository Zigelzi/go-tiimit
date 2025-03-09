package team

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Zigelzi/go-tiimit/player"
)

func TestDistribute(t *testing.T) {
	var tests = []struct {
		name          string
		goalies       []player.Player
		fieldPlayers  []player.Player
		wantErr       bool
		expectedErr   string
		expectedTeam1 Team
		expectedTeam2 Team
	}{
		{
			name: "Distribute teams with goalies and field players",
			goalies: []player.Player{
				{MyClubId: 1000, Name: "Matti Meikäläinen", IsGoalie: true},
				{MyClubId: 1001, Name: "Teppo Teikäläinen", IsGoalie: true},
				{MyClubId: 1002, Name: "Kaija Karppi", IsGoalie: true},
				{MyClubId: 1003, Name: "Saija Siirappi", IsGoalie: true},
			},
			fieldPlayers: []player.Player{
				{MyClubId: 1004, Name: "Lauri Laavu", IsGoalie: false},
				{MyClubId: 1005, Name: "Marja Myyry", IsGoalie: false},
				{MyClubId: 1006, Name: "Tarja Tyyry", IsGoalie: false},
				{MyClubId: 1007, Name: "Kalevi Kuurula", IsGoalie: false},
			},
			wantErr:     false,
			expectedErr: "",
			expectedTeam1: Team{
				name: "Team 1",
				players: []player.Player{
					{MyClubId: 1001, Name: "Teppo Teikäläinen", IsGoalie: true},
					{MyClubId: 1003, Name: "Saija Siirappi", IsGoalie: true},
					{MyClubId: 1005, Name: "Marja Myyry", IsGoalie: false},
					{MyClubId: 1007, Name: "Kalevi Kuurula", IsGoalie: false},
				},
			},
			expectedTeam2: Team{
				name: "Team 2",
				players: []player.Player{
					{MyClubId: 1000, Name: "Matti Meikäläinen", IsGoalie: true},
					{MyClubId: 1002, Name: "Kaija Karppi", IsGoalie: true},
					{MyClubId: 1004, Name: "Lauri Laavu", IsGoalie: false},
					{MyClubId: 1006, Name: "Tarja Tyyry", IsGoalie: false},
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			team1, team2, err := Distribute(testCase.goalies, testCase.fieldPlayers)
			if testCase.wantErr {
				if err == nil {
					t.Errorf("error is missing: got [nil] want [%s]", testCase.expectedErr)
					return
				}
				if !strings.Contains(err.Error(), testCase.expectedErr) {
					t.Errorf("error contents don't match: got [%s] want [%s]", err.Error(), testCase.expectedErr)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: got [%s] want [nil]", err)
			}

			if len(team1.players) != len(testCase.expectedTeam1.players) {
				t.Errorf("%s number of players don't match: got [%d] want [%d]", team1.name, len(team1.players), len(testCase.expectedTeam1.players))
			}

			if len(team2.players) != len(testCase.expectedTeam2.players) {
				t.Errorf("%s number of players don't match: got [%d] want [%d]", team2.name, len(team2.players), len(testCase.expectedTeam2.players))
			}

			if !reflect.DeepEqual(team1.players, testCase.expectedTeam1.players) {
				t.Errorf("%s players don't match: got [%v] want [%v]", team1.name, team1.players, testCase.expectedTeam1.players)
			}
			if !reflect.DeepEqual(team2.players, testCase.expectedTeam2.players) {
				t.Errorf("%s players don't match: got [%v] want [%v]", team2.name, team2.players, testCase.expectedTeam2.players)
			}
		})
	}
}
