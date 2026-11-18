package main

import (
	"net/http"
	"net/url"
	"testing"
)

// Scenario numbers refer to spec/milestones/00-test-harness.md.

// 00/S4: A coach adds a player and sees them in the list
func TestCoachAddsPlayerAndSeesThemInList(t *testing.T) {
	app := newTestApp(t)
	client := loggedInCoach(t, app)

	rec := client.htmxPost("/players", url.Values{"name": {"Uusi Keksitty"}, "myclub-id": {"4242"}})

	wantStatus(t, rec, http.StatusOK)
	if got := rec.Header().Get("HX-Redirect"); got != "/players" {
		t.Fatalf("after adding a player the coach is sent to %q, want the players list /players\n\nThe reply was:\n%s", got, rec.Body.String())
	}

	rec = client.get("/players")

	wantStatus(t, rec, http.StatusOK)
	wantBodyContains(t, rec, "Uusi Keksitty")
}
