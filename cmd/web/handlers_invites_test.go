package main

import (
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// Scenario numbers refer to spec/milestones/02-user-registration.md.
//
// The sign-up form swaps only its error area, as the login form does, so
// what the person typed and the password hint stay on the page by
// themselves. These tests check the reply; they do not re-check the form.

const (
	newUserPassword = "Kesä Treeni 2026"
	passwordHint    = "At least 10 characters. A few words separated with dashes (-) works well."
	usernameHint    = "5 to 20 letters. No numbers or spaces."
)

// 02/S3: New organizer joins with an invite link
func TestNewOrganizerJoinsWithInviteLink(t *testing.T) {
	app := newTestApp(t)
	path := createInvite(t, app)

	// WhatsApp's preview opens the link without sending the form.
	preview := newTestClient(t, app).get(path)
	wantStatus(t, preview, http.StatusOK)

	joiner := newTestClient(t, app)
	rec := joiner.get(path)

	wantStatus(t, rec, http.StatusOK)
	wantBodyContains(t, rec,
		`name="username"`, `name="password"`,
		`autocomplete="username"`, `autocomplete="new-password"`, `autocapitalize="none"`,
		usernameHint, passwordHint)
	// 02/R3: the token in the address never reaches PostHog.
	wantBodyLacks(t, rec, "analytics.js", "initPostHog")

	rec = joiner.signUp(path, " AinVir ", newUserPassword)

	wantStatus(t, rec, http.StatusOK)
	if got := rec.Header().Get("HX-Redirect"); got != "/" {
		t.Fatalf("after signing up the new organizer is sent to %q, want the front page /\n\nThe reply was:\n%s", got, visibleText(rec.Body.String()))
	}
	if joiner.session == nil {
		t.Fatal("signing up set no session cookie, so the new organizer is not logged in")
	}
	rec = joiner.get("/")
	wantBodyContains(t, rec, `href="/users"`)
	if got := inviteUsedBy(t, app, path); got != "ainvir" {
		t.Errorf("the invite records it created the account %q, want ainvir", got)
	}

	later := newTestClient(t, app)
	later.login("ainvir", newUserPassword)
	if later.session == nil {
		t.Errorf("logging in later as ainvir with %q failed", newUserPassword)
	}
}

// 02/S6: New organizer chooses a username that is taken
func TestNewOrganizerChoosesTakenUsername(t *testing.T) {
	app := newTestApp(t)
	createUser(t, app, "ainvir", "made-up-password-2")
	path := createInvite(t, app)
	joiner := newTestClient(t, app)

	rec := joiner.signUp(path, " AINVIR ", newUserPassword)

	wantStatus(t, rec, http.StatusOK)
	wantBodyContains(t, rec, "taken")
	wantNoSignUp(t, joiner, rec)
	if got := countUsers(t, app); got != 2 {
		t.Errorf("there are %d users, want 2 (the coach and the existing ainvir)", got)
	}
	wantLinkStillWorks(t, app, path)
}

// 02/S7: New organizer sends the form with a field that breaks the rules
func TestNewOrganizerSendsFieldThatBreaksRules(t *testing.T) {
	cases := []struct {
		name, username, password, message string
	}{
		{"empty username", "   ", newUserPassword, "Choose a username"},
		{"username shorter than 5 letters", "aino", newUserPassword, "5 to 20 letters"},
		{"username longer than 20 letters", strings.Repeat("a", 21), newUserPassword, "5 to 20 letters"},
		{"username with a dash", "ain-vir", newUserPassword, "only letters"},
		{"username with a digit", "ainvir2", newUserPassword, "only letters"},
		{"password shorter than 10 characters", "ainvir", "lyhyt-sa9", "at least 10 characters"},
		{"password longer than 64 characters", "ainvir", strings.Repeat("a", 65), "at most 64 characters"},
		{"password too long to store", "ainvir", strings.Repeat("ä", 37) + strings.Repeat("a", 27), "too long"},
		{"password same as the username", "ainovirtanen", "AinoVirtanen", "same as your username"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newTestApp(t)
			path := createInvite(t, app)
			joiner := newTestClient(t, app)

			rec := joiner.signUp(path, tc.username, tc.password)

			wantStatus(t, rec, http.StatusOK)
			wantBodyContains(t, rec, tc.message)
			wantNoSignUp(t, joiner, rec)
			if got := countUsers(t, app); got != 1 {
				t.Errorf("there are %d users, want only the coach", got)
			}
			wantLinkStillWorks(t, app, path)
		})
	}
}

// 02/S8: Organizer logs in with different capitals or spaces in the username
func TestOrganizerLogsInWithDifferentCapitals(t *testing.T) {
	app := newTestApp(t)
	createUser(t, app, "ainvir", newUserPassword)
	client := newTestClient(t, app)

	client.login(" AinVir ", newUserPassword)

	if client.session == nil {
		t.Errorf("logging in as %q failed, want it to match the user ainvir", " AinVir ")
	}
}

// 02/S12: Two people send the form from the same link
func TestTwoPeopleSendFormFromSameLink(t *testing.T) {
	app := newTestApp(t)
	path := createInvite(t, app)
	first := newTestClient(t, app)
	second := newTestClient(t, app)
	wantStatus(t, first.get(path), http.StatusOK)
	wantStatus(t, second.get(path), http.StatusOK)

	rec := first.signUp(path, "ainvir", newUserPassword)
	if first.session == nil {
		t.Fatalf("the first person was not logged in\n\nThe reply was:\n%s", visibleText(rec.Body.String()))
	}

	rec = second.signUp(path, "ainovirtanen", "toinen-salasana-2026")

	wantBodyContains(t, rec, "already used")
	wantNoSignUp(t, second, rec)
	if got := countUsers(t, app); got != 2 {
		t.Errorf("there are %d users, want 2 (the coach and the first person)", got)
	}
}

// createInvite creates the coach and an invite through the Users page, and
// returns the invite's path.
func createInvite(t *testing.T, app *webConfig) string {
	t.Helper()
	coach := loggedInCoach(t, app)
	rec := coach.htmxPost("/users/invites", url.Values{})
	wantStatus(t, rec, http.StatusOK)
	match := invitePath.FindString(rec.Body.String())
	if match == "" {
		t.Fatalf("creating an invite gave no invite path\n\nThe reply was:\n%s", visibleText(rec.Body.String()))
	}
	return match
}

// signUp sends the sign-up form on an invite page and keeps the session
// cookie if one is set.
func (c *testClient) signUp(path, username, password string) *httptest.ResponseRecorder {
	c.t.Helper()
	rec := c.htmxPost(path, url.Values{"username": {username}, "password": {password}})
	if cookie := sessionCookie(rec); cookie != nil {
		c.session = cookie
	}
	return rec
}

// wantNoSignUp fails the test when a sign-up reply logged someone in or
// moved them away from the form.
func wantNoSignUp(t *testing.T, c *testClient, rec *httptest.ResponseRecorder) {
	t.Helper()
	if c.session != nil {
		t.Errorf("the reply logged the person in, want it to keep them on the form")
	}
	if got := rec.Header().Get("HX-Redirect"); got != "" {
		t.Errorf("the reply sends the person to %q, want it to keep them on the form", got)
	}
}

// wantLinkStillWorks fails the test when a fresh person cannot sign up with
// the invite at path.
func wantLinkStillWorks(t *testing.T, app *webConfig, path string) {
	t.Helper()
	other := newTestClient(t, app)
	rec := other.signUp(path, "toinenuusi", "toinen-salasana-2026")
	if other.session == nil {
		t.Errorf("the link no longer works for another try\n\nThe reply was:\n%s", visibleText(rec.Body.String()))
	}
}

// countUsers returns how many accounts exist.
func countUsers(t *testing.T, app *webConfig) int {
	t.Helper()
	var n int
	if err := app.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&n); err != nil {
		t.Fatalf("counting users failed: %v", err)
	}
	return n
}

// inviteUsedBy returns the username of the account the invite at path
// created, or "" if none.
func inviteUsedBy(t *testing.T, app *webConfig, path string) string {
	t.Helper()
	token := strings.TrimPrefix(path, "/invites/")
	var username string
	err := app.db.QueryRow(`
		SELECT COALESCE(u.username, '')
		FROM user_invites i LEFT JOIN users u ON u.id = i.used_by_user_id
		WHERE i.token = ?`, token).Scan(&username)
	if err != nil {
		t.Fatalf("reading the invite failed: %v", err)
	}
	return username
}

// 02/S4: Someone opens an invite link that was already used
func TestSomeoneOpensUsedInviteLink(t *testing.T) {
	app := newTestApp(t)
	path := createInvite(t, app)
	joiner := newTestClient(t, app)
	joiner.signUp(path, "ainvir", newUserPassword)
	if joiner.session == nil {
		t.Fatal("signing up with the link failed, so the test cannot use it up")
	}

	rec := newTestClient(t, app).get(path)

	wantStatus(t, rec, http.StatusOK)
	wantBodyContains(t, rec, "Invitation already used", "This invitation was already used. Ask the organizer for a new one.")
	wantNoForm(t, rec)
}

// 02/S5: Someone opens an expired invite link
func TestSomeoneOpensExpiredInviteLink(t *testing.T) {
	created := time.Date(2026, 11, 21, 16, 40, 0, 0, time.UTC)
	cases := []struct {
		name    string
		openAt  time.Time
		expired bool
	}{
		{"one second before 2 days", created.Add(48*time.Hour - time.Second), false},
		{"exactly 2 days", created.Add(48 * time.Hour), true},
		{"more than 2 days", created.Add(72 * time.Hour), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newTestApp(t)
			app.now = func() time.Time { return created }
			path := createInvite(t, app)
			app.now = func() time.Time { return tc.openAt }

			rec := newTestClient(t, app).get(path)

			wantStatus(t, rec, http.StatusOK)
			if tc.expired {
				wantBodyContains(t, rec, "Invitation expired", "This invitation has expired. Ask the organizer for a new one.")
				wantNoForm(t, rec)
			} else {
				wantBodyContains(t, rec, `name="password"`)
				wantBodyLacks(t, rec, "Invitation expired")
			}
		})
	}
}

// 02/S11: Someone who is logged in opens an invite link
func TestLoggedInOrganizerOpensInviteLink(t *testing.T) {
	app := newTestApp(t)
	coach := loggedInCoach(t, app)
	rec := coach.htmxPost("/users/invites", url.Values{})
	path := invitePath.FindString(rec.Body.String())
	if path == "" {
		t.Fatalf("creating an invite gave no invite path\n\nThe reply was:\n%s", visibleText(rec.Body.String()))
	}

	rec = coach.get(path)

	wantStatus(t, rec, http.StatusOK)
	wantBodyContains(t, rec, "Already logged in", coachUsername, `hx-post="/logout"`)
	wantNoForm(t, rec)
	wantBodyLacks(t, rec, "analytics.js", "initPostHog")

	rec = newTestClient(t, app).get(path)
	wantBodyContains(t, rec, `name="password"`)
}

// 02/S13: Someone opens a mistyped or cut-off invite link
func TestSomeoneOpensMistypedInviteLink(t *testing.T) {
	cases := []struct {
		name string
		path func(real string) string
	}{
		{"end cut off", func(real string) string { return real[:len(real)-5] }},
		{"made-up code", func(string) string { return "/invites/abc123" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newTestApp(t)
			path := tc.path(createInvite(t, app))
			joiner := newTestClient(t, app)

			rec := joiner.get(path)

			wantStatus(t, rec, http.StatusOK)
			// templ escapes the apostrophe, so look for it as it appears in HTML.
			wantBodyContains(t, rec, html.EscapeString("Invitation doesn't work"), "This invitation does not work. Ask the organizer for a new one.")
			wantNoForm(t, rec)
			wantBodyLacks(t, rec, "analytics.js", "initPostHog")

			rec = joiner.signUp(path, "ainvir", newUserPassword)

			wantBodyContains(t, rec, "This invitation does not work. Ask the organizer for a new one.")
			wantNoSignUp(t, joiner, rec)
			if got := countUsers(t, app); got != 1 {
				t.Errorf("there are %d users, want only the coach", got)
			}
		})
	}
}

// wantNoForm fails the test when the page shows the form to create an
// account.
func wantNoForm(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	wantBodyLacks(t, rec, `name="username"`, `name="password"`)
}

// 02/S11: logging out works after signing up with a link. A cookie set
// without a path is scoped to the folder of the address that set it, so a
// session started at /invites/<token> would only reach /invites pages and
// never /logout.
func TestSessionCookieCoversWholeApp(t *testing.T) {
	app := newTestApp(t)
	path := createInvite(t, app)

	signUp := newTestClient(t, app).htmxPost(path, url.Values{"username": {"ainvir"}, "password": {newUserPassword}})
	login := newTestClient(t, app).htmxPost("/login", url.Values{"username": {coachUsername}, "password": {coachPassword}})

	for name, rec := range map[string]*httptest.ResponseRecorder{"signing up": signUp, "logging in": login} {
		cookie := sessionCookie(rec)
		if cookie == nil {
			t.Fatalf("%s set no session cookie", name)
		}
		if cookie.Path != "/" {
			t.Errorf("%s sets the session cookie for path %q, want / so every page and /logout get it", name, cookie.Path)
		}
	}
}

// 02/S11: logging out always takes you to the front page, even when the
// request carries no session, so the button never seems to do nothing.
func TestLogoutWithoutSessionGoesToFrontPage(t *testing.T) {
	app := newTestApp(t)

	rec := newTestClient(t, app).htmxPost("/logout", url.Values{})

	if got := rec.Header().Get("HX-Redirect"); got != "/" {
		t.Errorf("logging out without a session sends the browser to %q, want the front page /", got)
	}
}

// Code review: the sign-up and login forms post even before htmx has
// loaded, so a password never ends up in the address; and the invite page
// tells the browser not to send its address (with the token) to other sites.
func TestFormsPostWithoutHtmxAndInvitePageHidesReferrer(t *testing.T) {
	app := newTestApp(t)
	path := createInvite(t, app)

	invite := newTestClient(t, app).get(path)
	wantBodyContains(t, invite, `method="post"`, `action="`+path+`"`, `<meta name="referrer" content="same-origin"`)

	login := newTestClient(t, app).get("/login")
	wantBodyContains(t, login, `method="post"`, `action="/login"`)
}

// 02/S11, code review: someone logged in who sends the sign-up form (for
// example from a tab opened before logging in) gets no new account.
func TestLoggedInOrganizerSendsSignUpForm(t *testing.T) {
	app := newTestApp(t)
	coach := loggedInCoach(t, app)
	rec := coach.htmxPost("/users/invites", url.Values{})
	path := invitePath.FindString(rec.Body.String())
	before := coach.session

	rec = coach.signUp(path, "ainvir", newUserPassword)

	wantBodyContains(t, rec, "already logged in")
	if got := rec.Header().Get("HX-Redirect"); got != "" {
		t.Errorf("the reply sends the browser to %q, want it to stay on the page", got)
	}
	if coach.session != before {
		t.Errorf("sending the form switched the session, want the organizer to stay logged in as before")
	}
	if got := countUsers(t, app); got != 1 {
		t.Errorf("there are %d users, want only the coach", got)
	}
	newTestClient(t, app).signUp(path, "toinenuusi", newUserPassword)
	if got := countUsers(t, app); got != 2 {
		t.Errorf("the link no longer works for someone who is not logged in")
	}
}
