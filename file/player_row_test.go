package file

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewPlayerRow(t *testing.T) {
	var tests = []struct {
		name          string
		myClubId      string
		playerName    string
		wantErr       bool
		expectedError string
		expectedRow   PlayerRow
	}{
		{
			name:       "Valid player",
			myClubId:   "123",
			playerName: "Matti Meikäläinen",
			wantErr:    false,
			expectedRow: PlayerRow{
				MyClubId: 123,
				Name:     "Matti Meikäläinen",
			},
		},
		{
			name:          "Invalid missing name",
			myClubId:      "123",
			playerName:    "",
			wantErr:       true,
			expectedError: "player name can't be empty",
			expectedRow:   PlayerRow{},
		},
		{
			name:          "Invalid not integer MyClubId",
			myClubId:      "asdasd",
			playerName:    "Matti Meikäläinen",
			wantErr:       true,
			expectedError: "unable to convert MyClubId to integer",
			expectedRow:   PlayerRow{},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			playerRow, err := newPlayerRow(testCase.myClubId, testCase.playerName)

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
			if !reflect.DeepEqual(playerRow, testCase.expectedRow) {
				t.Errorf("player row doesn't match: got [%v] want [%v]", playerRow, testCase.expectedRow)
			}

		})

	}
}
