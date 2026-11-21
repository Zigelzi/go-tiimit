package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Zigelzi/go-tiimit/internal/auth"
)

// wantLoggedIn fails the test unless logging in with username and password
// sets a session.
func wantLoggedIn(t *testing.T, app *webConfig, username, password string) {
	t.Helper()
	client := newTestClient(t, app)
	rec := client.login(username, password)
	if client.session == nil {
		t.Fatalf("logging in as %q failed: status %d, no session cookie", username, rec.Code)
	}
	wantStatus(t, client.get("/users"), http.StatusOK)
}

// 03/S1: Developer starts Tiimit on an empty database with first-user settings
func TestStartupCreatesFirstUser(t *testing.T) {
	app := newTestApp(t)

	// Capitals and surrounding spaces: the username is normalized as at
	// sign-up (03/R5), the password is used exactly as set.
	outcome, err := app.ensureFirstUser(context.Background(), "  FirstCoach ", "made-up-first-password")

	if err != nil {
		t.Fatalf("ensureFirstUser returned %v, want no error", err)
	}
	if outcome != firstUserCreated {
		t.Fatalf("outcome = %v, want firstUserCreated", outcome)
	}
	wantLoggedIn(t, app, "firstcoach", "made-up-first-password")
}

// 03/S2: Developer starts Tiimit on an empty database with first-user settings that break the rules
func TestStartupRefusesFirstUserThatBreaksRules(t *testing.T) {
	cases := []struct {
		name        string
		username    string
		password    string
		wantSetting string
		wantReason  error
	}{
		{"username too short", "abc", "made-up-first-password", "FIRST_USER_USERNAME", auth.ErrUsernameLength},
		{"password too short", "firstcoach", "short-pw", "FIRST_USER_PASSWORD", auth.ErrPasswordShort},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newTestApp(t)

			_, err := app.ensureFirstUser(context.Background(), tc.username, tc.password)

			if err == nil {
				t.Fatal("ensureFirstUser returned no error, want one that stops startup")
			}
			if !strings.Contains(err.Error(), tc.wantSetting) {
				t.Errorf("error %q does not name %s", err, tc.wantSetting)
			}
			if !errors.Is(err, tc.wantReason) {
				t.Errorf("error %q does not say why: want %q", err, tc.wantReason)
			}
			if strings.Contains(err.Error(), tc.password) {
				t.Errorf("error %q shows the password", err)
			}
			if n := countUsers(t, app); n != 0 {
				t.Errorf("users in database = %d, want 0", n)
			}
		})
	}
}

// 03/S3: Developer starts Tiimit on an empty database without first-user settings
func TestStartupWithoutFirstUserSettings(t *testing.T) {
	app := newTestApp(t)

	outcome, err := app.ensureFirstUser(context.Background(), "", "")

	if err != nil {
		t.Fatalf("ensureFirstUser returned %v, want no error", err)
	}
	if outcome != firstUserNotSet {
		t.Fatalf("outcome = %v, want firstUserNotSet", outcome)
	}
	if n := countUsers(t, app); n != 0 {
		t.Errorf("users in database = %d, want 0", n)
	}
}

// 03/S4: Developer starts Tiimit on a database that already has users
func TestStartupIgnoresFirstUserSettingsWhenUsersExist(t *testing.T) {
	app := newTestApp(t)
	createUser(t, app, "ainvir", "made-up-password-2")

	outcome, err := app.ensureFirstUser(context.Background(), "firstcoach", "short")

	if err != nil {
		t.Fatalf("ensureFirstUser returned %v, want no error", err)
	}
	if outcome != firstUserNotNeeded {
		t.Fatalf("outcome = %v, want firstUserNotNeeded", outcome)
	}
	if n := countUsers(t, app); n != 1 {
		t.Errorf("users in database = %d, want 1", n)
	}
	wantLoggedIn(t, app, "ainvir", "made-up-password-2")
}

// 03/S5: Developer starts Tiimit on an empty database with only one first-user setting
func TestStartupRefusesHalfSetFirstUserSettings(t *testing.T) {
	cases := []struct {
		name        string
		username    string
		password    string
		wantSetting string
	}{
		{"password missing", "firstcoach", "", "FIRST_USER_PASSWORD"},
		{"username missing", "", "made-up-first-password", "FIRST_USER_USERNAME"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newTestApp(t)

			_, err := app.ensureFirstUser(context.Background(), tc.username, tc.password)

			if err == nil {
				t.Fatal("ensureFirstUser returned no error, want one that stops startup")
			}
			if !strings.Contains(err.Error(), tc.wantSetting) || !strings.Contains(err.Error(), "missing") {
				t.Errorf("error %q does not say %s is missing", err, tc.wantSetting)
			}
			if tc.password != "" && strings.Contains(err.Error(), tc.password) {
				t.Errorf("error %q shows the password", err)
			}
			if n := countUsers(t, app); n != 0 {
				t.Errorf("users in database = %d, want 0", n)
			}
		})
	}
}
