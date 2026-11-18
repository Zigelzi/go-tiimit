# Milestone 00: Test harness

**Status:** done · **Branch:** chore/test-harness

<!-- Status: draft → red-teamed → in-progress → done -->

## Why

Before building any features with the agentic workflow, every later milestone must be able to write its acceptance scenarios as Go tests that drive the real web app. Today only the `internal/` packages have tests, so nothing checks what a coach actually sees, and a change to a handler or template can break a page without any test failing.
This milestone serves no flow step directly; it is the groundwork for all of them, and goal 1 (fast feedback) and goal 2 (quality) depend on it. It takes on no ideas.

Thin slice: a test can start the app on its own fresh, empty database, log in as a coach and check a real page. That is proven on a few behaviours that already exist, and `make test` runs everything.

## Decided

- Follow Roundtrip's harness (`../roundtrip/cmd/web/helpers_test.go`): a test builds the app on a fresh temp SQLite file and calls the router directly; tests run one at a time.
- The Tailwind standalone CLI is pinned and downloaded into a gitignored `./bin/` (version: see the building decisions below).

## Scope

- **In:**
  - The test setup: a fresh database per test in a temp directory that never touches `tiimit.db`, a logged-in coach helper, and made-up players (never real names). Test data is written in the test code; nothing is copied from a MyClub export. PostHog stays off in tests: the test builds the app with an empty key, never reading the environment.
  - A refactor with no visible change, so a test can build the app: the routes and middleware move out of `main()` into a function that returns the handler, and the database opens from a path passed in (`main` and the CLI read `DB_PATH` and pass it).
  - Prove it on existing behaviour:
    - a logged-out visitor opening the players page is sent to login;
    - a logged-out visitor opening the front page sees no player or practice data;
    - a coach logs in and sees the front page with practices and the player count;
    - wrong credentials show an error;
    - a coach adds a player and sees them in the list.
  - `make test`, `make generate` (sqlc, templ, Tailwind) and `make run`. `make test` runs templ generation first (all it needs); the builds and `make run` run `make generate` first, so a fresh checkout works with one command.
  - Switch Tailwind from `bunx` to the standalone CLI and remove `package.json`, `bun.lock`, `node_modules` and the unused `tailwind.config.js`, which also drops prettier and its Tailwind plugin. `make dev`, `make prod/build-arm64` and `./deploy.sh` must keep working, and the pages must look exactly as before.
  - `./deploy.sh` runs `make generate test` instead of `go test ./...`, so it works on a fresh checkout and tests the code it ships. Its only change.
  - Update `CLAUDE.md`'s command list and `ARCHITECTURE.md`'s test convention (and its "registered in `main.go`" row) once this lands.
- **Out:**
  - Any change a coach can see, including fixing bugs the new tests uncover. Inactive players still being distributed is one; write those down as ideas instead.
  - Moving practice creation out of the handler (the "Not moved yet" line in `ARCHITECTURE.md`).
  - Deploy changes (I20), including reading the address and environment from config.
  - Tests for the practice pages, and uploading Excel files in tests.
  - Committing generated templ code (`*_templ.go`): `make generate` covers fresh checkouts instead.
  - Reaching the dev server from a phone on the home network.
  - CI.
  - The login cookie's flags and CSRF protection (I22, I23).

## Acceptance conditions (BDD)

### S1: A logged-out visitor is sent to the login page

- Given nobody is logged in
- When a visitor opens the players page
- Then they see the login page

### S2: A coach logs in and sees the front page

- Given a coach account
- And the team has players
- And the team has earlier practices
- When the coach logs in with the right username and password
- Then they are taken to the front page
- And it shows the earlier practices
- And it shows the number of active players in the team

### S3: Wrong credentials show an error

- Given a coach account
- When someone logs in with a wrong password
- Then they stay on the login page
- And they see an error message

### S4: A coach adds a player and sees them in the list

- Given a logged-in coach
- When the coach adds a player
- Then the player appears in the players list

### S5: A fresh checkout tests with one command

- Given a fresh copy of the code with no generated files
- When the developer runs the one test command
- Then the code is generated first
- And every test runs and passes
- Verification: manual

### S6: Pages look exactly as before the Tailwind switch

- Given the app built with the standalone Tailwind CLI
- When a coach opens each page
- Then it looks the same as before the switch
- Verification: manual

### S7: The front page shows a logged-out visitor no club data

- Given nobody is logged in
- And the club has players and practices
- When a visitor opens the front page
- Then they only see a login link
- And they don't see practices
- And they don't see the number of active players
- And they can't start setting up a practice

## Test notes

- All: tests live in `cmd/web/*_test.go`. A helper builds the app on `t.TempDir()/test.db`, runs the migrations (Goose logger silenced) and calls the router with `httptest.NewRecorder`, as Roundtrip does. No `t.Parallel`: Goose keeps package-level state. Players and practices are inserted with made-up names through the sqlc queries.
- Coach helper: creates the user the way `cmd/cli` does (hash the password, insert the user), then logs in through the real `POST /login` and keeps the session cookie for later requests, so every logged-in test goes through the real cookie path.
- S1: `GET /players` with no cookie answers a redirect to `/login`. Also: the same request sent as htmx answers 401 with `HX-Redirect: /login`.
- S2: `POST /login` sent as htmx answers `HX-Redirect: /` and sets the session cookie. Then `GET /` with that cookie shows the seeded practices' dates and the active-player count.
- S3: the reply has no session cookie and no `HX-Redirect`, and contains "Incorrect username or password". The reply is a fragment, so "stays on the login page" means "not redirected".
- S4: `POST /players` as the logged-in coach, then `GET /players` contains the new made-up name.
- S5: Human-PM checks it in a throwaway `git clone` of the branch in a temp directory, never by cleaning the working copy (`git clean -X` would delete `.env`, `tiimit.db` and Excel files). The developer machine needs Go, the templ CLI and sqlc; the Tailwind binary downloads itself.
- S6: mostly automatic: before removing bun, build the CSS once with bun, then with the standalone CLI with the same flags, and compare the two files. Identical output means identical pages. Human-PM glances at the pages as a backstop.
- S7: seed players and a practice; `GET /` with no cookie contains a link to `/login` and neither the practice's date, the player count nor a link to `/practices/new`.

## Open questions

None open.

## Decisions log

- **R1** A test can't build the app without risking `tiimit.db`: routing lives in `main()` and the database path comes only from `DB_PATH`. Accepted: refactor to a routes function and a database opened from a path (In).
- **R2** The login reply sends the coach to the front page, not the players list. Human-PM rewrote S2 to the front page.
- **R3** Simulating a fresh checkout with `git clean` destroys local data. Checked against Roundtrip: it verifies its one-command build by hand, not in a test. S5 is manual, run in a throwaway clone; the command name left the scenario text.
- **R4** `make test` needs only templ output. Accepted: `make test` runs templ generation only, `make generate` runs all three. Human-PM updates the installed templ CLI to match `go.mod` (v0.3.1001); no `go tool` pin.
- **R5** Accepted: pinned Tailwind downloaded into `./bin/`, `tailwind.config.js` deleted (v4 ignores it), S6 proven by comparing the built CSS.
- **R6** Accepted: `./deploy.sh` runs `make test`, its only change.
- **R7** Goose keeps package-level state, so parallel tests would race. Checked against Roundtrip: its tests run one at a time and silence Goose's logger. Same here.
- **R8** S7 could not fail because the front page never shows names. Human-PM rewrote S7 to check practices, the player count and the setup link.
- **R9** S2's Given joined two conditions. Declined: the intent was clear (S2 was rewritten anyway).
- **R10** No decision needed: the coach is created as the CLI does and logs in through the real login (Test notes).
- **R11 / Q1** S1–S4 and S7 are enough to prove the harness; the logged-out htmx request is one extra check in S1's test, not a new scenario.
- **R12** Accepted: the cookie's missing `Secure` flag and the missing CSRF protection become ideas I22 and I23 (Out); test data is never copied from an export (In).
- Tailwind pinned to 4.1.17, not 4.2.1: `bun.lock` said 4.2.1, but the Pi build ran 4.1.17 from `node_modules` and `make dev` fetched the latest (4.3.5). 4.1.17 is what production serves, and the standalone build's CSS is byte-for-byte identical to it, which proves S6.
- `make run` runs `make generate`, then the server with the `.env` settings, as Roundtrip's does.
- The Tailwind download picks the binary for the machine it runs on (Linux or macOS, x64 or arm64); the CSS is the same on every platform, so the Pi's arm64 binary needs no Tailwind of its own.
- Drop prettier along with `package.json`: no Node toolchain stays in the repo (constitution: Tailwind via the standalone CLI); `templ fmt` formats templates.
- Logged-out access is tested on both the players page (S1, redirect) and the front page (S7), because the front page is public and only hides data.
- Branch prefix `chore/`: the milestone changes no behaviour a coach can see.
- Code review, round 1 (`/code-review high`):
  - Fixed: `./deploy.sh` tested before sqlc regenerated, so it could ship database code the tests never ran (`make generate test`); S2 and S7 checked the practice by its date, which the logged-in header also shows on that day (now by its link); the server's start error was dropped; a failed database ping leaked the handle; an unknown platform downloaded the x64 Tailwind binary (now a clear error); the deploy guide's flowchart and "no `set -e`" paragraph were stale; `/node_modules/` stays ignored for older clones.
  - Declined: sqlc out of `make generate` (Human-PM's decision in R4; the version churn goes away when the installed sqlc matches the committed 1.30.0); Goose's instance-based provider instead of serial tests (Human-PM chose Roundtrip's approach in R7; revisit if the suite gets slow); a shared constant for the session cookie name (the tests state the cookie the app must set, written before the code, as Roundtrip's do).

## Slices

1. **Harness and the first acceptance tests.** The routes function and the database opened from a path; the test helpers; tests for S1, S2, S3, S4, S7; `make test` (templ generation, then `go test ./...`). Covers S1, S2, S3, S4, S7.
   Checklist:
   1. Run `make test`: everything passes, and `tiimit.db`'s modification time hasn't changed.
   2. Run `make dev`, log in, open the players page and add a player: it works as before.
   3. Log out and open `/players`: you land on the login page.
2. **Standalone Tailwind and one-command builds.** The pinned Tailwind download, `make generate` and `make run`, the dev and Pi builds switched to them, `./deploy.sh` on `make test`, bun, `package.json` and `tailwind.config.js` removed, `CLAUDE.md` and `ARCHITECTURE.md` updated. Covers S5, S6.
   Checklist:
   1. In a temp directory, `git clone` the repo, check out `chore/test-harness`, run `make test`: it generates, then all tests pass.
   2. Run `make dev` and open the front page, the players page and a practice: they look as before.
   3. Run `make prod/build-arm64`: `./build/web` is built.
