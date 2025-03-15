package practice

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Zigelzi/go-tiimit/player"
)

func TestAddPlayer(t *testing.T) {
	var tests = []struct {
		name           string
		myClubId       int
		status         string
		wantErr        bool
		expectedErr    string
		expectedStatus AttendanceStatus
	}{
		{
			name:           "Add player with attendance in",
			myClubId:       1400,
			status:         "Osallistuu",
			wantErr:        false,
			expectedStatus: AttendanceIn,
		},
		{
			name:           "Add player with attendance out",
			myClubId:       1400,
			status:         "Ei osallistu",
			wantErr:        false,
			expectedStatus: AttendanceOut,
		},
		{
			name:           "Add player with attendance unknown",
			myClubId:       1400,
			status:         "Ei vastausta",
			wantErr:        false,
			expectedStatus: AttendanceUnknown,
		},
		{
			name:        "Return error for invalid status",
			myClubId:    1400,
			status:      "quo",
			wantErr:     true,
			expectedErr: "invalid status parsed for for player with MyClub ID:",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			p := New()
			err := p.AddPlayer(testCase.myClubId, testCase.status)

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

			status := p.AttendingPlayers[testCase.myClubId]
			if status != testCase.expectedStatus {
				t.Errorf("statuses don't match: got [%v] want [%v]", status, testCase.expectedStatus)
			}
		})
	}
}

func TestGetPlayersByStatus(t *testing.T) {
	var tests = []struct {
		name            string
		setupPractice   func() *Practice
		status          AttendanceStatus
		wantErr         bool
		expectedErr     string
		expectedPlayers []player.Player
	}{
		{
			name: "Return attending players from practice with players in all statuses",
			setupPractice: func() *Practice {
				p := New()
				p.AddPlayer(1000, "Osallistuu")
				p.AddPlayer(1001, "Osallistuu")
				p.AddPlayer(1002, "Ei osallistu")
				p.AddPlayer(1003, "Ei osallistu")
				p.AddPlayer(1004, "Ei vastausta")
				p.AddPlayer(1005, "Ei vastausta")
				return &p
			},
			status:  AttendanceIn,
			wantErr: false,
			expectedPlayers: []player.Player{
				{MyClubId: 1000, Name: "Matti Meikäläinen"},
				{MyClubId: 1001, Name: "Teppo Teikäläinen"},
			},
		},
		{
			name: "Return no players from practice without attending players",
			setupPractice: func() *Practice {
				p := New()
				p.AddPlayer(1002, "Ei osallistu")
				p.AddPlayer(1003, "Ei osallistu")
				p.AddPlayer(1004, "Ei vastausta")
				p.AddPlayer(1005, "Ei vastausta")
				return &p
			},
			status:          AttendanceIn,
			wantErr:         false,
			expectedPlayers: []player.Player{},
		},
		{
			name: "Return unknown players from practice with players in all statuses",
			setupPractice: func() *Practice {
				p := New()
				p.AddPlayer(1000, "Osallistuu")
				p.AddPlayer(1001, "Osallistuu")
				p.AddPlayer(1002, "Ei osallistu")
				p.AddPlayer(1003, "Ei osallistu")
				p.AddPlayer(1004, "Ei vastausta")
				p.AddPlayer(1005, "Ei vastausta")
				return &p
			},
			status:  AttendanceUnknown,
			wantErr: false,
			expectedPlayers: []player.Player{
				{MyClubId: 1004, Name: "Tero Taapu"},
				{MyClubId: 1005, Name: "Lauri Laatu"},
			},
		},
		{
			name: "Return no players from practice without unknown players",
			setupPractice: func() *Practice {
				p := New()
				p.AddPlayer(1000, "Osallistuu")
				p.AddPlayer(1001, "Osallistuu")
				p.AddPlayer(1002, "Ei osallistu")
				p.AddPlayer(1003, "Ei osallistu")
				return &p
			},
			status:          AttendanceUnknown,
			wantErr:         false,
			expectedPlayers: []player.Player{},
		},
		{
			name: "Return out players from practice with players in all statuses",
			setupPractice: func() *Practice {
				p := New()
				p.AddPlayer(1000, "Osallistuu")
				p.AddPlayer(1001, "Osallistuu")
				p.AddPlayer(1002, "Ei osallistu")
				p.AddPlayer(1003, "Ei osallistu")
				p.AddPlayer(1004, "Ei vastausta")
				p.AddPlayer(1005, "Ei vastausta")
				return &p
			},
			status:  AttendanceOut,
			wantErr: false,
			expectedPlayers: []player.Player{
				{MyClubId: 1002, Name: "Seppo Seikäläinen"},
				{MyClubId: 1003, Name: "Kati Kaapu"},
			},
		},
		{
			name: "Return no players from practice without not attending players",
			setupPractice: func() *Practice {
				p := New()
				p.AddPlayer(1000, "Osallistuu")
				p.AddPlayer(1001, "Osallistuu")
				p.AddPlayer(1004, "Ei vastausta")
				p.AddPlayer(1005, "Ei vastausta")
				return &p
			},
			status:          AttendanceOut,
			wantErr:         false,
			expectedPlayers: []player.Player{},
		},
		{
			name: "Return no players from practice without players",
			setupPractice: func() *Practice {
				p := New()

				return &p
			},
			status:          AttendanceOut,
			wantErr:         false,
			expectedPlayers: []player.Player{},
		},
		{
			name: "Return attending players from practice and errors with players in all statuses",
			setupPractice: func() *Practice {
				p := New()
				p.AddPlayer(9999, "Osallistuu") // Unknown player in mock
				p.AddPlayer(1001, "Osallistuu")
				p.AddPlayer(1002, "Ei osallistu")
				p.AddPlayer(1003, "Ei osallistu")
				p.AddPlayer(1004, "Ei vastausta")
				p.AddPlayer(1005, "Ei vastausta")
				return &p
			}, status: AttendanceIn,
			wantErr:     true,
			expectedErr: "unable to get player",
			expectedPlayers: []player.Player{
				{MyClubId: 1001, Name: "Teppo Teikäläinen"},
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			p := testCase.setupPractice()
			players, err := p.GetPlayersByStatus(testCase.status, mockPlayerGetter)
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

			if len(players) != len(testCase.expectedPlayers) {
				t.Errorf("number of players in status [%s] doesn't match: got [%d] want [%d]", testCase.status, len(players), len(testCase.expectedPlayers))
			}
			if !reflect.DeepEqual(players, testCase.expectedPlayers) {
				t.Errorf("players in status [%s] don't match", testCase.status)
				t.Errorf("got: %v", players)
				t.Errorf("want: %v", testCase.expectedPlayers)
			}
		})
	}
}

func mockPlayerGetter(id int64) (player.Player, error) {
	players := map[int64]player.Player{
		1000: {MyClubId: 1000, Name: "Matti Meikäläinen"},
		1001: {MyClubId: 1001, Name: "Teppo Teikäläinen"},
		1002: {MyClubId: 1002, Name: "Seppo Seikäläinen"},
		1003: {MyClubId: 1003, Name: "Kati Kaapu"},
		1004: {MyClubId: 1004, Name: "Tero Taapu"},
		1005: {MyClubId: 1005, Name: "Lauri Laatu"},
	}

	queriedPlayer, exists := players[id]
	if !exists {
		return player.Player{}, fmt.Errorf("unable to query player with MyClub ID")
	}

	return queriedPlayer, nil
}
