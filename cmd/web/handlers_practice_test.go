package main

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"
)

// Scenario numbers carry their milestone: 00 is
// spec/milestones/00-test-harness.md, 01 is
// spec/milestones/01-flexible-attendance-import.md.

// 00/S7: The front page shows a logged-out visitor no club data
func TestFrontPageShowsLoggedOutVisitorNoClubData(t *testing.T) {
	app := newTestApp(t)
	first := addPlayer(t, app, "Testi Pelaaja", 9001)
	addPlayer(t, app, "Keksitty Nimi", 9002)
	practiceID := addPractice(t, app, time.Date(2025, time.March, 12, 18, 0, 0, 0, time.UTC), first)
	client := newTestClient(t, app)

	rec := client.get("/")

	wantStatus(t, rec, http.StatusOK)
	wantBodyContains(t, rec, `href="/login"`)
	wantBodyLacks(t, rec,
		practiceLink(practiceID), // the practice
		"12 March",               // its date
		"active in the club",     // the player count
		`href="/practices/new"`,  // setting up a practice
		"Testi Pelaaja", "Keksitty Nimi",
	)
}

// reportFileName carries the practice date, which the import still reads
// from the file name.
const reportFileName = "testiseura-kuntofutis-2026-11-16.xlsx"

// reportRows returns the rows of an attendance report sheet: the title rows
// a real MyClub export has above the header, then the header, then players.
func reportRows(header []string, players ...[]string) [][]string {
	rows := [][]string{{"Testiseura Kuntofutis  (16.11.2026)"}, {"Kuntofutis"}, {}, header}
	return append(rows, players...)
}

// importReport builds an Excel file with one sheet holding rows and uploads
// it as the set-up page's form does.
func (c *testClient) importReport(sheetName string, rows [][]string) *httptest.ResponseRecorder {
	c.t.Helper()
	workbook := excelize.NewFile()
	defer workbook.Close()
	if err := workbook.SetSheetName("Sheet1", sheetName); err != nil {
		c.t.Fatalf("naming the report sheet failed: %v", err)
	}
	for i, row := range rows {
		cells := row
		if err := workbook.SetSheetRow(sheetName, fmt.Sprintf("A%d", i+1), &cells); err != nil {
			c.t.Fatalf("writing report row %d failed: %v", i+1, err)
		}
	}
	excelBytes, err := workbook.WriteToBuffer()
	if err != nil {
		c.t.Fatalf("writing the report failed: %v", err)
	}

	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	part, err := form.CreateFormFile("attendance-list", reportFileName)
	if err != nil {
		c.t.Fatalf("adding the report to the form failed: %v", err)
	}
	part.Write(excelBytes.Bytes())
	form.Close()

	req := httptest.NewRequest(http.MethodPost, "/practice", &body)
	req.Header.Set("Content-Type", form.FormDataContentType())
	req.Header.Set("HX-Request", "true")
	return c.do(req)
}

// wantPracticeCreated fails the test unless rec sends the browser on to a
// new practice, and returns that practice's path.
func wantPracticeCreated(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	wantStatus(t, rec, http.StatusOK)
	path := rec.Header().Get("HX-Redirect")
	if path == "" {
		t.Fatalf("the import should open the new practice, but it did not.\n\nThe page showed:\n%s", visibleText(rec.Body.String()))
	}
	return path
}

// wantNoPractice fails the test when any practice exists.
func wantNoPractice(t *testing.T, app *webConfig, rec *httptest.ResponseRecorder) {
	t.Helper()
	if path := rec.Header().Get("HX-Redirect"); path != "" {
		t.Errorf("the import should not open a practice, but it sent the browser to %s", path)
	}
	practices, err := app.queries.GetNewestPractices(context.Background(), 10)
	if err != nil {
		t.Fatalf("listing practices failed: %v", err)
	}
	if len(practices) != 0 {
		t.Errorf("no practice should exist, but %d do", len(practices))
	}
}

// 01/S1: Organizer imports a report with reordered and extra columns
func TestOrganizerImportsReportWithReorderedAndExtraColumns(t *testing.T) {
	app := newTestApp(t)
	names := []string{"Testi Pelaaja", "Keksitty Nimi", "Kuvitteellinen Henkilö", "Esimerkki Ihminen"}
	for i, name := range names {
		addPlayer(t, app, name, int64(9001+i))
	}
	client := loggedInCoach(t, app)

	rec := client.importReport("Tapahtuma", reportRows(
		[]string{"Osallistuu", "Jäsennumero", "Lempiväri", "Suku- ja etunimi"},
		[]string{"Osallistuu", "9001", "Sininen", "Testi Pelaaja"},
		[]string{"Osallistuu", "9002", "Punainen", "Keksitty Nimi"},
		[]string{"Osallistuu", "9003", "Vihreä", "Kuvitteellinen Henkilö"},
		[]string{"Osallistuu", "9004", "Keltainen", "Esimerkki Ihminen"},
	))

	practicePath := wantPracticeCreated(t, rec)
	wantBodyContains(t, bothTeams(t, client, practicePath), names...)
}

// bothTeams returns both teams' views of a practice in one reply, since the
// practice page shows one team at a time.
func bothTeams(t *testing.T, client *testClient, practicePath string) *httptest.ResponseRecorder {
	t.Helper()
	both := httptest.NewRecorder()
	for _, team := range []string{"1", "2"} {
		rec := client.get(practicePath + "/teams/" + team)
		wantStatus(t, rec, http.StatusOK)
		both.Body.Write(rec.Body.Bytes())
	}
	return both
}

// 01/S2: Organizer imports a report with a missing column
func TestOrganizerImportsReportWithMissingColumn(t *testing.T) {
	app := newTestApp(t)
	addPlayer(t, app, "Testi Pelaaja", 9001)
	client := loggedInCoach(t, app)

	rec := client.importReport("Tapahtuma", reportRows(
		[]string{"Suku- ja etunimi", "Sukupuoli", "Osallistuu"},
		[]string{"Testi Pelaaja", "Mies", "Osallistuu"},
	))

	wantNoPractice(t, app, rec)
	wantBodyContains(t, rec,
		"missing the column Jäsennumero",
		"Jäsennumero, Suku- ja etunimi and Osallistuu",
	)
}

// 01/S3: Organizer imports a report exported in English
func TestOrganizerImportsReportExportedInEnglish(t *testing.T) {
	app := newTestApp(t)
	addPlayer(t, app, "Testi Pelaaja", 9001)
	client := loggedInCoach(t, app)

	rec := client.importReport("Event", reportRows(
		[]string{"Member number", "Last and first name", "Gender", "Coming"},
		[]string{"9001", "Testi Pelaaja", "Male", "Coming"},
	))

	wantNoPractice(t, app, rec)
	wantBodyContains(t, rec,
		"in Finnish",
		"Jäsennumero, Suku- ja etunimi and Osallistuu",
	)
}

// 01/S4: Organizer imports a report with an unreadable attending player
func TestOrganizerImportsReportWithUnreadableAttendingPlayer(t *testing.T) {
	app := newTestApp(t)
	names := []string{"Testi Pelaaja", "Keksitty Nimi", "Kuvitteellinen Henkilö", "Esimerkki Ihminen"}
	for i, name := range names {
		addPlayer(t, app, name, int64(9001+i))
	}
	client := loggedInCoach(t, app)

	rec := client.importReport("Tapahtuma", reportRows(
		[]string{"Jäsennumero", "Suku- ja etunimi", "Sukupuoli", "Osallistuu"},
		[]string{"9001", "Testi Pelaaja", "Mies", "Osallistuu"},
		[]string{"9002", "Keksitty Nimi", "Nainen", "Osallistuu"},
		[]string{"9003", "Kuvitteellinen Henkilö", "Mies", "Osallistuu"},
		[]string{"", "Esimerkki Ihminen", "Mies", "Osallistuu"},
	))

	wantStatus(t, rec, http.StatusOK)
	if path := rec.Header().Get("HX-Redirect"); path != "" {
		t.Fatalf("the organizer should stay on the import result, but was sent to %s", path)
	}
	practices, err := app.queries.GetNewestPractices(context.Background(), 10)
	if err != nil || len(practices) != 1 {
		t.Fatalf("one practice should exist, got %d (error: %v)", len(practices), err)
	}
	practicePath := fmt.Sprintf("/practices/%d", practices[0].ID)

	wantBodyContains(t, rec,
		"Esimerkki Ihminen",
		"No MyClub ID",
		fmt.Sprintf(`href="%s"`, practicePath),
	)
	teams := bothTeams(t, client, practicePath)
	wantBodyContains(t, teams, names[:3]...)
	wantBodyLacks(t, teams, names[3])
}
