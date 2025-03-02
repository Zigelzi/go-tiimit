package practice

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Zigelzi/go-tiimit/player"
)

func TestNewAttendingPlayer(t *testing.T) {
	var tests = []struct {
		name                    string
		player                  player.Player
		attendanceStatus        string
		wantErr                 bool
		expectedErr             string
		expectedAttendingPlayer AttendingPlayer
	}{
		{
			name: "Valid attendance in",
			player: player.Player{
				MyClubId: 123,
				Name:     "Matti Meikäläinen",
			},
			attendanceStatus: "Osallistuu",
			wantErr:          false,
			expectedAttendingPlayer: AttendingPlayer{
				Player: player.Player{
					MyClubId: 123,
					Name:     "Matti Meikäläinen",
				},
				Attendance: AttendanceIn,
			},
		},
		{
			name: "Valid attendance out",
			player: player.Player{
				MyClubId: 123,
				Name:     "Matti Meikäläinen",
			},
			attendanceStatus: "Ei osallistu",
			wantErr:          false,
			expectedAttendingPlayer: AttendingPlayer{
				Player: player.Player{
					MyClubId: 123,
					Name:     "Matti Meikäläinen",
				},
				Attendance: AttendanceOut,
			},
		},
		{
			name: "Valid attendance unknown",
			player: player.Player{
				MyClubId: 123,
				Name:     "Matti Meikäläinen",
			},
			attendanceStatus: "Ei vastausta",
			wantErr:          false,
			expectedAttendingPlayer: AttendingPlayer{
				Player: player.Player{
					MyClubId: 123,
					Name:     "Matti Meikäläinen",
				},
				Attendance: AttendanceUnknown,
			},
		},
		{
			name: "Invalid attendance not known",
			player: player.Player{
				MyClubId: 123,
				Name:     "Matti Meikäläinen",
			},
			attendanceStatus: "not known",
			wantErr:          true,
			expectedErr:      "invalid attendance status",
			expectedAttendingPlayer: AttendingPlayer{
				Player:     player.Player{},
				Attendance: AttendanceInvalid,
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			attendingPlayer, err := NewAttendingPlayer(testCase.player, testCase.attendanceStatus)

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
			if !reflect.DeepEqual(attendingPlayer, testCase.expectedAttendingPlayer) {
				t.Errorf("attending players doesn't match: got [%v] want [%v]", attendingPlayer, testCase.expectedAttendingPlayer)
			}
		})
	}
}
