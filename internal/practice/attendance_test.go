package practice

import (
	"database/sql"
	"slices"
	"strings"
	"testing"

	"github.com/Zigelzi/go-tiimit/internal/db"
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
		name                     string
		setupPractice            func() *Practice
		status                   AttendanceStatus
		wantErr                  bool
		expectedErr              string
		expectedPlayers          []db.Player
		expectedUnknownPlayerIds []int
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
			expectedPlayers: []db.Player{
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
			expectedPlayers: []db.Player{},
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
			expectedPlayers: []db.Player{
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
			expectedPlayers: []db.Player{},
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
			expectedPlayers: []db.Player{
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
			expectedPlayers: []db.Player{},
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
			wantErr: false,
			expectedPlayers: []db.Player{
				{MyClubId: 1001, Name: "Teppo Teikäläinen"},
			},
			expectedUnknownPlayerIds: []int{9999},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			p := testCase.setupPractice()
			players, unknownPlayerIds, err := p.GetPlayersByStatus(testCase.status, mockPlayerGetter)
			if testCase.wantErr {
				if err == nil {
					t.Errorf("error is missing: got [nil] want [%s]", testCase.expectedErr)
					// return
				}
				if !strings.Contains(err.Error(), testCase.expectedErr) {
					t.Errorf("error contents don't match: got [%s] want [%s]", err.Error(), testCase.expectedErr)
				}
				// return
			} else {
				if err != nil {
					t.Errorf("unexpected error: got [%s] want [nil]", err)
				}
			}

			if len(players) != len(testCase.expectedPlayers) {
				t.Errorf("number of players in status [%s] doesn't match: got [%d] want [%d]", testCase.status, len(players), len(testCase.expectedPlayers))
			}
			assertEqualPlayers(t, players, testCase.expectedPlayers)

			if len(unknownPlayerIds) != len(testCase.expectedUnknownPlayerIds) {
				t.Errorf("number of unknown player IDs don't match: got [%d] want [%d]", len(unknownPlayerIds), len(testCase.expectedUnknownPlayerIds))
			}

			assertEqualUnknownPlayerIds(t, unknownPlayerIds, testCase.expectedUnknownPlayerIds)
		})
	}
}

func mockPlayerGetter(id int64) (db.Player, error) {
	players := map[int64]db.Player{
		1000: {MyClubId: 1000, Name: "Matti Meikäläinen"},
		1001: {MyClubId: 1001, Name: "Teppo Teikäläinen"},
		1002: {MyClubId: 1002, Name: "Seppo Seikäläinen"},
		1003: {MyClubId: 1003, Name: "Kati Kaapu"},
		1004: {MyClubId: 1004, Name: "Tero Taapu"},
		1005: {MyClubId: 1005, Name: "Lauri Laatu"},
	}

	queriedPlayer, exists := players[id]
	if !exists {
		return db.Player{}, sql.ErrNoRows
	}

	return queriedPlayer, nil
}

func assertEqualPlayers(t *testing.T, got, want []db.Player) {
	t.Helper()

	for _, wantPlayer := range want {
		if !slices.Contains(got, wantPlayer) {
			t.Errorf("wanted player [%d] %s missing from got players", wantPlayer.MyClubId, wantPlayer.Name)
			t.Errorf("got %v", got)
		}
	}
}

func assertEqualUnknownPlayerIds(t *testing.T, got, want []int) {
	t.Helper()

	for _, wantMyClubId := range want {
		if !slices.Contains(got, wantMyClubId) {
			t.Errorf("unknown MyClubId not found: want [%d]", wantMyClubId)
			t.Errorf("got MyClubIds: %v", got)
		}
	}
}
