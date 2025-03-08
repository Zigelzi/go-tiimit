package practice

import (
	"strings"
	"testing"
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
