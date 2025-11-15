package team

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Zigelzi/go-tiimit/internal/player"
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
				if len(actualTeam1.Players) != testCase.expectedTeam1Length {
					t.Errorf("number of players in %s don't match: got [%d] want [%d]", actualTeam1.Name, len(actualTeam1.Players), testCase.expectedTeam1Length)
				}
				if len(actualTeam2.Players) != testCase.expectedTeam2Length {
					t.Errorf("number of players in %s don't match: got [%d] want [%d]", actualTeam2.Name, len(actualTeam2.Players), testCase.expectedTeam2Length)
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
				Name: "Team 1",
				Players: []player.Player{
					{MyClubId: 1001, Name: "Teppo Teikäläinen", IsGoalie: true},
					{MyClubId: 1003, Name: "Saija Siirappi", IsGoalie: true},
					{MyClubId: 1005, Name: "Marja Myyry", IsGoalie: false},
					{MyClubId: 1007, Name: "Kalevi Kuurula", IsGoalie: false},
				},
			},
			expectedTeam2: Team{
				Name: "Team 2",
				Players: []player.Player{
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

			if len(team1.Players) != len(testCase.expectedTeam1.Players) {
				t.Errorf("%s number of players don't match: got [%d] want [%d]", team1.Name, len(team1.Players), len(testCase.expectedTeam1.Players))
			}

			if len(team2.Players) != len(testCase.expectedTeam2.Players) {
				t.Errorf("%s number of players don't match: got [%d] want [%d]", team2.Name, len(team2.Players), len(testCase.expectedTeam2.Players))
			}

			if !reflect.DeepEqual(team1.Players, testCase.expectedTeam1.Players) {
				t.Errorf("%s players don't match: got [%v] want [%v]", team1.Name, team1.Players, testCase.expectedTeam1.Players)
			}
			if !reflect.DeepEqual(team2.Players, testCase.expectedTeam2.Players) {
				t.Errorf("%s players don't match: got [%v] want [%v]", team2.Name, team2.Players, testCase.expectedTeam2.Players)
			}
		})
	}
}
