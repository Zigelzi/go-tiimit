# Milestone 03: Remove the CLI

**Status:** done · **Branch:** chore/remove-cli
<!-- Status: draft → red-teamed → in-progress → done -->

## Why
Serves **User management, step 2**. Takes on **I17** (deprecate the CLI). Split out of milestone 02 (`02-user-registration.md`, 02/Q5) and starts only after milestone 02 is merged, because until then the CLI is the only way to create an account.

Once organizers join with invite links (milestone 02), the CLI is a second, worse way to create accounts: it needs SSH to the Raspberry Pi, and it stores passwords lowercased and cut to their first word (tech review #4, `docs/tech-review-2026-11.md`). Keeping it means two account-creation paths with different rules, plus a dependency (`promptui`) used by nothing else.

Assumed impact: Tiimit needs no terminal tool for anything an organizer does, and there is one place where accounts are created.

## Decided
<!-- Choices Human-PM made at kickoff, written down so they live
     in the spec rather than in one session's memory. -->
- **No terminal at all.** The CLI (`cmd/cli`) is removed, with its dependency (`promptui`).
- **The first user of an empty database** is created from two settings in the server's environment, `FIRST_USER_USERNAME` and `FIRST_USER_PASSWORD` (Q4), read on startup only while there are no users. The live database already has users, so this is for a fresh install, a developer machine, or a recovery onto an empty database. A setup page shown while there are no users was rejected: the app is on the public web, so whoever reached it first would become the first user (02, kickoff).
- **Users are counted before the settings are looked at** (R1). On a database with users, the settings are ignored entirely, even when they break the rules, so a leftover setting can never stop the live app on a deploy.
- **The first user follows milestone 02's username and password rules**, from the same code, so there is still only one set of rules. The username is normalized exactly as at sign-up; the password is used exactly as set, with nothing trimmed (R5).
- **Bad or half-set settings stop Tiimit from starting** (Q2, R3), which can only happen on an empty database. The message names the setting to fix and never shows the password (R5).
- **The first-user password in plain text is an accepted risk** (Q3, carried over from 02/Q3): it only acts on an empty database, which is rare. The deployment guide says to remove it after a recovery (R9).
- No migration: this milestone changes no table.

## Scope
- **In:**
  - Creating the first user from the two settings on startup. In `.env.example` the two settings are present but commented out, so a copied `.env` starts as in S3 (R9).
  - A short "Recovering onto an empty database" section in `docs/deployment-guide.md`: where to add the two settings on the Pi, and to remove them after the first login (R9).
  - Removing `cmd/cli`, `promptui` from `go.mod`, and the CLI-only `auth.PasswordMinLength` and `auth.IsWeakPassword`, which contradict the real 10-character rule (R7).
  - The test helper `createCoach` (`cmd/web/helpers_test.go`) creates the coach through the new first-user code instead of copying the CLI (red-team, optional).
  - Removing or rewording every mention of this app's CLI (R8): `CLAUDE.md` ("server/CLI startup", "plus a small CLI"), `README.md`, `ARCHITECTURE.md` (the diagram and the other CLI lines), `spec/domain-model.md` ("Users are created only with the CLI"), `spec/user-flows.md`, `bruno/README.md` and `bruno/auth/login.bru`. Mentions of the Tailwind or templ CLI stay. `deploy.sh` has no CLI step.
  - Marking I17 done in `spec/product-notes.md`.
- **Out:**
  - Changing the first user's password later (I26).
  - Roles: the first user is not an admin, since nobody is (I21, I27).
  - Cleaning up other unused settings in `.env.example` (tech review #16).
  - Guarding against two copies of Tiimit starting at the same moment (R10): only one runs on the Pi, and the unique username already makes a race create one user at most.

## Acceptance conditions (BDD)

### S1: Developer starts Tiimit on an empty database with first-user settings
- Given a database with no users
- And the first-user username and password are set and follow the rules
- When the developer starts Tiimit on the server
- Then Tiimit is running
- And the first user can log in with that username and password

### S2: Developer starts Tiimit on an empty database with first-user settings that break the rules
- Given a database with no users
- And the first-user username or password breaks the sign-up rules
- When the developer starts Tiimit on the server
- Then Tiimit stops with a message
- And the message says which setting to fix and why
- And the message does not show the password
- And no user is created

### S3: Developer starts Tiimit on an empty database without first-user settings
- Given a database with no users
- And neither first-user setting is set
- When the developer starts Tiimit on the server
- Then Tiimit is running
- And the log says that nobody can log in until both first-user settings are set

### S4: Developer starts Tiimit on a database that already has users
- Given a database with one user
- And the first-user settings are set, with a password that breaks the rules
- When the developer starts Tiimit on the server
- Then Tiimit is running
- And no user is created
- And the existing user can still log in with their own password

### S5: Developer starts Tiimit on an empty database with only one first-user setting
- Given a database with no users
- And the first-user username is set
- And the first-user password is not set
- When the developer starts Tiimit on the server
- Then Tiimit stops with a message
- And the message says the first-user password is missing
- And no user is created

## Test notes
- All five drive one startup function (in `cmd/web`) against a fresh database, not the whole `main`. The function returns an error for S2 and S5 and a "nobody can log in" result for S3; `main` does the logging and stopping (R6). The tests check the returned values, not log output.
- S1 and S4 then log in through the real `POST /login` with the milestone 00 harness.
- S2: two cases only, one username rule and one password rule, to prove that a rules error stops startup and names the right setting. The rules themselves are already unit-tested in `internal/auth/credentials_test.go` (R4). Each case checks that the error text does not contain the password.
- S1 also checks normalization: a username set with capitals and surrounding spaces logs in in its normalized form (R5).
- The count and the create run in one transaction (R10). No test for the race.
- Removing the CLI has no scenario: `go build ./...`, `go vet ./...` and `make test` pass without `cmd/cli`, and `go mod tidy` drops `promptui`. The mention check is `grep -rn 'cmd/cli\|promptui\|the CLI\|IsWeakPassword\|\<PasswordMinLength' . --exclude-dir=.git --exclude-dir=bin`. Its expected hits: `spec/legacy/`, finished milestone specs (00 to 02), this spec, `docs/tech-review-2026-11.md`, the I17 and I26 lines of `spec/product-notes.md` (R8), and the success measure in `spec/user-flows.md` (see the decisions log).

## Open questions
<!-- For the red-team pass (workflow step 2) and Human-PM to resolve (step 3).
     A resolved question keeps its number and moves to the decisions log. -->
None open.

## Decisions log
<!-- Red-team findings (R1, R2, …) with Human-PM's decision on each, resolved
     questions (Q…), and decisions made while building. One line of why each. -->
- S1 is 02/S10, moved here when milestone 02 was split (02/Q5).
- **Q1** No log reminder to remove the password after creating the first user: more code to maintain for little gain. R9's deployment-guide note covers it instead.
- **Q2** Bad first-user settings stop Tiimit from starting, rather than starting and creating nobody. Accepted because, with R1, it can only happen on an empty database, where nobody can log in anyway.
- **Q3** The first-user password stays in plain text in the server's environment file. Human-PM accepted the risk: the setting only acts on an empty database, a rare case.
- **Q4** The settings are named `FIRST_USER_USERNAME` and `FIRST_USER_PASSWORD`.
- **R1** (high) Accepted: count users before looking at the settings; new S4 covers a database with users and bad settings. Why: otherwise a leftover setting could stop the live app on a deploy.
- **R2** Accepted: S1's "later start" line was a second situation; it is now S4.
- **R3** Accepted: only one setting set stops startup like S2 (new S5). Why: it is almost always a typo.
- **R4** Accepted: S2 tests one username case and one password case, not all of 02's rules again.
- **R5** Accepted: username normalized as at sign-up, password used exactly as set, the message never shows the password.
- **R6** Accepted: the startup function returns values; `main` logs and stops. Why: `log.Fatalf` cannot be tested.
- **R7** Accepted: remove `auth.PasswordMinLength` and `auth.IsWeakPassword`, used only by the CLI and contradicting the 10-character rule.
- **R8** Accepted: corrected the list of files to update (Tailwind CLI mentions stay) and made the grep check precise.
- **R9** Accepted: recovery section in the deployment guide; settings commented out in `.env.example`. Unverified claim to check on the Pi: whether the systemd unit file is readable by all users (`ls -l /etc/systemd/system/tiimit.service`).
- **R10** Noted, no design: one transaction for the count and the create; a race fails safe on the unique username.
- Red-team optional finding accepted: `createCoach` uses the new first-user code.
- While building (slice 1): 02/S6's test created another user before the coach; its two setup lines were swapped, since the coach must now be the first user. What it checks is unchanged.
- While building (slice 2): `spec/user-flows.md` keeps "without SSH or the CLI". It is a success measure in Human-PM's file and is still true; agents edit that file only to keep references in sync.
- While building (slice 2): the recovery section uses a systemd drop-in (`systemctl edit`), so removing the password afterwards is one command (`systemctl revert`).

## Slices
1. **First user from settings.** Startup counts users, then creates the first user, refuses, or reports that nobody can log in. `createCoach` switches to it. Covers S1 to S5.
   - Human-PM checklist:
     1. With an empty `DB_PATH` and both settings set in `.env`, run `make run` and log in with them.
     2. Stop, set the password to `short`, run again: Tiimit stops and the message names `FIRST_USER_PASSWORD` without showing it.
     3. Comment out both settings, run again on the same database (it now has a user): Tiimit starts and the user still logs in.
2. **Remove the CLI.** Delete `cmd/cli`, `promptui`, `PasswordMinLength` and `IsWeakPassword`; update the docs listed in Scope; add the recovery section; mark I17 done. No scenario; see Test notes.
   - Human-PM checklist:
     1. `make test` passes.
     2. Read the new recovery section in `docs/deployment-guide.md`: could you follow it on the Pi without help?
     3. `README.md` and `CLAUDE.md` describe Tiimit as a web app only.
