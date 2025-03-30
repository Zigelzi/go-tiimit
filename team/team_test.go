package team

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Zigelzi/go-tiimit/player"
)

func TestDistribute(t *testing.T) {
	goalies := []player.Player{
		{MyClubId: 1000, Name: "Matti Meikäläinen", IsGoalie: true},
		{MyClubId: 1001, Name: "Teppo Teikäläinen", IsGoalie: true},
		{MyClubId: 1002, Name: "Kaija Karppi", IsGoalie: true},
		{MyClubId: 1003, Name: "Saija Siirappi", IsGoalie: true},
	}
	fieldPlayers := []player.Player{
		{MyClubId: 1004, Name: "Lauri Laavu", IsGoalie: false},
		{MyClubId: 1005, Name: "Marja Myyry", IsGoalie: false},
		{MyClubId: 1006, Name: "Tarja Tyyry", IsGoalie: false},
		{MyClubId: 1007, Name: "Kalevi Kuurula", IsGoalie: false},
	}
	t.Run("distributing players returns even teams", func(t *testing.T) {
		var testCases = []struct {
			name                string
			initialGoalies      []player.Player
			initialFieldPlayers []player.Player
			expectedTeam1Length int
			expectedTeam2Length int
		}{
			{
				name:                "even number of goalies and field players",
				initialGoalies:      goalies[:4],
				initialFieldPlayers: fieldPlayers[:4],
				expectedTeam1Length: 4,
				expectedTeam2Length: 4,
			},
			{
				name: "zero goalies and even number of field players",

				initialGoalies:      []player.Player{},
				initialFieldPlayers: fieldPlayers[:4],
				expectedTeam1Length: 2,
				expectedTeam2Length: 2,
			},
			{
				name: "uneven number of goalies and even number of field players",

				initialGoalies:      goalies[:3],
				initialFieldPlayers: fieldPlayers[:4],
				expectedTeam1Length: 3,
				expectedTeam2Length: 4,
			},
			{
				name: "uneven number of goalies and uneven number of field players",

				initialGoalies:      goalies[:3],
				initialFieldPlayers: fieldPlayers[:3],
				expectedTeam1Length: 3,
				expectedTeam2Length: 3,
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				actualTeam1, actualTeam2, actualErr := Distribute(testCase.initialGoalies, testCase.initialFieldPlayers)

				if actualErr != nil {
					t.Errorf("unexpected error: got [%v] want [nil]", actualErr)
				}
				if len(actualTeam1.players) != testCase.expectedTeam1Length {
					t.Errorf("number of players in %s don't match: got [%d] want [%d]", actualTeam1.name, len(actualTeam1.players), testCase.expectedTeam1Length)
				}
				if len(actualTeam2.players) != testCase.expectedTeam2Length {
					t.Errorf("number of players in %s don't match: got [%d] want [%d]", actualTeam2.name, len(actualTeam2.players), testCase.expectedTeam2Length)
				}
			})
		}
	})
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
