package file

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewClubPlayerRow(t *testing.T) {
	var tests = []struct {
		name          string
		myClubId      string
		playerName    string
		runPower      string
		ballHandling  string
		wantErr       bool
		expectedError string
		expectedRow   ClubPlayerRow
	}{
		{
			name:         "Valid perfect both skills",
			myClubId:     "123",
			playerName:   "Matti Meikäläinen",
			runPower:     "10",
			ballHandling: "10",
			wantErr:      false,
			expectedRow: ClubPlayerRow{
				PlayerRow: PlayerRow{
					MyClubId: 123,
					Name:     "Matti Meikäläinen",
				},
				RunPower:     10.0,
				BallHandling: 10.0,
			},
		},
		{
			name:         "Valid zero run power",
			myClubId:     "123",
			playerName:   "Matti Meikäläinen",
			runPower:     "0",
			ballHandling: "10",
			wantErr:      false,
			expectedRow: ClubPlayerRow{
				PlayerRow: PlayerRow{
					MyClubId: 123,
					Name:     "Matti Meikäläinen",
				},
				RunPower:     0.0,
				BallHandling: 10.0,
			},
		},
		{
			name:         "Valid zero ball handling",
			myClubId:     "123",
			playerName:   "Matti Meikäläinen",
			runPower:     "10",
			ballHandling: "0",
			wantErr:      false,
			expectedRow: ClubPlayerRow{
				PlayerRow: PlayerRow{
					MyClubId: 123,
					Name:     "Matti Meikäläinen",
				},
				RunPower:     10.0,
				BallHandling: 0.0,
			},
		},
		{
			name:          "Invalid run power parsing failure",
			myClubId:      "123",
			playerName:    "Matti Meikäläinen",
			runPower:      "notfloat",
			ballHandling:  "10",
			wantErr:       true,
			expectedError: "failed to parse run power",
			expectedRow:   ClubPlayerRow{},
		},
		{
			name:          "Invalid ball handling parsing failure",
			myClubId:      "123",
			playerName:    "Matti Meikäläinen",
			runPower:      "10",
			ballHandling:  "notfloat",
			wantErr:       true,
			expectedError: "failed to parse ball handling",
			expectedRow:   ClubPlayerRow{},
		},
		{
			name:          "Invalid too low run power",
			myClubId:      "123",
			playerName:    "Matti Meikäläinen",
			runPower:      "-10",
			ballHandling:  "10",
			wantErr:       true,
			expectedError: "run power too low: needs to be between 0-10",
			expectedRow:   ClubPlayerRow{},
		},
		{
			name:          "Invalid too high run power",
			myClubId:      "123",
			playerName:    "Matti Meikäläinen",
			runPower:      "100",
			ballHandling:  "10",
			wantErr:       true,
			expectedError: "run power too high: needs to be between 0-10",
			expectedRow:   ClubPlayerRow{},
		},
		{
			name:          "Invalid too low ball handling",
			myClubId:      "123",
			playerName:    "Matti Meikäläinen",
			runPower:      "10",
			ballHandling:  "-10",
			wantErr:       true,
			expectedError: "ball handling too low: needs to be between 0-10",
			expectedRow:   ClubPlayerRow{},
		},
		{
			name:          "Invalid too high ball handling",
			myClubId:      "123",
			playerName:    "Matti Meikäläinen",
			runPower:      "10",
			ballHandling:  "100",
			wantErr:       true,
			expectedError: "ball handling too high: needs to be between 0-10",
			expectedRow:   ClubPlayerRow{},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			clubPlayerRow, err := newClubPlayerRow(testCase.myClubId, testCase.playerName, testCase.runPower, testCase.ballHandling)
			if testCase.wantErr {
				if err == nil {
					t.Errorf("error is missing: got [nil] want [%s]", testCase.expectedError)
					return
				}
				if !strings.Contains(err.Error(), testCase.expectedError) {
					t.Errorf("error contents don't match: got [%s] want [%s]", err.Error(), testCase.expectedError)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: got [%s] want [nil]", err)
			}
			if !reflect.DeepEqual(clubPlayerRow, testCase.expectedRow) {
				t.Errorf("player row doesn't match: got [%v] want [%v]", clubPlayerRow, testCase.expectedRow)
			}
		})
	}
}
