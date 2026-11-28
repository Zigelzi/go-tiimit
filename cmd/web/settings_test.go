package main

import (
	"runtime/debug"
	"strings"
	"testing"
)

// fakeSettings returns a getenv that reads from the given map, as os.Getenv
// reads from the environment.
func fakeSettings(values map[string]string) func(string) string {
	return func(name string) string { return values[name] }
}

// staticFileCacheControl returns the Cache-Control header the app sends with a
// static file when it runs in the given environment.
func staticFileCacheControl(t *testing.T, env string) string {
	t.Helper()
	app := newTestApp(t)
	app.env = env
	rec := newTestClient(t, app).get("/static/tabs.js")
	wantStatus(t, rec, 200)
	return rec.Header().Get("Cache-Control")
}

// 04/S7: The app runs with the settings it is given
func TestServerSettingsAreReadFromTheEnvironment(t *testing.T) {
	settings, err := readServerSettings(fakeSettings(map[string]string{
		"LISTEN_ADDRESS": "127.0.0.1:9090",
		"ENV":            "production",
	}))
	if err != nil {
		t.Fatalf("reading the settings failed: %v", err)
	}

	if settings.address != "127.0.0.1:9090" {
		t.Errorf("address = %q, want %q", settings.address, "127.0.0.1:9090")
	}
	if settings.env != "production" {
		t.Errorf("env = %q, want %q", settings.env, "production")
	}
	// Production mode lets browsers cache the static files.
	if got := staticFileCacheControl(t, settings.env); got == "no-store" {
		t.Errorf("Cache-Control in production = %q, want the files cacheable", got)
	}
}

// 04/S11: The app runs as it does today without settings
func TestServerSettingsDefaultToLocalDevelopment(t *testing.T) {
	settings, err := readServerSettings(fakeSettings(nil))
	if err != nil {
		t.Fatalf("reading the settings failed: %v", err)
	}

	if settings.address != "127.0.0.1:8080" {
		t.Errorf("address = %q, want %q", settings.address, "127.0.0.1:8080")
	}
	if settings.env != "development" {
		t.Errorf("env = %q, want %q", settings.env, "development")
	}
	if got := staticFileCacheControl(t, settings.env); got != "no-store" {
		t.Errorf("Cache-Control in development = %q, want %q", got, "no-store")
	}
}

// A mistyped environment stops the app instead of quietly running it in
// development mode on the Pi.
func TestServerSettingsRejectAnUnknownEnvironment(t *testing.T) {
	_, err := readServerSettings(fakeSettings(map[string]string{"ENV": "prod"}))

	if err == nil {
		t.Fatal("got no error for ENV=prod, want one")
	}
	if !strings.Contains(err.Error(), "ENV") {
		t.Errorf("error %q does not name the ENV setting", err)
	}
}

// 04/S1 (the log line): the version names the commit the binary was built from.
func TestBuildVersion(t *testing.T) {
	const commit = "0123456789abcdef0123456789abcdef01234567"
	tests := []struct {
		name string
		info *debug.BuildInfo
		ok   bool
		want string
	}{
		{"no build info", nil, false, "unknown"},
		{"built without git, as go run and go test are", &debug.BuildInfo{}, true, "unknown"},
		{
			"built from a commit",
			&debug.BuildInfo{Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: commit},
				{Key: "vcs.modified", Value: "false"},
			}},
			true, commit,
		},
		{
			"built with uncommitted changes",
			&debug.BuildInfo{Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: commit},
				{Key: "vcs.modified", Value: "true"},
			}},
			true, commit + " with uncommitted changes",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildVersion(tt.info, tt.ok); got != tt.want {
				t.Errorf("buildVersion = %q, want %q", got, tt.want)
			}
		})
	}
}
