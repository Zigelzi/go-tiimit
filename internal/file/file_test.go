package file

import (
	"reflect"
	"strings"
	"testing"
	"time"
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

func TestParsingEmptyRowsReturnsError(t *testing.T) {
	rows := [][]string{}
	_, actualErr := parseAttendanceRows(rows)
	expectedErr := "no players found from the rows"
	if actualErr == nil {
		t.Fatalf("error is missing: got [nil] want [%s]", expectedErr)
	}
	if !strings.Contains(actualErr.Error(), expectedErr) {
		t.Errorf("error contents don't match: got [%s] want [%s]", actualErr.Error(), expectedErr)
	}
}

func TestParsingValidRowsReturnsCorrectNumberOfAttendanceRows(t *testing.T) {
	rows := [][]string{
		{"1111", "Matti Meikäläinen", "Mies", "Osallistuu"},
		{"1111", "Seppo Seikäläinen", "Mies", "Osallistuu"},
		{"1111", "Kaija Kaarela", "Nainen", "Osallistuu"},
	}

	actualAttendanceRows, actualErr := parseAttendanceRows(rows)

	if actualErr != nil {
		t.Errorf("unexpected error: got [%v] want [nil]", actualErr.Error())
	}

	if len(actualAttendanceRows) != len(rows) {
		t.Errorf("number of rows and AttendancePlayers don't match: got [%d] want [%d]", len(actualAttendanceRows), len(rows))
	}
}

func TestParsingRowsWithMissingFieldsReturnsError(t *testing.T) {
	testCases := []struct {
		name        string
		rows        [][]string
		expectedErr string
	}{
		{
			name: "returns error when MyClubId is missing",
			rows: [][]string{
				{"", "Matti Meikäläinen", "Mies", "Osallistuu"},
				{"", "Seppo Seikäläinen", "Mies", "Osallistuu"},
				{"", "Kaija Kaarela", "Nainen", "Osallistuu"},
			},
			expectedErr: "unable to create new player row",
		},
		{
			name: "returns error when name is missing",
			rows: [][]string{
				{"1111", "", "Mies", "Osallistuu"},
				{"2222", "", "Mies", "Osallistuu"},
				{"3333", "", "Nainen", "Osallistuu"},
			},
			expectedErr: "unable to create new player row",
		},
		{
			name: "returns error when attendance is missing",
			rows: [][]string{
				{"1111", "Matti Meikäläinen", "Mies", ""},
				{"2222", "Seppo Seikäläinen", "Mies", ""},
				{"3333", "Kaija Kaarela", "Nainen", ""},
			},
			expectedErr: "unable to create new player row",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, actualErr := parseAttendanceRows(testCase.rows)
			if actualErr == nil {
				t.Fatalf("error is missing: got [nil] want [%s]", testCase.expectedErr)
			}
			if !strings.Contains(actualErr.Error(), testCase.expectedErr) {
				t.Errorf("error contents don't match: got [%s] want [%s]", actualErr.Error(), testCase.expectedErr)
			}
		})
	}

}
func TestParsingValidRowsReturnsEqualAttendanceRows(t *testing.T) {

	var testCases = []struct {
		name                   string
		rows                   [][]string
		expectedAttendanceRows []AttendancePlayerRow
	}{
		{
			name: "rows with attendance in",
			rows: [][]string{
				{"1111", "Matti Meikäläinen", "Mies", "Osallistuu"},
				{"2222", "Seppo Seikäläinen", "Mies", "Osallistuu"},
				{"3333", "Kaija Kaarela", "Nainen", "Osallistuu"},
			},
			expectedAttendanceRows: []AttendancePlayerRow{
				{PlayerRow: PlayerRow{MyClubId: 1111, Name: "Matti Meikäläinen"}, Attendance: "Osallistuu"},
				{PlayerRow: PlayerRow{MyClubId: 2222, Name: "Seppo Seikäläinen"}, Attendance: "Osallistuu"},
				{PlayerRow: PlayerRow{MyClubId: 3333, Name: "Kaija Kaarela"}, Attendance: "Osallistuu"},
			},
		},
		{
			name: "rows with attendance out",
			rows: [][]string{
				{"1111", "Matti Meikäläinen", "Mies", "Ei osallistu"},
				{"2222", "Seppo Seikäläinen", "Mies", "Ei osallistu"},
				{"3333", "Kaija Kaarela", "Nainen", "Ei osallistu"},
			},
			expectedAttendanceRows: []AttendancePlayerRow{
				{PlayerRow: PlayerRow{MyClubId: 1111, Name: "Matti Meikäläinen"}, Attendance: "Ei osallistu"},
				{PlayerRow: PlayerRow{MyClubId: 2222, Name: "Seppo Seikäläinen"}, Attendance: "Ei osallistu"},
				{PlayerRow: PlayerRow{MyClubId: 3333, Name: "Kaija Kaarela"}, Attendance: "Ei osallistu"},
			},
		},
		{
			name: "rows with unknown attendance",
			rows: [][]string{
				{"1111", "Matti Meikäläinen", "Mies", "Ei vastausta"},
				{"2222", "Seppo Seikäläinen", "Mies", "Ei vastausta"},
				{"3333", "Kaija Kaarela", "Nainen", "Ei vastausta"},
			},
			expectedAttendanceRows: []AttendancePlayerRow{
				{PlayerRow: PlayerRow{MyClubId: 1111, Name: "Matti Meikäläinen"}, Attendance: "Ei vastausta"},
				{PlayerRow: PlayerRow{MyClubId: 2222, Name: "Seppo Seikäläinen"}, Attendance: "Ei vastausta"},
				{PlayerRow: PlayerRow{MyClubId: 3333, Name: "Kaija Kaarela"}, Attendance: "Ei vastausta"},
			},
		},
		{
			name: "rows with all attendance statuses",
			rows: [][]string{
				{"1111", "Matti Meikäläinen", "Mies", "Osallistuu"},
				{"2222", "Seppo Seikäläinen", "Mies", "Ei osallistu"},
				{"3333", "Kaija Kaarela", "Nainen", "Ei vastausta"},
			},
			expectedAttendanceRows: []AttendancePlayerRow{
				{PlayerRow: PlayerRow{MyClubId: 1111, Name: "Matti Meikäläinen"}, Attendance: "Osallistuu"},
				{PlayerRow: PlayerRow{MyClubId: 2222, Name: "Seppo Seikäläinen"}, Attendance: "Ei osallistu"},
				{PlayerRow: PlayerRow{MyClubId: 3333, Name: "Kaija Kaarela"}, Attendance: "Ei vastausta"},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualAttendanceRows, actualErr := parseAttendanceRows(testCase.rows)

			if actualErr != nil {
				t.Errorf("unexpected error: got [%v] want [nil]", actualErr.Error())
			}

			for i, actualAttendanceRow := range actualAttendanceRows {
				if !reflect.DeepEqual(actualAttendanceRow, testCase.expectedAttendanceRows[i]) {
					t.Errorf("row and attendance row don't match: got [%v] want [%v]", actualAttendanceRow, testCase.expectedAttendanceRows[i])
				}
			}

		})
	}
}

func TestParseDate(t *testing.T) {
	type testCase struct {
		name         string
		fileName     string
		expectedDate time.Time
		expectedErr  string
	}
	t.Run("returns date type when file name contains a date", func(t *testing.T) {
		testCases := []testCase{
			{
				"first day of the year",
				"vjs-Kuntofutis-2025-01-01-VJS Kuntofutis.xlsx",
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				"",
			},
			{
				"last day of the year",
				"vjs-Kuntofutis-2025-12-31-VJS Kuntofutis.xlsx",
				time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
				"",
			},
			{
				"leap year",
				"vjs-Kuntofutis-2024-02-28-VJS Kuntofutis.xlsx",
				time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
				"",
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				actualDate, err := parseDate(testCase.fileName)

				if err != nil {
					t.Errorf("unexpected error: got [%v] want [nil]", err)
				}

				if actualDate != testCase.expectedDate {
					t.Errorf("parsed dates don't match: got [%v] want [%v]", actualDate, testCase.expectedDate)
				}
			})
		}
	})

	t.Run("returns error when file name has invalid or missing date", func(t *testing.T) {
		testCases := []testCase{
			{
				"missing date",
				"vjs-Kuntofutis-VJS Kuntofutis.xlsx",
				time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC),
				"failed to find date",
			},
			{
				"date that doesn't exist",
				"vjs-Kuntofutis-2025-01-33-VJS Kuntofutis.xlsx",
				time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC),
				"failed to parse date",
			},
			{
				"date formatted in yyyymmdd",
				"vjs-Kuntofutis-20250131-VJS Kuntofutis.xlsx",
				time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC),
				"failed to find date",
			},
			{
				"date formatted in yyyy/mm/dd",
				"vjs-Kuntofutis-2025/01/31-VJS Kuntofutis.xlsx",
				time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC),
				"failed to find date",
			},
		}

		for _, testCase := range testCases {
			actualDate, err := parseDate(testCase.fileName)

			if err == nil {
				t.Errorf("missing error: got [nil] want [%v]", testCase.expectedErr)
			}

			if actualDate != testCase.expectedDate {
				t.Errorf("parsed dates don't match: got [%v] want [%v]", actualDate, testCase.expectedDate)
			}

			if !strings.Contains(err.Error(), testCase.expectedErr) {
				t.Errorf("errors don't match: got [%s] want [%s]", err, testCase.expectedErr)
			}
		}
	})
}

func TestFindDate(t *testing.T) {
	type testCase struct {
		name        string
		inputStr    string
		expectedStr string
		expectedErr string
	}

	t.Run("returns string matching yyyy-mm-dd format when text contains such", func(t *testing.T) {
		testCases := []testCase{
			{
				"date in the beginning",
				"2025-01-31-vjs-Kuntofutis-VJS Kuntofutis.xlsx",
				"2025-01-31",
				"",
			},
			{
				"date in the end",
				"vjs-Kuntofutis-VJS Kuntofutis-2025-01-31.xlsx",
				"2025-01-31",
				"",
			},
			{
				"date in the middle",
				"vjs-Kuntofutis-2025-01-31-VJS Kuntofutis.xlsx",
				"2025-01-31",
				"",
			},
		}

		for _, testCase := range testCases {
			actualStr, err := findDate(testCase.inputStr)

			if err != nil {
				t.Errorf("unexpected error: got [%v] want [nil]", err)
			}

			if actualStr != testCase.expectedStr {
				t.Errorf("returned string doesn't match: got [%s] want [%s]", actualStr, testCase.expectedStr)
			}
		}
	})

	t.Run("returns error when string doesn't contain yyyy-mm-dd date", func(t *testing.T) {
		testCases := []testCase{
			{
				"missing date",
				"vjs-KuntofutisVJS Kuntofutis.xlsx",
				"",
				"vjs-KuntofutisVJS Kuntofutis.xlsx doesn't contain date with pattern yyyy-mm-dd",
			},
			{
				"date in yyyy/mm/dd",
				"vjs-Kuntofutis-2025/01/01-VJS Kuntofutis.xlsx",
				"",
				"vjs-Kuntofutis-2025/01/01-VJS Kuntofutis.xlsx doesn't contain date with pattern yyyy-mm-dd",
			},
			{
				"date in dd-mm-yyyy",
				"vjs-Kuntofutis-01-01-2025-VJS Kuntofutis.xlsx",
				"",
				"vjs-Kuntofutis-01-01-2025-VJS Kuntofutis.xlsx doesn't contain date with pattern yyyy-mm-dd",
			},
		}

		for _, testCase := range testCases {
			actualStr, err := findDate(testCase.inputStr)
			if err == nil {
				t.Errorf("missing error: got [nil] want [%s]", testCase.expectedErr)
			}

			if actualStr != testCase.expectedStr {
				t.Errorf("returned string doesn't match: got [%s] want [%s]", actualStr, testCase.expectedStr)
			}

			if !strings.Contains(err.Error(), testCase.expectedErr) {
				t.Errorf("error doesn't match: got [%s] want [%s]", err, testCase.expectedErr)
			}
		}
	})
}
