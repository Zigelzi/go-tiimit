package main

import (
	"net/http"
	"testing"
	"time"
)

// Scenario numbers refer to spec/milestones/00-test-harness.md.

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
