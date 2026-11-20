# Milestone 01: Flexible attendance import

**Status:** done · **Branch:** feature/flexible-attendance-import
<!-- Status: draft → red-teamed → in-progress → done -->

## Why
Serves **Setting up the practice, step 3** (sub-step 1, importing the MyClub attendance report). Takes on no `I…` idea; I3 is the closest neighbour and stays out (see Scope). The fuller review step after import is I25.

Assumed impact: fewer failed practice set-ups, and an organizer other than the primary one can set up a practice without help.

Today setting up a practice works reliably only because the primary organizer knows exactly which columns to export from MyClub and in which order. MyClub lets the exporter choose the columns and reorder them, so any other organizer who exports the report their own way gets a file Tiimit can't read. They get no feedback about what went wrong (or, for some files, no response at all), so they can't set up the practice.

## Decided
<!-- Choices Human-PM made at kickoff and step 3, written down so they live
     in the spec rather than in one session's memory. -->
- The import finds the columns it needs by their header names, not by position. MyClub lets the exporter pick which fields to export and reorder them, so position is not stable.
- The three required columns are `Jäsennumero`, `Suku- ja etunimi` and `Osallistuu`. These are the same labels the organizer ticks in MyClub's export dialog (R11), so the error message can name them as they are. Any other column is ignored.
- The name column is required even though creating the practice matches players by MyClub ID only: it makes the skipped-player message actionable (R5).
- Only the Finnish export is in scope. The English export translates the sheet name (`Tapahtuma` becomes `Event`) and every header (`Jäsennumero` becomes `Member number`, `Osallistuu` becomes `Coming`).
- `Osallistuu` is a fixed column name in MyClub, not something set per event.
- **The header row** is the first row in which any cell exactly matches one of the three required names, ignoring surrounding spaces. Player rows are every row below it (R2). Searching for "a row with all three" would leave a report with a missing column without any header row, so the missing column could not be named. The text `Osallistuu` also appears as a status in player rows, but the header row always comes first.
- The import reads the first sheet of the file, whatever its name (R10).
- Which rows are skipped and reported (R4): a row whose attendance value is not one of `Osallistuu`, `Ei osallistu`, `Ei vastausta`; or a row marked `Osallistuu` whose MyClub ID is missing or not a number. A row with an empty attendance cell counts as not attending and is not reported (R3).
- A malformed player row is skipped and reported; it does not fail the whole import. If no player row can be read at all, no practice is created and the organizer gets a message (R9).
- When players were skipped, the organizer stays on the result instead of being sent to the practice: the page lists each skipped player by name with the reason, and a link opens the new practice. Without skipped players, the import goes to the practice as today (R8).
- Kickoff context: tech review item #3 (`docs/tech-review-2026-11.md`), where a short row or a near-empty sheet panics the import.

## Scope
- **In:**
  - The import reads the three required columns by their Finnish header names, in any order, from wherever the header row sits, and ignores any other columns.
  - A report missing a required column is rejected with one short message that names the missing columns and lists the three columns to include, exported in Finnish.
  - An attending player Tiimit can't read is skipped, the rest of the practice is set up, and the organizer sees who was skipped and why.
  - No file makes the import crash; every failure ends in a message to the organizer (tech review #3).
  - Server logs about skipped rows record row numbers and MyClub IDs, never player names (R8).
- **Out:**
  - Reading English or Swedish reports. They are treated like any other unreadable file: the message lists the Finnish columns (R1).
  - Alternative ways to give the name (separate first and last name columns).
  - Showing attending players who are not in Tiimit's roster (I3). Today they are left out silently, and this milestone does not change that.
  - A review step where the organizer checks the import before the practice is created (I25).
  - The date in the file name. It is still required as it is today.

## Acceptance conditions (BDD)
<!-- Group conditions into named scenarios. Name each by what the user
     does, in the situation that makes it distinct; the result belongs
     in the Then lines.
     Exception: a scenario with no real user (e.g. an infrastructure or
     build check) is named by the behaviour being proven, not forced into
     a "coach" frame. Each scenario becomes a `go test` case, written
     before implementation. Plain language only, readable without code
     knowledge; every scenario has Given, When and Then, one condition per
     line. Mark one a Go test can't observe with `Verification: manual`. -->

### S1: Organizer imports a report with reordered and extra columns
- Given four players are in Tiimit's roster
- And the organizer has a Finnish attendance report with the attendance column first, the name column last and an extra column in between
- And the report shows all four players attending
- When the organizer imports the report
- Then a new practice is created
- And all four players are distributed into the two teams

### S2: Organizer imports a report with a missing column
- Given the organizer has a Finnish attendance report with the name and attendance columns but no MyClub ID column
- When the organizer imports the report
- Then no practice is created
- And the organizer sees that the column "Jäsennumero" is missing
- And the organizer sees the three columns to include when exporting

### S3: Organizer imports a report exported in English
- Given the organizer has an attendance report exported from MyClub in English
- When the organizer imports the report
- Then no practice is created
- And the organizer sees the three Finnish columns to include
- And the organizer sees that the report must be exported in Finnish

### S4: Organizer imports a report with an unreadable attending player
- Given four players are in Tiimit's roster
- And the organizer has a Finnish attendance report where all four players are attending
- And the fourth player's MyClub ID is missing from the report
- When the organizer imports the report
- Then a new practice is created with the other three players
- And the organizer sees the fourth player's name
- And the organizer sees that the MyClub ID was missing
- And the organizer can open the new practice from the same page

## Test notes
<!-- How each scenario is tested, by number: fixtures, what the test drives,
     how a manual one is checked. Mechanics live here, not in the scenarios. -->
- All: acceptance tests drive the real practice-creation handler through `httptest` with a multipart upload. That needs a new upload helper; it lives in `cmd/web/handlers_practice_test.go` with the tests, since only they use it (ARCHITECTURE.md conventions). Each test builds its report in memory with excelize from a small table of made-up names, so no real export is ever committed. The uploaded file name carries a date (`yyyy-mm-dd`), since the handler still requires one. Players are put in the roster with `addPlayer`.
- S1: asserts the practice page lists all four made-up names. It also replaces the plain "required columns in the usual order" case, which S1 covers.
- S2: asserts the page names `Jäsennumero` as missing and lists all three columns, and that no practice exists afterwards. The other single missing columns are unit tests.
- S3: the fixture uses the English headers and sheet name seen in a real export (`Event`; `Member number`, `Last and first name`, `Gender`, `Coming`). The same message covers any spreadsheet that is not a MyClub report, so there is no separate scenario for that.
- S4: asserts the response is not a redirect, contains the fourth made-up name and the reason, and links to the new practice; and that the practice holds the other three players.
- Unit tests (table-driven, in `internal/file`), written with the acceptance tests:
  - finding the header row: header on the usual row, with an extra title row above it, and with only some of the required columns;
  - column positions in any order, extra columns ignored;
  - each required column missing on its own, and all three missing;
  - a sheet with no header row, a sheet with fewer rows than the usual header position, a header with no player rows, and a file where every player row is skipped (R9);
  - row skipping (R4, R3): unknown attendance value is skipped; attending with a missing or non-numeric MyClub ID is skipped; a short row with no attendance cell is not attending and not skipped (regression for tech review #3); a row with `Ei vastausta` is kept.
- Existing tests in `internal/file/file_test.go` that assume fixed column positions are rewritten to the new behaviour, not kept alongside it.

## Open questions
<!-- For the red-team pass (workflow step 2) and Human-PM to resolve (step 3).
     A resolved question keeps its number and moves to the decisions log. -->
None open.

## Decisions log
<!-- Red-team findings (R1, R2, …) with Human-PM's decision on each, resolved
     questions (Q…), and decisions made while building. One line of why each. -->
- **Scenario numbering:** the draft's S1 to S9 were merged and renumbered to S1 to S4 at step 3, at Human-PM's request. This breaks the "never renumber" rule on purpose: nothing outside this spec referred to the draft numbers. The R entries below quote the draft numbers.
- **R1** Accepted. Draft S7 (Swedish) was untestable, and S6/S7/S8 tested one behaviour. Merged into S3; every rejection lists the Finnish columns.
- **R2** Accepted. Header row rule written into Decided; it is what lets S2 name the missing column.
- **R3** Accepted. S4 uses an attending player with a missing MyClub ID; an empty attendance cell means not attending, silently.
- **R4** Accepted. Defines which rows are skipped (Decided).
- **R5** Name column stays required, so the organizer gets a skipped player's name rather than a row number.
- **R6** Accepted. Draft S1 folded into S5, draft S2 to S4 into one scenario with per-column unit tests; then renumbered (see above).
- **R7** Accepted at first, then reversed by Human-PM: names describe what the organizer does, which reads more clearly; the outcome is in the Then lines.
- **R8** Accepted. No redirect when players were skipped; names never in URLs or logs. A full review step is a later idea (I25).
- **R9** Accepted. No readable player rows means no practice and a message; unit test only.
- **R10** Accepted. Read the first sheet.
- **R11** Confirmed by Human-PM: the headers match the labels in MyClub's export dialog.
- **R12** Done. Removed the duplicated Open questions section left from the template.
- **Built, slice 1:** the upload helper lives with the tests in `handlers_practice_test.go`, not in `helpers_test.go`: only those tests use it.
- **Built, slice 2:** an empty name on an attending row with a valid MyClub ID is imported, since players are matched by MyClub ID (R4). The old row constructors that rejected it were removed.
- **Code review, round 1** (`/code-review high`): fixed: a report with no readable attending player now gets the no-players message instead of a later "try again" (blocking); the import result replaces the upload form so the report can't be sent twice; logs show the MyClub ID only when it is a number, since the cell may hold free text; skipped rows are checked on the error path in unit tests; tech review #3 marked done. Declined: listing skipped players with the no-players message (rare, R9 chose a plain message; revisit with I25); IDs with digit grouping (real exports read as plain digits).
- **Q1** Resolved by R8. **Q2** by R10. **Q3** by R9. **Q4** by R1. **Q5** by R2 (unit test, no scenario). **Q6** by R4.

## Slices
<!-- Written at refinement (workflow step 4). Each slice is one commit and
     ends with Human-PM's test checklist before the next starts. -->
1. **Read columns by header name.** The import finds the header row and the three columns by name on the first sheet, and rejects a report with missing columns with the Finnish column list. No input panics. Covers S1, S2, S3 and the header, column and empty-sheet unit tests.
2. **Skip and report unreadable players.** Rows are skipped by the R4 rule instead of failing the import; when any were skipped, the organizer sees them with the reason and a link to the practice instead of being redirected. Logs carry row numbers and MyClub IDs only. Covers S4 and the row-skipping unit tests.
