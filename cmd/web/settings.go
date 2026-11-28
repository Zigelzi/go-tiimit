package main

import (
	"fmt"
	"runtime/debug"
)

// The settings that say where the server listens and which environment it runs in.
const (
	listenAddressSetting = "LISTEN_ADDRESS"
	environmentSetting   = "ENV"
)

// The defaults keep the server on this computer only, in development mode, so
// running it without settings never exposes it or caches stale files.
const (
	defaultListenAddress = "127.0.0.1:8080"
	defaultEnvironment   = "development"
)

// serverSettings is where the server listens and which environment it runs in.
type serverSettings struct {
	address string
	env     string
}

// readServerSettings returns the listen address and the environment from
// getenv, using the defaults for a setting that is empty. It returns an error
// naming the setting when the environment is neither development nor
// production, so a typo stops the app instead of running it in the wrong mode.
func readServerSettings(getenv func(string) string) (serverSettings, error) {
	settings := serverSettings{
		address: getenv(listenAddressSetting),
		env:     getenv(environmentSetting),
	}
	if settings.address == "" {
		settings.address = defaultListenAddress
	}
	if settings.env == "" {
		settings.env = defaultEnvironment
	}
	if settings.env != "development" && settings.env != "production" {
		return serverSettings{}, fmt.Errorf("%s is %q, want development or production", environmentSetting, settings.env)
	}
	return settings, nil
}

// buildVersion returns the commit the binary was built from, as Go records it
// in a binary built inside a git checkout, or "unknown" when there is none (go
// run and go test record no commit). Pass it the results of debug.ReadBuildInfo.
func buildVersion(info *debug.BuildInfo, ok bool) string {
	if !ok || info == nil {
		return "unknown"
	}
	revision, modified := "", false
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}
	if revision == "" {
		return "unknown"
	}
	if modified {
		return revision + " with uncommitted changes"
	}
	return revision
}
