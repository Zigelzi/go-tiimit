package main

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Zigelzi/go-tiimit/internal/auth"
	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/pressly/goose/v3"
)

// Helpers shared by the HTTP acceptance tests: building the app on a fresh
// database, sending requests as a browser would, and setting up data.
// Helpers that only one file's tests use stay in that file.
//
// Every name in test data is made up. Real player names are personal data
// and never appear in the repository.
//
// Tests do not call t.Parallel: Goose keeps its settings in package-level
// variables, so two tests migrating at once would race.

const (
	coachUsername = "testcoach"
	coachPassword = "made-up-password"
)

// newTestApp builds the app on a fresh, migrated SQLite file in a temp
// directory, so a test never touches tiimit.db. The PostHog key is left
// empty, which turns analytics off.
func newTestApp(t *testing.T) *webConfig {
	t.Helper()
	goose.SetLogger(goose.NopLogger())
	database, err := db.InitDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("creating the test database failed: %v", err)
	}
	t.Cleanup(func() { database.Close() })
	if err := db.RunMigrations(database); err != nil {
		t.Fatalf("migrating the test database failed: %v", err)
	}
	return &webConfig{
		queries: db.New(database),
		db:      database,
		env:     "test",
		now:     time.Now,
	}
}

// testClient sends requests to the app's router, carrying the session cookie
// once it has one, as a browser would.
type testClient struct {
	t       *testing.T
	handler http.Handler
	session *http.Cookie
}

func newTestClient(t *testing.T, app *webConfig) *testClient {
	t.Helper()
	return &testClient{t: t, handler: app.routes()}
}

// do sends req with the session cookie, if any, and returns the response
// without following redirects.
func (c *testClient) do(req *http.Request) *httptest.ResponseRecorder {
	c.t.Helper()
	if c.session != nil {
		req.AddCookie(c.session)
	}
	rec := httptest.NewRecorder()
	c.handler.ServeHTTP(rec, req)
	return rec
}

// get sends a plain browser GET.
func (c *testClient) get(path string) *httptest.ResponseRecorder {
	c.t.Helper()
	return c.do(httptest.NewRequest(http.MethodGet, path, nil))
}

// htmxPost submits a form the way the app's htmx forms do.
func (c *testClient) htmxPost(path string, form url.Values) *httptest.ResponseRecorder {
	c.t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	return c.do(req)
}

// login submits the login form and keeps the session cookie if one is set.
func (c *testClient) login(username, password string) *httptest.ResponseRecorder {
	c.t.Helper()
	rec := c.htmxPost("/login", url.Values{"username": {username}, "password": {password}})
	if cookie := sessionCookie(rec); cookie != nil {
		c.session = cookie
	}
	return rec
}

// sessionCookie returns the session cookie the response sets, or nil.
func sessionCookie(rec *httptest.ResponseRecorder) *http.Cookie {
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "user_session_id" && cookie.Value != "" {
			return cookie
		}
	}
	return nil
}

// loggedInCoach creates the coach account and returns a client logged in
// through the real login form.
func loggedInCoach(t *testing.T, app *webConfig) *testClient {
	t.Helper()
	createCoach(t, app)
	client := newTestClient(t, app)
	client.login(coachUsername, coachPassword)
	if client.session == nil {
		t.Fatal("logging in as the coach failed: the reply set no session cookie")
	}
	return client
}

// createCoach creates the coach account the way cmd/cli does.
func createCoach(t *testing.T, app *webConfig) {
	t.Helper()
	hashed, err := auth.HashPassword(coachPassword)
	if err != nil {
		t.Fatalf("creating the coach account failed: hashing the password: %v", err)
	}
	_, err = app.queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:       coachUsername,
		HashedPassword: hashed,
	})
	if err != nil {
		t.Fatalf("creating the coach account failed: %v", err)
	}
}

// addPlayer adds an active player with default ratings and returns its id.
func addPlayer(t *testing.T, app *webConfig, name string, myClubID int64) int64 {
	t.Helper()
	ctx := context.Background()
	err := app.queries.AddPlayer(ctx, db.AddPlayerParams{
		Name:         name,
		MyclubID:     myClubID,
		RunPower:     5,
		BallHandling: 5,
	})
	if err != nil {
		t.Fatalf("adding player %q failed: %v", name, err)
	}
	player, err := app.queries.GetPlayerByMyclubID(ctx, myClubID)
	if err != nil {
		t.Fatalf("finding the added player %q by MyClub ID %d failed: %v", name, myClubID, err)
	}
	return player.ID
}

// addPractice creates a practice on date with the players on team 1, and
// returns its id.
func addPractice(t *testing.T, app *webConfig, date time.Time, playerIDs ...int64) int64 {
	t.Helper()
	ctx := context.Background()
	practiceID, err := app.queries.CreatePractice(ctx, date)
	if err != nil {
		t.Fatalf("creating the practice on %s failed: %v", date.Format("2 January 2006"), err)
	}
	for _, playerID := range playerIDs {
		err := app.queries.AddPlayerToPractice(ctx, db.AddPlayerToPracticeParams{
			PracticeID: practiceID,
			PlayerID:   playerID,
			TeamNumber: 1,
		})
		if err != nil {
			t.Fatalf("adding player %d to the practice failed: %v", playerID, err)
		}
	}
	return practiceID
}

// practiceLink returns the link the front page shows for a practice.
func practiceLink(practiceID int64) string {
	return fmt.Sprintf(`href="/practices/%d"`, practiceID)
}

// wantStatus fails the test when the response's status is not want.
func wantStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("the reply's status is %d %s, want %d %s\n\nThe reply was:\n%s",
			rec.Code, http.StatusText(rec.Code), want, http.StatusText(want), rec.Body.String())
	}
}

// wantBodyContains fails the test when the response body lacks any of texts.
// It lists every missing text, then prints the body once.
func wantBodyContains(t *testing.T, rec *httptest.ResponseRecorder, texts ...string) {
	t.Helper()
	body := rec.Body.String()
	var missing []string
	for _, text := range texts {
		if !strings.Contains(body, text) {
			missing = append(missing, fmt.Sprintf("%q", text))
		}
	}
	if len(missing) > 0 {
		t.Errorf("the page should show %s, but does not.\n\nThe page showed:\n%s",
			strings.Join(missing, ", "), visibleText(body))
	}
}

// wantBodyLacks fails the test when the response body contains any of texts.
// It lists every text found, then prints the body once.
func wantBodyLacks(t *testing.T, rec *httptest.ResponseRecorder, texts ...string) {
	t.Helper()
	body := rec.Body.String()
	var found []string
	for _, text := range texts {
		if strings.Contains(body, text) {
			found = append(found, fmt.Sprintf("%q", text))
		}
	}
	if len(found) > 0 {
		t.Errorf("the page should not show %s, but does.\n\nThe page showed:\n%s",
			strings.Join(found, ", "), visibleText(body))
	}
}

var (
	hiddenParts = regexp.MustCompile(`(?is)<head.*?</head>|<script.*?</script>|<style.*?</style>`)
	tags        = regexp.MustCompile(`<[^>]*>`)
)

// visibleText returns the text a person sees on an HTML page, one element
// per line, so a failing test prints something readable instead of markup.
// It is a rough reading for error messages, not something to assert on.
func visibleText(page string) string {
	page = hiddenParts.ReplaceAllString(page, "")
	page = tags.ReplaceAllString(page, "\n")
	var lines []string
	for _, line := range strings.Split(html.UnescapeString(page), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, "    "+line)
		}
	}
	return strings.Join(lines, "\n")
}
