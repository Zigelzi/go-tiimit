package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Zigelzi/go-tiimit/internal/auth"
	"github.com/Zigelzi/go-tiimit/internal/db"
)

// Scenario numbers refer to spec/milestones/02-user-registration.md.

// invitePath matches the path of an invite link in a reply. The token is
// the hex string auth.MakeToken returns.
var invitePath = regexp.MustCompile(`/invites/([0-9a-f]+)`)

// 02/S1: Organizer views existing users
func TestOrganizerViewsExistingUsers(t *testing.T) {
	app := newTestApp(t)
	client := loggedInCoach(t, app)
	createUser(t, app, "ainvir", "made-up-password-2")

	rec := client.get("/")

	wantStatus(t, rec, http.StatusOK)
	wantBodyContains(t, rec, `href="/users"`, "Users", "2 users")

	rec = client.get("/users")

	wantStatus(t, rec, http.StatusOK)
	wantBodyContains(t, rec, "ainvir", "testcoach", `href="/"`, "Home", `hx-post="/users/invites"`)
	// The coach was created first, so a list in creation order would put
	// testcoach first. Scripts are left out because PostHog's set-up carries
	// the logged-in username.
	text := visibleText(rec.Body.String())
	if a, b := strings.Index(text, "ainvir"), strings.Index(text, "testcoach"); a == -1 || b == -1 || a > b {
		t.Errorf("the Users page should list ainvir before testcoach (alphabetical order), but does not.\n\nThe page showed:\n%s", text)
	}
}

// 02/S2: Organizer creates an invite link
func TestOrganizerCreatesInviteLink(t *testing.T) {
	app := newTestApp(t)
	// 16:40 UTC is 18:40 in Finland in November (UTC+2).
	app.now = func() time.Time { return time.Date(2026, 11, 21, 16, 40, 0, 0, time.UTC) }
	client := loggedInCoach(t, app)

	rec := client.htmxPost("/users/invites", url.Values{})

	wantStatus(t, rec, http.StatusOK)
	wantBodyContains(t, rec, "Share", "readonly", "Works once, until Mon 23 Nov 18:40")
	match := invitePath.FindStringSubmatch(rec.Body.String())
	if match == nil {
		t.Fatalf("the reply should carry an invite path /invites/<token>, but does not.\n\nThe page showed:\n%s", visibleText(rec.Body.String()))
	}
	if token := match[1]; len(token) < 64 {
		t.Errorf("the invite token %q has %d characters, want at least 64 (32 random bytes in hex, as auth.MakeToken makes)", token, len(token))
	}
	if got := countInvites(t, app); got != 1 {
		t.Errorf("creating an invite stored %d invites, want 1", got)
	}
	// The browser builds the full link from its own address (02/R1), so the
	// server must not guess one.
	wantBodyLacks(t, rec, "127.0.0.1", "localhost")
}

// 02/S2, 02/R11: creating an invite is a POST only, so following a link or
// a preview bot can never create one.
func TestCreatingInviteNeedsPost(t *testing.T) {
	app := newTestApp(t)
	client := loggedInCoach(t, app)

	rec := client.get("/users/invites")

	if rec.Code == http.StatusOK {
		t.Errorf("a GET to /users/invites answered 200 OK, want it refused")
	}
	if got := countInvites(t, app); got != 0 {
		t.Errorf("a GET to /users/invites stored %d invites, want 0", got)
	}
}

// 02/S9: Someone who is not logged in tries to see users or create an invite
func TestNotLoggedInCannotSeeUsersOrCreateInvite(t *testing.T) {
	cases := []struct {
		name string
		send func(c *testClient) *http.Response
	}{
		{"opens the Users page", func(c *testClient) *http.Response {
			return c.get("/users").Result()
		}},
		{"creates an invite from an old page", func(c *testClient) *http.Response {
			return c.do(formPost("/users/invites", false)).Result()
		}},
		{"creates an invite from an old page with htmx", func(c *testClient) *http.Response {
			return c.do(formPost("/users/invites", true)).Result()
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newTestApp(t)
			createCoach(t, app)
			client := newTestClient(t, app)

			res := tc.send(client)

			if got := redirectTarget(res); got != "/login" {
				t.Errorf("someone not logged in is sent to %q, want the login page /login", got)
			}
			if got := countInvites(t, app); got != 0 {
				t.Errorf("someone not logged in stored %d invites, want 0", got)
			}
		})
	}
}

// createUser creates an account directly, the way createCoach does.
func createUser(t *testing.T, app *webConfig, username, password string) {
	t.Helper()
	hashed, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("creating user %q failed: hashing the password: %v", username, err)
	}
	_, err = app.queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:       username,
		HashedPassword: hashed,
	})
	if err != nil {
		t.Fatalf("creating user %q failed: %v", username, err)
	}
}

// countInvites returns how many invites are stored.
func countInvites(t *testing.T, app *webConfig) int {
	t.Helper()
	var n int
	if err := app.db.QueryRow("SELECT COUNT(*) FROM user_invites").Scan(&n); err != nil {
		t.Fatalf("counting invites failed: %v", err)
	}
	return n
}

// formPost builds an empty form POST, with or without htmx's header.
func formPost(path string, htmx bool) *http.Request {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if htmx {
		req.Header.Set("HX-Request", "true")
	}
	return req
}

// redirectTarget returns where a reply sends the browser: htmx's
// HX-Redirect header, or a plain redirect's Location.
func redirectTarget(res *http.Response) string {
	if target := res.Header.Get("HX-Redirect"); target != "" {
		return target
	}
	if res.StatusCode >= 300 && res.StatusCode < 400 {
		return res.Header.Get("Location")
	}
	return ""
}
