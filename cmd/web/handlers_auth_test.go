package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Scenario numbers refer to spec/milestones/00-test-harness.md.

// 00/S1: A logged-out visitor is sent to the login page
func TestLoggedOutVisitorIsSentToLogin(t *testing.T) {
	app := newTestApp(t)
	client := newTestClient(t, app)

	rec := client.get("/players")

	wantStatus(t, rec, http.StatusSeeOther)
	if got := rec.Header().Get("Location"); got != "/login" {
		t.Errorf("the visitor is sent to %q, want /login", got)
	}
	wantBodyContains(t, client.get("/login"), "Log in")

	// The same request sent by htmx gets the redirect as a header instead.
	req := httptest.NewRequest(http.MethodGet, "/players", nil)
	req.Header.Set("HX-Request", "true")
	rec = client.do(req)

	wantStatus(t, rec, http.StatusUnauthorized)
	if got := rec.Header().Get("HX-Redirect"); got != "/login" {
		t.Errorf("the htmx reply sends the visitor to %q, want /login", got)
	}
}

// 00/S2: A coach logs in and sees the front page
func TestCoachLogsInAndSeesFrontPage(t *testing.T) {
	app := newTestApp(t)
	createCoach(t, app)
	first := addPlayer(t, app, "Testi Pelaaja", 9001)
	addPlayer(t, app, "Keksitty Nimi", 9002)
	practiceID := addPractice(t, app, time.Date(2025, time.March, 12, 18, 0, 0, 0, time.UTC), first)
	client := newTestClient(t, app)

	rec := client.login(coachUsername, coachPassword)

	if client.session == nil {
		t.Fatal("logging in with the right password set no session cookie")
	}
	if got := rec.Header().Get("HX-Redirect"); got != "/" {
		t.Fatalf("after logging in the coach is sent to %q, want the front page /", got)
	}

	rec = client.get("/")

	wantStatus(t, rec, http.StatusOK)
	// The practice is checked by its link: the page header shows today's
	// date, which on 12 March would match the practice's date as well.
	wantBodyContains(t, rec, practiceLink(practiceID), "2 active in the club")
}

// 00/S3: Wrong credentials show an error
func TestWrongCredentialsShowAnError(t *testing.T) {
	app := newTestApp(t)
	createCoach(t, app)
	client := newTestClient(t, app)

	rec := client.login(coachUsername, "not-the-password")

	wantStatus(t, rec, http.StatusOK)
	if cookie := sessionCookie(rec); cookie != nil {
		t.Errorf("a wrong password logged the visitor in: the reply set a session cookie")
	}
	if got := rec.Header().Get("HX-Redirect"); got != "" {
		t.Errorf("a wrong password sent the visitor to %q, want them to stay on the login page", got)
	}
	wantBodyContains(t, rec, "Incorrect username or password")
}
