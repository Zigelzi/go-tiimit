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

func TestNewAttendancePlayerRow(t *testing.T) {
	var tests = []struct {
		name             string
		myClubId         string
		playerName       string
		attendanceStatus string
		wantErr          bool
		expectedError    string
		expectedRow      AttendancePlayerRow
	}{
		{
			name:             "Valid information AttendanceIn",
			myClubId:         "123",
			playerName:       "Matti Meikäläinen",
			attendanceStatus: "Osallistuu",
			wantErr:          false,
			expectedRow: AttendancePlayerRow{
				PlayerRow: PlayerRow{
					MyClubId: 123,
					Name:     "Matti Meikäläinen",
				},
				Attendance: "Osallistuu",
			},
		},
		{
			name:             "Valid information AttendanceOut",
			myClubId:         "123",
			playerName:       "Matti Meikäläinen",
			attendanceStatus: "Ei osallistu",
			wantErr:          false,
			expectedRow: AttendancePlayerRow{
				PlayerRow: PlayerRow{
					MyClubId: 123,
					Name:     "Matti Meikäläinen",
				},
				Attendance: "Ei osallistu",
			},
		},
		{
			name:             "Valid information AttendanceUnknown",
			myClubId:         "123",
			playerName:       "Matti Meikäläinen",
			attendanceStatus: "Ei vastausta",
			wantErr:          false,
			expectedRow: AttendancePlayerRow{
				PlayerRow: PlayerRow{
					MyClubId: 123,
					Name:     "Matti Meikäläinen",
				},
				Attendance: "Ei vastausta",
			},
		},
		{
			name:             "Invalid information Unknown attendance status",
			myClubId:         "123",
			playerName:       "Matti Meikäläinen",
			attendanceStatus: "Wutisdis",
			wantErr:          true,
			expectedError:    "unknown attendance status",
			expectedRow:      AttendancePlayerRow{},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			attendancePlayerRow, err := newAttendancePlayerRow(testCase.myClubId, testCase.playerName, testCase.attendanceStatus)

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
			if !reflect.DeepEqual(attendancePlayerRow, testCase.expectedRow) {
				t.Errorf("player row doesn't match: got [%v] want [%v]", attendancePlayerRow, testCase.expectedRow)
			}

		})

	}
}
