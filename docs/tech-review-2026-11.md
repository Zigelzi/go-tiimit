# Technical review, 2026-11-18

A one-off review of the whole codebase at `340c56f`, made when switching to the agentic workflow. It lists the technical problems that make Tiimit harder to change safely, by a human or an LLM, in the order Claude suggests to fix them.

This file is a dated snapshot, not a tracker. The spec rule is "one fact lives in exactly one file", so once you have reviewed it, move each accepted topic into `spec/product-notes.md` (Technical section, next free `I` number) or straight into a `fix/` or `chore/` branch, and delete this file or keep it as history. Topics already tracked there are listed in [Already tracked](#already-tracked) and not repeated.

Everything was read by hand: handlers, middleware, `internal/`, SQL, templates, view, CLI, tests, Makefile and deploy script. Items marked **verified** were reproduced (a throwaway test, `make test`, or `git log`); the rest are from reading the code.

## How the priorities were set

Same order as the constitution's goals, applied to technical work:

1. **Harm now.** Anything that can lose or corrupt the club's data, leak personal data, or crash the core flow comes first, however small.
2. **Leverage for the next milestones.** Almost every open idea in "Setting up the practice" and "Running the practice" (I3, I4, I5, I7, I10, I24) changes practice creation or the practice page. That code has no HTTP tests and lives in one 140-line handler. Making it safe to change pays back on each of those milestones.
3. **Cost of reading.** Stale docs, duplication and missing doc comments make every session (human or LLM) re-derive things. Fix these opportunistically, when a milestone touches the code anyway.

Sizes: **XS** a few lines, one commit. **S** one small `fix/` or `chore/` branch. **M** a milestone of its own.

## Summary

| #   | Topic                                                         | Size           | Why this position                                                            |
| --- | ------------------------------------------------------------- | -------------- | ---------------------------------------------------------------------------- |
| 1   | ~~Docs describe a player import that no longer exists~~ (done) | XS             | Every session reads these docs first; wrong facts there cost every milestone |
| 2   | ~~Creating a practice ignores the commit error~~ (done)       | XS             | Can report success for a practice that was never saved                       |
| 3   | ~~Attendance import panics on a short row or a near-empty sheet~~ (done) | XS             | Crash in the core flow, triggered by a slightly different export             |
| 4   | CLI lowercases and truncates usernames and passwords          | XS             | Weakens passwords silently and makes mixed-case logins fail                  |
| 5   | PostHog may receive player names                              | S              | Privacy rule in the constitution; needs a check, then a config change        |
| 6   | Production runs with `env = "development"`                    | XS             | Static files are never cached on phones on the pitch                         |
| 7   | SQLite opened without `busy_timeout`, `foreign_keys` or WAL   | XS–S           | Two coaches ticking vests at once can get "database is locked"               |
| 8   | No acceptance tests for the practice flows                    | S              | Prerequisite for safely changing the code most ideas touch                   |
| 9   | Move practice creation out of the handler                     | S              | Already a "Not moved yet" line; do it right after #8                         |
| 10  | `make check`: vet, staticcheck and pinned tool versions       | XS             | Faster, cheaper feedback for every slice; staticcheck can't run today        |
| 11  | One convention for error responses and "not found"            | S              | Today a missing player renders an empty row and full pages fail blank        |
| 12  | Handler boilerplate and duplicated helpers                    | S, incremental | Biggest share of handler lines; each copy drifts separately                  |
| 13  | Doc comments on exported identifiers                          | incremental    | `go doc` is the documented discovery tool, but it shows bare signatures      |
| 14  | Dead code and stale TODOs                                     | XS             | Misleads readers about what is used                                          |
| 15  | Server hardening for a public endpoint                        | S              | Real but low-probability risks behind nginx                                  |
| 16  | Small things                                                  | XS each        | Fix in passing                                                               |

## Priority 1: defects and risks (small fix branches)

### 1. Docs describe a player import that no longer exists (verified) — done

**Status (2026-11-19):** fixed on `chore/stale-player-import-docs`. The player-import claims were removed from `CLAUDE.md`, `README.md`, `spec/domain-model.md` and `ARCHITECTURE.md`; the two remaining I19 items (the raw session token in the logout log, the "failed toget" typo) were fixed in code and I19 was removed from `product-notes.md`. On review, the quotes below for `README.md`, `ARCHITECTURE.md` and `domain-model.md` did not all match the files at `340c56f`, so the text that follows is the original finding, kept as history.

Bulk-importing players from a MyClub export was removed in `bd04715` ("Reduce CLI to user creation only"). Players are now only added one by one in the web app. Four files still say otherwise:

- `CLAUDE.md:19`: "A coach imports the club's players from a MyClub Excel export".
- `README.md:15`: "Import players from MyClub to Tiimit by using Excel exported from MyClub."
- `spec/domain-model.md`, "Imports (not stored)": "the club's player export (to create players)", and line 36 "import and edit players".
- `ARCHITECTURE.md:58`: "imports into players" as a job of `internal/player`.

Also, I19 in `product-notes.md` points at `legacy/improvements.md#error-handling-cleanups`, whose file:line references have drifted. Checked against today's code: the broken format verb, the repeated inactive/active message and the bare `log.Println` calls are already fixed; the raw session token in the log (`handlers_auth.go:82`) and the "failed toget" typo (`handlers_practice.go:320`, `:364`) are still there.

**Why first:** the workflow's whole premise is that a fresh session learns the project from these files. A wrong fact in `CLAUDE.md` is read in every session and can steer a milestone (e.g. a red-team asking about player imports). **Fix:** correct the four places; rewrite I19 to list only the two remaining items without line numbers.

### 2. Creating a practice ignores the commit error

**Status (2026-11-19):** fixed on `fix/practice-commit-error`: the commit error is checked and shown as the error banner, the transaction uses the request context, and the player logs carry the error and the id but no name. The text below is the original finding, kept as history.

`cmd/web/handlers_practice.go:183` calls `tx.Commit()` without checking its error, then sends `HX-Redirect` to the new practice. If the commit fails (for example a locked database, see #7), the coach is sent to a practice that does not exist and gets a blank 404. Two smaller issues in the same block: `cfg.db.Begin()` (`:143`) ignores the request context, and the two `AddPlayerToPractice` failure logs (`:165`, `:178`) omit `err` but include the player's name (personal data in logs).

**Fix:** check the commit error and render the error banner; use `BeginTx(r.Context(), nil)`; log ids and `err`, not names. A test for the commit path is hard to force; the code fix is three lines and reads correctly on its own. Better still, fold it into #9.

### 3. Attendance import panics on short rows or a near-empty sheet (verified), done

**Status (2026-11-20):** fixed on `feature/flexible-attendance-import` (milestone 01). The import now finds the header row and columns by name, reads cells with a bounds check, skips and reports unreadable attending players, and gives every failure a message. The text that follows is the original finding, kept as history.

`internal/file/import_attendance_rows.go:28` slices `rows[4:]` and `:46` indexes `row[3]` with no length checks. excelize's `GetRows` drops trailing empty cells, so:

- a player row with an empty attendance cell panics with `index out of range [3] with length 3`;
- a "Tapahtuma" sheet with fewer than 4 rows panics with `slice bounds out of range [4:1]`.

Both reproduced with a throwaway test. `net/http` recovers the panic, so the server survives, but the coach gets no response and no error banner, only a dead button. The existing `TestParsingRowsWithMissingFieldsReturnsError` passes because it uses rows with an empty string in column 4, which a real export does not produce.

**Fix:** check lengths and return an error (the handler already turns that into a banner); add both cases to the table test. Consider whether one malformed row should fail the whole import or be skipped and reported (that part is a product question, close to I3).

### 4. The CLI lowercases and truncates usernames and passwords

`cmd/cli/main.go:76` and `:80` pass input through `cleanInput`, which lowercases it and keeps only the first space-separated word. A password `Kesä Treeni 2026` is stored as the hash of `kesä`. The login form does not lowercase, so anyone who types their password as they chose it is rejected, and the stored password is weaker than intended. The password is also echoed to the terminal, and an empty username is accepted.

**Fix:** read the username with `strings.TrimSpace` only, read the password without echo (`golang.org/x/term.ReadPassword`), reject an empty username. Existing users keep working because their stored hashes already match the lowercase first word. Until I17 replaces the CLI, it is the only way to create users, so it is worth the ten lines.

### 5. PostHog may receive player names

`cmd/web/static/analytics.js:54` initialises PostHog with its defaults and `person_profiles: "always"`, and identifies the user with their username. PostHog's default autocapture records the text of clicked elements. On the practice page, ticking a vest is a click on a `<label>` whose text is the player's name (`components/player_card.templ`), so names are likely sent to PostHog with every tick. If session replay is switched on in the PostHog project, whole pages are recorded. The constitution treats player names as personal data that only logged-in users may see.

Also: when `POSTHOG_API_KEY` is empty, `base.templ:20` still calls `initPostHog("")`, which loads PostHog's script from its CDN on every page.

**Fix:** first check in PostHog's activity view whether `$el_text` contains names. Then either turn autocapture off and capture the few events you care about by name, or set `mask_all_text: true` and `mask_all_element_attributes: true`; and skip the init when the key is empty. Decide whether usernames should go to PostHog at all (the user id alone is enough to count users).

### 6. Production runs with `env = "development"`

`cmd/web/main.go:42-43` hardcodes the address and `env: "development"`, so `disableCacheInDevMode` (`:90`) sends `Cache-Control: no-store` for every static file in production. Phones on the pitch re-download the stylesheet and scripts on every page. I20 already plans to read the address and environment from an env file; this is the cheap part of it that can land on its own: read `APP_ENV` (default `production`) and `ADDR` from the environment, and add them to `.env.example`.

### 7. SQLite is opened without `busy_timeout`, `foreign_keys` or WAL

`internal/db/config.go:23` opens the database with no pragmas and allows 10 open connections. With SQLite's defaults, a write that arrives while another connection holds the write lock fails immediately with `SQLITE_BUSY` instead of waiting a few milliseconds. Two coaches ticking vests at the same moment (I14 says it can happen) is enough. Foreign keys are off, which `domain-model.md` already notes for deletes.

**Fix:** open with `?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)` (modernc's DSN syntax). Before switching foreign keys on against the live database, run `PRAGMA foreign_key_check;` on a copy: no deletes exist, so it should be clean, but the deploy also migrates live data. Add a short test that two concurrent vest toggles both succeed.

## Priority 2: enablers before the next practice milestones

### 8. No acceptance tests for the practice flows

Milestone 00 deliberately stopped at login and adding a player. The HTTP tests cover 6 behaviours; nothing drives creating a practice from an upload, viewing it, moving a player, ticking a vest or switching team tabs, nor editing a player or setting one inactive. That is the flow the app exists for, and the one I3, I4, I5, I7, I10 and I24 all change.

**Fix:** a `chore/practice-tests` milestone (spec with scenarios for today's behaviour, no visible change). It needs one new helper: building an attendance `.xlsx` in memory with excelize (as the probe for #3 did) with made-up names, and posting it as multipart. Doing this before #9 is what makes #9 a safe refactor.

### 9. Move practice creation out of the handler

Already recorded in `ARCHITECTURE.md` under "Not moved yet". `handleCreatePractice` (`handlers_practice.go:48-186`) parses the file, matches rows to players one query at a time, splits goalies, distributes, and writes in a transaction, with a `// REPO START` comment marking where a function wants to be. It also holds a 28-line commented-out block for unknown players (I3).

**Why here:** each of I3 (show who was skipped), I4/I5 (add players afterwards), I10 (re-roll) and I24 (skip inactive players) needs to change this code. In `internal/practice` as something like `Create(ctx, store, rows) (practiceID, skipped, error)`, each becomes a small change with a unit test; in the handler, each is another branch in a long function tested only end to end. **When:** as the first slice of whichever of those milestones comes first, right after #8. Fix #2 inside it.

### 10. `make check`: vet, staticcheck and pinned tool versions

`make test` runs `go test` only. `go vet` is clean today but nothing runs it. `staticcheck` is installed but built with Go 1.24, so it refuses to run on this Go 1.26 module (verified). `templ` and `sqlc` are taken from `PATH`, so the generated code depends on whichever version each machine has; `340c56f` upgraded them by hand.

**Fix:** declare the tools in `go.mod` with the `tool` directive (`go get -tool github.com/a-h/templ/cmd/templ@v0.3.1118`, the same for sqlc and staticcheck) and call them as `go tool templ`, `go tool sqlc`, `go tool staticcheck`. Add `make check` = generate, vet, staticcheck, test, and have `deploy.sh` run it. A session then has one command that says "done or not", and the versions live in one place. This is a `chore` straight on main per your memory rule for tool upgrades.

### 11. One convention for error responses and "not found"

Handlers fail in three different ways:

- full-page handlers write a bare status with an empty body (`handleIndexPage`, `handleViewPracticePage`, `handleViewPlayersPage`), so the coach sees a blank page;
- htmx handlers call `renderError`, which replies `200 OK` with `HX-Reswap: none`;
- `handleEditPlayerRow` (`handlers_players.go:146-151`) and `handleViewPlayerRow` (`:258-264`) treat "no such player" as success and render a row for an empty player with id 0.

`ARCHITECTURE.md` says errors go through the error banner, but there is no rule for full pages, for status codes, or for "not found". Tests cannot assert a failure by status. **Fix:** write the rule into `ARCHITECTURE.md` conventions (e.g. full pages render a small error page with the right status; htmx fragments keep the banner; a missing id is a 404 in both), then apply it handler by handler. Fix the two empty-row bugs now; they are two lines each.

## Priority 3: clean up when a milestone touches the code

### 12. Handler boilerplate and duplicated helpers

- Parsing an id from the path appears 11 times, 7 with `strconv.Atoi` and 4 with `ParseInt`, each with its own message. One `pathID(r, name) (int64, error)` helper.
- The loop db rows → domain → view types is copied in four handlers (`handlers_players.go:39-50`, `:316-324`, `handlers_practice.go:325-329`, `:369-373`).
- `handleSetPlayerInactive` (`handlers_players.go:268-327`) has two near-identical branches.
- `handleMovePlayer` (`handlers_practice.go:236-287`) reads the team, then writes the other one, in two near-identical `case` blocks. That is the read-then-write pattern `ARCHITECTURE.md` rule 5 forbids for vests. One `UPDATE ... SET team_number = 3 - team_number ... RETURNING team_number` fixes both. Do it with I7.
- `SortByScore` exists in both `player` and `practice`; `roundToTwoDecimals` is defined twice (`player/scoring.go:15`, `practice/practice.go:112`); building a `PracticePlayer` from a row is written twice (`practice.go:43` has a TODO saying so).
- `handleIndexPage` loads every active player only to count them; a `COUNT(*)` query says what it means.

### 13. Doc comments on exported identifiers

`CLAUDE.md` and `ARCHITECTURE.md` both say "run `go doc ./internal/practice` to see what exists", and the comment convention is written down, but almost no exported identifier in `internal/` or `cmd/web/view` has a doc comment, so `go doc` shows bare signatures. For example, `Distribute` doesn't say that goalies are placed after field players, or that ties go to team 1. Add them as each package is touched; start with `internal/practice` and `internal/player`, where the rules are.

### 14. Dead code and stale TODOs

`practice.FromDB` and `Practice.UnknownPlayers` (`practice/practice.go:19`, `:28`), `file.FileName` (`file/models.go:9`), the commented-out block in `handleCreatePractice`, the `error` that `file.GetAttendanceRowsByStatus` can never return, the negative-score check in `Player.Score` that `New`'s clamping makes unreachable, and the TODOs in `file/models.go:25` and `import_attendance_rows.go:27`. Each one makes a reader (or a red-team agent) ask "is this used?". Delete them; git keeps the history. A `deadcode` run would confirm the list (it needs downloading `golang.org/x/tools/cmd/deadcode`, which I did not do).

### 15. Server hardening for a public endpoint

nginx sits in front, so these are lower priority, but the app is on the public web:

- `http.Server` (`main.go:47`) has no timeouts; set at least `ReadHeaderTimeout`.
- Login has no throttling, and the minimum password length is 6. A per-username delay after a few failures is a few lines.
- The upload has no size cap (`http.MaxBytesReader`); only logged-in users can reach it.
- No graceful shutdown; `log.Fatal` skips the deferred `Close`. Harmless with SQLite today, worth having once WAL is on.

I22 and I23 (cookie flags, CSRF) belong here too and are already tracked. Note for them: Go 1.25+ has `http.CrossOriginProtection`, which covers I23 by wrapping the mux in one line, so both together are an XS change.

### 16. Small things

- `player_card.templ` wraps each `<li>` in a `<div>` inside the `<ol>` (invalid HTML), and uses the bare player id as the checkbox's element id.
- `view.PracticeSummary.DaysInPast` and `IsUpcoming`, and `index.templ`, call `time.Now()` directly, which makes the front page untestable around dates. Pass "today" in.
- The app is "Tiumo" in page titles and `deploy.sh`, and "Tiimit" everywhere else. Pick one so a grep finds everything.
- Scores are whole numbers in the domain model, but the edit form accepts any float.
- `.env.example` carries `GOOSE_*` variables the app does not read. Say they are for the goose CLI, or remove them.
- Session expiry relies on `expires_at` being written in UTC and compared as text with `CURRENT_TIMESTAMP`; it is correct today (`handlers_auth.go:54` uses UTC) but worth a comment.

## Already tracked

Not repeated above; listed so you can see the whole picture in one place.

- I18: the auth middleware treats a database error as "logged out".
- ~~I19: error-handling cleanups.~~ Done, see #1.
- I20: deploy like Roundtrip, including database snapshots before each deploy. This is the most valuable item on the whole list for data safety; it is P1 in spirit, already well specified, and an M.
- I22, I23: cookie flags and CSRF. See #15 for a shortcut.
- I24: inactive players are distributed. Fits naturally with #9.
- `domain-model.md`: foreign keys are off. See #7.

## What is in good shape

Worth keeping as it is, so the list above doesn't read as "everything is wrong":

- The layering is clear and mostly followed, and `ARCHITECTURE.md` honestly lists where it isn't.
- The spec workflow, the numbering rules and the "one fact in one file" rule make sessions cheap to start.
- The test harness from milestone 00 is readable and fast, and its failure messages print the visible page.
- `practice.Distribute` is pure and has property-style tests (sizes, goalies, everyone placed once, balance).
- The vest toggle is a single SQL statement, as rule 5 asks.
- sqlc, embedded migrations and the pinned, checksum-verified Tailwind binary keep the build reproducible.
