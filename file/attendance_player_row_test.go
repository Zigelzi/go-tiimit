package file

import (
	"reflect"
	"strings"
	"testing"
)

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
