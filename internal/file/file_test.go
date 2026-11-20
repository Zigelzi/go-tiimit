package file

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

// finnishHeader is the header row of a Finnish MyClub attendance report in
// the column order MyClub uses by default.
var finnishHeader = []string{"Jäsennumero", "Suku- ja etunimi", "Sukupuoli", "Osallistuu"}

// report returns the rows of an attendance report sheet: the title rows a
// real export has above the header, then the header, then the players.
func report(header []string, players ...[]string) [][]string {
	rows := [][]string{{"Testiseura Kuntofutis  (16.11.2026)"}, {"Kuntofutis"}, {}, header}
	return append(rows, players...)
}

func attending(id int, name string) AttendancePlayerRow {
	return AttendancePlayerRow{PlayerRow: PlayerRow{MyclubID: id, Name: name}, Attendance: AttendanceIn}
}

func TestParseAttendanceSheet(t *testing.T) {
	allColumns := []string{"Jäsennumero", "Suku- ja etunimi", "Osallistuu"}

	tests := []struct {
		name        string
		rows        [][]string
		want        []AttendancePlayerRow
		wantSkipped []SkippedRow
		wantMissing []string // the columns the error names as missing
		wantErr     error    // a specific error the result must wrap
	}{
		{
			name: "reads players below the header in MyClub's default layout",
			rows: report(finnishHeader,
				[]string{"9001", "Testi Pelaaja", "Mies", "Osallistuu"},
				[]string{"9002", "Keksitty Nimi", "Nainen", "Ei osallistu"},
				[]string{"9003", "Kuvitteellinen Henkilö", "Mies", "Ei vastausta"},
			),
			want: []AttendancePlayerRow{
				attending(9001, "Testi Pelaaja"),
				{PlayerRow: PlayerRow{MyclubID: 9002, Name: "Keksitty Nimi"}, Attendance: AttendanceOut},
				{PlayerRow: PlayerRow{MyclubID: 9003, Name: "Kuvitteellinen Henkilö"}, Attendance: AttendanceUnknown},
			},
		},
		{
			name: "finds the header below an extra title row",
			rows: append([][]string{{"Lisärivi"}}, report(finnishHeader,
				[]string{"9001", "Testi Pelaaja", "Mies", "Osallistuu"},
			)...),
			want: []AttendancePlayerRow{attending(9001, "Testi Pelaaja")},
		},
		{
			name: "finds the header as the first row",
			rows: [][]string{
				finnishHeader,
				{"9001", "Testi Pelaaja", "Mies", "Osallistuu"},
			},
			want: []AttendancePlayerRow{attending(9001, "Testi Pelaaja")},
		},
		{
			name: "reads columns in any order and ignores extra columns",
			rows: report([]string{"Osallistuu", "Jäsennumero", "Lempiväri", "Suku- ja etunimi"},
				[]string{"Osallistuu", "9001", "Sininen", "Testi Pelaaja"},
				[]string{"Osallistuu", "9002", "Punainen", "Keksitty Nimi"},
			),
			want: []AttendancePlayerRow{attending(9001, "Testi Pelaaja"), attending(9002, "Keksitty Nimi")},
		},
		{
			name: "ignores spaces around header names and cell values",
			rows: report([]string{" Jäsennumero ", "Suku- ja etunimi ", " Osallistuu"},
				[]string{" 9001 ", " Testi Pelaaja ", " Osallistuu "},
			),
			want: []AttendancePlayerRow{attending(9001, "Testi Pelaaja")},
		},
		{
			name:        "names the MyClub ID column when it is missing",
			rows:        report([]string{"Suku- ja etunimi", "Osallistuu"}, []string{"Testi Pelaaja", "Osallistuu"}),
			wantMissing: []string{"Jäsennumero"},
		},
		{
			name:        "names the name column when it is missing",
			rows:        report([]string{"Jäsennumero", "Osallistuu"}, []string{"9001", "Osallistuu"}),
			wantMissing: []string{"Suku- ja etunimi"},
		},
		{
			name:        "names the attendance column when it is missing",
			rows:        report([]string{"Jäsennumero", "Suku- ja etunimi"}, []string{"9001", "Testi Pelaaja"}),
			wantMissing: []string{"Osallistuu"},
		},
		{
			name: "names all three columns for an English report",
			rows: report([]string{"Member number", "Last and first name", "Gender", "Coming"},
				[]string{"9001", "Testi Pelaaja", "Male", "Coming"},
			),
			wantMissing: allColumns,
		},
		{
			name:        "names all three columns for an empty sheet",
			rows:        [][]string{},
			wantMissing: allColumns,
		},
		{
			name:        "names all three columns for a sheet with fewer rows than the usual header position",
			rows:        [][]string{{"Testiseura Kuntofutis  (16.11.2026)"}},
			wantMissing: allColumns,
		},
		{
			name:    "reports no players for a header without player rows",
			rows:    report(finnishHeader),
			wantErr: ErrNoPlayers,
		},
		{
			name: "treats a row without an attendance cell as not attending, without reporting it",
			rows: report(finnishHeader,
				[]string{"9001", "Testi Pelaaja"},
				[]string{"9002", "Keksitty Nimi", "Nainen", "Osallistuu"},
			),
			want: []AttendancePlayerRow{attending(9002, "Keksitty Nimi")},
		},
		{
			name: "skips and reports a row with an unknown attendance answer",
			rows: report(finnishHeader,
				[]string{"9001", "Testi Pelaaja", "Mies", "Ehkä"},
				[]string{"9002", "Keksitty Nimi", "Nainen", "Osallistuu"},
			),
			want:        []AttendancePlayerRow{attending(9002, "Keksitty Nimi")},
			wantSkipped: []SkippedRow{{Row: 5, MyClubID: "9001", Name: "Testi Pelaaja", Reason: SkipUnknownAttendance}},
		},
		{
			name: "skips and reports an attending player without a MyClub ID",
			rows: report(finnishHeader,
				[]string{"9001", "Testi Pelaaja", "Mies", "Osallistuu"},
				[]string{"", "Keksitty Nimi", "Nainen", "Osallistuu"},
			),
			want:        []AttendancePlayerRow{attending(9001, "Testi Pelaaja")},
			wantSkipped: []SkippedRow{{Row: 6, MyClubID: "", Name: "Keksitty Nimi", Reason: SkipMissingMemberNumber}},
		},
		{
			name: "skips and reports an attending player whose MyClub ID is not a number",
			rows: report(finnishHeader,
				[]string{"9001", "Testi Pelaaja", "Mies", "Osallistuu"},
				[]string{"90O2", "Keksitty Nimi", "Nainen", "Osallistuu"},
			),
			want:        []AttendancePlayerRow{attending(9001, "Testi Pelaaja")},
			wantSkipped: []SkippedRow{{Row: 6, MyClubID: "90O2", Name: "Keksitty Nimi", Reason: SkipInvalidMemberNumber}},
		},
		{
			name: "leaves out a player who isn't attending and has no MyClub ID, without reporting it",
			rows: report(finnishHeader,
				[]string{"", "Testi Pelaaja", "Mies", "Ei osallistu"},
				[]string{"", "Keksitty Nimi", "Nainen", "Ei vastausta"},
				[]string{"9003", "Kuvitteellinen Henkilö", "Mies", "Osallistuu"},
			),
			want: []AttendancePlayerRow{attending(9003, "Kuvitteellinen Henkilö")},
		},
		{
			name: "keeps an attending player whose name is empty, since players are matched by MyClub ID",
			rows: report(finnishHeader, []string{"9001", "", "Mies", "Osallistuu"}),
			want: []AttendancePlayerRow{attending(9001, "")},
		},
		{
			name: "reports no players when every row is skipped",
			rows: report(finnishHeader,
				[]string{"", "Testi Pelaaja", "Mies", "Osallistuu"},
				[]string{"9002", "Keksitty Nimi", "Nainen", "Ehkä"},
			),
			wantErr:     ErrNoPlayers,
			wantSkipped: []SkippedRow{
				{Row: 5, MyClubID: "", Name: "Testi Pelaaja", Reason: SkipMissingMemberNumber},
				{Row: 6, MyClubID: "9002", Name: "Keksitty Nimi", Reason: SkipUnknownAttendance},
			},
		},
		{
			name: "reports no players when the only readable rows aren't attending, and keeps the skipped rows",
			rows: report(finnishHeader,
				[]string{"", "Testi Pelaaja", "Mies", "Osallistuu"},
				[]string{"9002", "Keksitty Nimi", "Nainen", "Ei osallistu"},
				[]string{"9003", "Kuvitteellinen Henkilö", "Mies"},
			),
			wantErr:     ErrNoPlayers,
			wantSkipped: []SkippedRow{{Row: 5, MyClubID: "", Name: "Testi Pelaaja", Reason: SkipMissingMemberNumber}},
		},
		{
			name:    "reports no players when nobody is attending",
			rows:    report(finnishHeader, []string{"9001", "Testi Pelaaja", "Mies", "Ei osallistu"}, []string{"9002", "Keksitty Nimi", "Nainen"}),
			wantErr: ErrNoPlayers,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, skipped, err := parseAttendanceSheet(tc.rows)

			switch {
			case tc.wantMissing != nil:
				var missingErr *MissingColumnsError
				if !errors.As(err, &missingErr) {
					t.Fatalf("want an error naming the missing columns %v, got %v", tc.wantMissing, err)
				}
				if !reflect.DeepEqual(missingErr.Columns, tc.wantMissing) {
					t.Errorf("missing columns: got %v want %v", missingErr.Columns, tc.wantMissing)
				}
			case tc.wantErr != nil:
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("error: got %v want %v", err, tc.wantErr)
				}
				if len(skipped) != 0 || len(tc.wantSkipped) != 0 {
					if !reflect.DeepEqual(skipped, tc.wantSkipped) {
						t.Errorf("skipped rows:\ngot  %v\nwant %v", skipped, tc.wantSkipped)
					}
				}
			default:
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("rows:\ngot  %v\nwant %v", got, tc.want)
				}
				if len(skipped) != 0 || len(tc.wantSkipped) != 0 {
					if !reflect.DeepEqual(skipped, tc.wantSkipped) {
						t.Errorf("skipped rows:\ngot  %v\nwant %v", skipped, tc.wantSkipped)
					}
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
				actualDate, err := ParseDate(testCase.fileName)

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
			actualDate, err := ParseDate(testCase.fileName)

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

func TestPlayerAttendanceFiltering(t *testing.T) {
	type testCase struct {
		name                   string
		attendancePlayerRows   []AttendancePlayerRow
		attendanceStatus       AttendanceStatus
		expectedAttendanceRows []AttendancePlayerRow
		expectedErr            error
	}

	defaultAttendancePlayerRows := []AttendancePlayerRow{
		{PlayerRow: PlayerRow{MyclubID: 1111, Name: "Matti Meikäläinen"}, Attendance: AttendanceIn},
		{PlayerRow: PlayerRow{MyclubID: 2222, Name: "Kaija Kaarela"}, Attendance: AttendanceIn},
		{PlayerRow: PlayerRow{MyclubID: 3333, Name: "Teppo Teikäläinen"}, Attendance: AttendanceOut},
		{PlayerRow: PlayerRow{MyclubID: 4444, Name: "Saija Sappi"}, Attendance: AttendanceOut},
		{PlayerRow: PlayerRow{MyclubID: 5555, Name: "Aamos Kukko"}, Attendance: AttendanceUnknown},
		{PlayerRow: PlayerRow{MyclubID: 6666, Name: "Kati Paarila"}, Attendance: AttendanceUnknown},
	}

	testCases := []testCase{
		{
			name:                 "retrieves players marked as attending",
			attendancePlayerRows: defaultAttendancePlayerRows,
			attendanceStatus:     AttendanceIn,
			expectedAttendanceRows: []AttendancePlayerRow{
				{PlayerRow: PlayerRow{MyclubID: 1111, Name: "Matti Meikäläinen"}, Attendance: AttendanceIn},
				{PlayerRow: PlayerRow{MyclubID: 2222, Name: "Kaija Kaarela"}, Attendance: AttendanceIn},
			},
			expectedErr: nil,
		},
		{
			name:                 "retrieves players marked as not attending",
			attendancePlayerRows: defaultAttendancePlayerRows,
			attendanceStatus:     AttendanceOut,
			expectedAttendanceRows: []AttendancePlayerRow{
				{PlayerRow: PlayerRow{MyclubID: 3333, Name: "Teppo Teikäläinen"}, Attendance: AttendanceOut},
				{PlayerRow: PlayerRow{MyclubID: 4444, Name: "Saija Sappi"}, Attendance: AttendanceOut},
			},
			expectedErr: nil,
		},
		{
			name:                 "retrieves players who haven't responded",
			attendancePlayerRows: defaultAttendancePlayerRows,
			attendanceStatus:     AttendanceUnknown,
			expectedAttendanceRows: []AttendancePlayerRow{
				{PlayerRow: PlayerRow{MyclubID: 5555, Name: "Aamos Kukko"}, Attendance: AttendanceUnknown},
				{PlayerRow: PlayerRow{MyclubID: 6666, Name: "Kati Paarila"}, Attendance: AttendanceUnknown},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualAttendanceRows, err := GetAttendanceRowsByStatus(testCase.attendancePlayerRows, testCase.attendanceStatus)

			if testCase.expectedErr != nil {
				if err == nil {
					t.Error("expected error, got [nil]")
					return
				}

				if errors.Is(err, testCase.expectedErr) == false {
					t.Errorf("errors don't match: got [%s] want [%s]", err, testCase.expectedErr)
				}

			} else {
				if err != nil {
					t.Errorf("unexpected error: got [%v] want [nil]", err)
					return
				}
			}

			if len(actualAttendanceRows) != len(testCase.expectedAttendanceRows) {
				t.Errorf("number of returned attendance player rows don't match: got [%d] want [%d]", len(actualAttendanceRows), len(testCase.expectedAttendanceRows))
			}
			if reflect.DeepEqual(actualAttendanceRows, testCase.expectedAttendanceRows) == false {
				t.Errorf("returned attendance player rows don't match: got [%v] want [%v]", actualAttendanceRows, testCase.expectedAttendanceRows)
			}
		})
	}
}
