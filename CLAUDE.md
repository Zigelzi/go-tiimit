# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Tiimit distributes football (soccer) teams for practice sessions. A coach imports a club's players from a MyClub Excel export, then for each practice imports an attendance-report Excel export; the app splits confirmed players into two balanced teams based on player scores. It's a webapp (server-rendered with templ + HTMX) plus a companion CLI for administrative tasks (creating users, importing players, managing goalie status). Built for a single club with fewer than 10 users — the code is intentionally simple rather than scaled.

## Commands

Dev environment (runs templ, air, tailwind, and asset-sync watchers together):
```
make dev
```

Individual watchers, if you only need one:
```
make dev/templ      # templ generate -watch, proxies :8080 -> :7331 for live reload
make dev/server      # air: rebuilds/reruns the Go server on .go changes
make dev/tailwind    # tailwind CLI watch, writes cmd/web/static/tailwind.css
make dev/sync_assets # notifies the templ proxy when static CSS changes
```

Production build (arm64, for Raspberry Pi deploy target):
```
make prod/build-arm64
```

Deploy to the configured Raspberry Pi (requires `.deploy.conf`, copy from `.deploy.conf.example`):
```
./deploy.sh
```

Tests:
```
go test ./...                                 # all packages
go test ./internal/practice/...               # single package
go test ./internal/practice/ -run TestDistribute  # single test
```

templ files must be generated before `go build`/`go test`/`go vet` will see current markup — `*_templ.go` files are committed, but if you edit a `.templ` file, run `templ generate` (or have `make dev/templ` running) before compiling.

Database schema changes: add a SQL file to `sql/schema/` (goose migration, timestamp-prefixed filename) and/or a query to `sql/queries/`, then regenerate the Go query code:
```
sqlc generate
```
Migrations run automatically against `DB_PATH` on every server/CLI startup (`db.RunMigrations`, via goose, embedded FS in `sql/migrations.go`).

Required env vars (see `.env.example`): `DB_PATH`, `GOOSE_MIGRATION_DIR`/`GOOSE_DRIVER`/`GOOSE_DBSTRING` (goose CLI use), `POSTHOG_API_KEY` (optional — server logs a warning and runs without analytics if unset).

## Architecture

**Two binaries, shared internal packages.**
- `cmd/web` — the HTTP server (net/http `ServeMux`, no framework). `main.go` wires routes; `handlers_*.go` hold handlers grouped by resource (players, practice, auth); `middleware.go` has auth/analytics middleware; `components/` holds `.templ` templates (HTML) and their generated `_templ.go`; `view/` holds presentation-layer structs (e.g. `view.Player`, `view.Practice`) that adapt `internal/` domain types for templates — handlers convert domain -> view, never render domain types directly.
- `cmd/cli` — an interactive REPL (`promptui`) for user creation and player management, sharing `internal/db` and `internal/auth`/`internal/player` with the web server.
- `internal/` — domain logic, framework-free: `player` (Player type, scoring, MyClub import), `practice` (Practice/team distribution logic), `file` (parsing MyClub/attendance Excel exports via `excelize`), `auth` (password hashing, sessions, request-context user info), `db` (sqlc-generated queries + DB init/migration glue).

**Data flow for the core "create practice" path** (see `handleCreatePractice` in `cmd/web/handlers_practice.go`): uploaded attendance Excel -> `internal/file` parses rows and filters to confirmed attendees -> matched against existing players in DB by MyClub ID -> `internal/player.GetPreferences` splits into goalies/field players -> `internal/practice.Distribute` alternates players into two teams (evening out team sizes, not just scores) -> practice + team assignments persisted in one DB transaction (`cfg.db.Begin()` / `queryTx := cfg.queries.WithTx(tx)`).

**Database**: SQLite (`modernc.org/sqlite`, pure Go, no cgo). Schema lives in `sql/schema/*.sql` (goose migrations), queries in `sql/queries/*.sql`, generated code in `internal/db/*.sql.go` via `sqlc` (config: `sqlc.yaml`). Never hand-edit `internal/db/*.sql.go` — edit the `.sql` and regenerate.

**Auth**: session-cookie based (`user_session_id`), not JWT. `authMiddleware` (in `cmd/web/middleware.go`) runs on every request, looks up the session, and stashes `auth.UserInfo` on the request context regardless of login state; `requireAuth` wraps individual routes and redirects to `/login` (or sends `HX-Redirect` for HTMX requests) if not logged in. Users are provisioned only via the CLI (`cmd/cli`), not self-service signup — see `spec/auth.md` for the original auth requirements.

**Frontend**: server-rendered `templ` components styled with Tailwind v4 (via bunx, no bundler/build step beyond the CLI), interactivity via HTMX (`HX-Redirect` headers, `HX-Request` header checks in handlers/middleware) rather than client-side JS/a SPA framework.

**Analytics**: PostHog, loaded conditionally based on `POSTHOG_API_KEY`; the key is threaded through request context (`cmd/web/analytics`) so templates can decide whether to emit the tracking snippet.

## Notes

- `spec/` contains lightweight product specs (e.g. `spec/auth.md`) with checkbox-style requirements — check here for the intended behavior/rationale behind a feature before changing it.
- `spec/improvements.md` is a running backlog of not-yet-built ideas, organized by user flow.
- `docs/` has deployment and security-checklist references for the Raspberry Pi production setup.
- Use `docs/graphs.md` to draw any Mermaid diagrams that need to be easily referenced during conversations.