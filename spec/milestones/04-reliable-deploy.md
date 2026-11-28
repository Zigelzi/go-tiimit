# Milestone 04: Reliable deploy

**Status:** done · **Branch:** chore/reliable-deploy
<!-- Status: draft → red-teamed → in-progress → done -->

## Why
Every deploy migrates the club's only copy of the data, and nothing takes a copy first. A deploy that breaks the app still ends looking fine, and nothing shows which version is running. As the person deploying, I want to deploy without worrying about losing the data, and to find out right away when a new version does not start. Coaches do not see this milestone. They only notice that their players and practices are still there after every deploy.

It supports every flow in `../user-flows.md` because it protects the data they all use. It does not change any step. It takes on I20, which covers deployment-guide gaps 3 to 8, 11 and 12 (`../../docs/deployment-guide.md#known-gaps--improvement-backlog`).

## Decided
Ported from Roundtrip's milestone 18 (`../roundtrip/spec/milestones/18-pi-deploy.md`, files in `../roundtrip/deploy/`), which runs on the same Pi, and changed where a public app with personal data needs more:

**Deploying**
- `make deploy` replaces `./deploy.sh`. It runs the tests, refuses work that is not committed, builds, uploads, and then runs the one root script on the Pi.
- A deploy refuses uncommitted changes and any branch other than `main`. A committed branch can be deployed on purpose with `ALLOW_BRANCH=1`, so a milestone can be tried on the Pi before it merges, but only when it changes no migration file compared with `main` (R1, R30).
- A new migration must be dated after the newest one on `main`. Goose refuses to start on an older-dated migration it has not run (R21). The rule goes in `ARCHITECTURE.md` at merge.
- The deploy is not left half done when the laptop's connection drops: the root script keeps going, and the guide says to run the deploy again (R22).

**The root script** (`deploy/tiimit-pi`, installed by hand to `/usr/local/sbin`, owned by root)
- It is the only command the deploy user may run with sudo without a password, through a sudoers file kept in the repository (gaps 4, 7). The deploy user keeps its normal sudo with a password, which the rollback steps use (R15, Q6).
- It takes no paths or other input from the caller, installs only a regular file from a fixed upload folder owned by the deploy user, never runs the uploaded program as root, and refuses to run unless the program folder is owned by root and writable by nobody else (R5, R16).
- It stops when the database file is missing, because a missing file means a wrong setting, and the app would otherwise start on a new, empty database (R18).
- It keeps exactly one previous program, not a pile of timestamped ones (gap 4).

**Backups**
- SQLite's own backup, never a file copy. Before every deploy and every Monday at 04:00 Europe/Helsinki; each kind keeps its newest 3 (Q5, Q10). The weekly timer catches up after a power cut.
- Only root can read or delete them: they live in `/var/backups/tiimit`. The app's own account makes the copy and root moves it in, so a break-in through the app cannot delete them (R19, Q11).
- They hold live login sessions, so a restore ends by deleting all sessions and all unused invitations (R2, R23). The restore runs as the app's own account, never as root, so the database keeps its owner (R23).

**After the restart**
- The script checks on the Pi that the home page answers with a normal page, retrying for a few seconds. The home page reads the database even for a visitor who is not logged in, so an answer proves the app and the database both work (Q2, R17).
- The deploy then checks `https://tiumo.fi` from the laptop, as a warning only: it fails when nginx or DNS is broken while the app is fine (R17, Q7).
- When the app does not answer, the deploy says which step failed and prints the rollback steps. Rolling back is not automatic. An older program does start on a newer database, so after a deploy that ran a migration the steps put back the backup too, and they say that this loses anything saved since (R4, R21).
- The app writes the commit it was built from to the Pi's log at startup (gap 11).

**The Pi's setup**
- The systemd unit is kept in the repository, with Roundtrip's hardening: no gaining rights, a read-only system, write access to `/var/lib/tiimit` only (R6). Its working folder stays `/opt/tiimit`, owned by root, so nobody but root can put a `.env` file there (R20).
- The settings live in `/etc/tiimit/tiimit.env`, which only root can read, with an example file in the repository (R4, Q4). The app reads its listen address and its environment from them, defaulting to today's `127.0.0.1:8080` and development (gap 8, R14). The example keeps `127.0.0.1`, so the app is never reachable around nginx (R31).
- Nobody but root can change the installed program; only the app's own account can read its data (R9, R10).
- SSH accepts keys only. The setting goes in its own file that sorts before Ubuntu's `50-cloud-init.conf`, which would otherwise win (R11, R26, Q8).
- The one-time setup clears out the old backup programs, the removed CLI and the folder of attendance reports (Q1, R12).
- The one-time setup follows a fixed order, so it can never leave the site down or Human-PM locked out halfway (R25). It is in the guide; Slices 2 and 3 follow it.

**Other**
- Nothing about this milestone is pushed to GitHub before the setup is done and S8 passes, because the decisions log lists weak points that are open until then (R31).
- Different from Roundtrip: Tiimit is public through nginx and has login. The nginx and firewall setup stay as they are, so Roundtrip's home-network-only parts are not ported.
- `/code-review high`, because the script runs as root on a server that is public.

## Scope
- **In:**
  - `make deploy`, its checks, and the root script and sudoers file in the repository
  - Backups before every deploy and weekly, protected from every account but root
  - The check that the app answers, and the rollback steps printed when it does not
  - The commit written to the Pi's log at startup
  - The systemd unit, the example environment file, and the app reading its address and environment from settings
  - The one-time setup on the Pi, written down in order: file owners and permissions, settings, unit, SSH without passwords, clearing out leftovers
  - The deployment guide rewritten around `make deploy`: one-time setup, deploy, rollback, restore
  - A restore rehearsed once on the Pi (Q9)
- **Out:**
  - Automatic rollback (Roundtrip left it out on purpose: restoring a backup can lose data)
  - Copies of backups off the Pi. Accepted risk: if the SD card fails, the data and every backup are lost (R8, I32)
  - The `Secure` flag on the login cookie and the username in the failed-login log: one small fix right after this milestone (R3, R31, I22)
  - Storing login and invitation tokens hashed, so a backup is no longer a working login (R2, I30)
  - Changes to nginx, TLS or the firewall, including HSTS and a rate limit on login (I31)
  - Quieting the `.env` warning at startup (gap 10)
  - CI, containers, deploying from anywhere but the laptop

## Acceptance conditions (BDD)

Most of these check the Pi and the deploy, which a `go test` cannot see, so they are checked by hand. S7 and S11 are `go test`s. How each one is checked is in Test notes. S4, S9 and S10 were merged into other scenarios at refinement (R28); their numbers stay free.

### S1: A deploy keeps the data and leaves a backup
- Given the Pi holds the club's players and practices
- When a new version is deployed
- Then the same players and practices are there afterwards
- And a backup of the database from just before the deploy is on the Pi
- And the Pi's log names the commit the new version was built from

Verification: manual

### S2: A version that does not start is reported
- Given a new version that does not answer after it starts
- When it is deployed
- Then the deploy ends with a failure message
- And the message names the step that failed
- And the message says how to put the previous version and its backup back

Verification: manual

### S3: Work that is not committed to main is not deployed
- Given the work on the laptop is not committed to `main`
- When the deploy is started without naming a branch
- Then the deploy stops and says why
- And nothing is copied to the Pi

Verification: manual

### S5: The Pi keeps the newest 3 backups of each kind
- Given weekly backups have been taken more than 3 times
- And more than 3 deploys have happened in between
- When the backups on the Pi are listed
- Then there are exactly the newest 3 weekly backups
- And there are exactly the newest 3 deploy backups

Verification: manual

### S6: A backup is restored on the Pi
- Given a backup taken a minute earlier
- And a change made to the data after it
- When the written restore steps are followed
- Then Tiumo shows the data as it was when the backup was taken
- And everyone has to log in again
- And invitations sent before the restore no longer work

Verification: manual

### S7: The app runs with the settings it is given
- Given the settings name a listen address and the production environment
- When the app starts
- Then it listens on that address
- And it runs in production mode

### S8: The Pi's accounts can do only their own job
- Given the deploy user, the app's own account and Roundtrip's account on the Pi
- When each one tries to do more than its job
- Then the deploy user can run only the deploy script with root rights without a password
- And only the app's own account can read the club's data
- And only root can change the installed app or read the backups
- And nobody can log in to the Pi over SSH with a password

Verification: manual

### S11: The app runs as it does today without settings
- Given no listen address and no environment in the settings
- When the app starts
- Then it listens on this computer only
- And it runs in development mode

### S12: A branch that changes the database is not deployed
- Given a committed branch that changes a migration file compared with `main`
- When the deploy is started with that branch named on purpose
- Then the deploy stops and says why
- And nothing is copied to the Pi

Verification: manual

### S13: The Pi takes a weekly backup on its own
- Given the Pi holds the club's players and practices
- When Monday 04:00 passes
- Then a new weekly backup is on the Pi
- And it holds the players and practices as they were at that time
- And nobody had to start it

Verification: manual

### S14: A deploy stops when the Pi's database cannot be found
- Given the settings on the Pi point to a database that does not exist
- When a new version is deployed
- Then the deploy stops and says the database is missing
- And the running version keeps working
- And no new, empty database is created

Verification: manual

## Test notes
- S1: count players and practices on the phone before and after. `sudo ls /var/backups/tiimit/deploy` on the Pi shows the new backup. `sudo journalctl -u tiimit -n 20` names the same commit as `git log -1` on the laptop. Pruning is checked in S5.
- S2: give the listen address setting a value the app cannot use, so it exits at startup, and deploy. Not the database path: a wrong path would start the app on an empty database (R18). Undo it afterwards. The deploy exits with an error, and the site is back after the printed steps are followed.
- S3: make an uncommitted change and deploy; then switch to a branch and deploy without `ALLOW_BRANCH`. Both stop before anything is copied. A branch with `ALLOW_BRANCH=1` and no migration change goes through.
- S5: run the weekly backup job by hand 4 times and deploy 4 times, alternating. `sudo ls /var/backups/tiimit/weekly /var/backups/tiimit/deploy` shows 3 in each, the newest ones.
- S6: on the Pi, at a quiet time (Q9). Take a deploy, change something small (tick a vest), create an invitation, follow the restore steps, and check the tick is gone, the phone is on the login page and the invitation page says it is not valid. Nothing but that tick and everyone's login is lost. A backup that has never been restored is not known to work.
- S7, S11: `go test` on how the app reads its settings. The production cache headers on a static file are also seen once on the Pi.
- S8: as the deploy user, `sudo -n -l` lists only the deploy script and `sudo -n true` is refused; `sudo -v` with the password still works (R15). Reading the database and writing to `/opt/tiimit` fail. As the app's account (`sudo -u tiimit`), renaming the installed program and listing `/var/backups/tiimit` fail. As Roundtrip's account, reading the database fails. `ssh -o PubkeyAuthentication=no` is refused without asking for a password, and `sudo sshd -T` shows password login off. `systemd-analyze security tiimit` scores better than before, and saving a practice still works with the hardening on (R27).
- S12: on a branch with an empty new migration file, deploy with `ALLOW_BRANCH=1`. It stops before anything is copied.
- S13: before acceptance, `systemctl list-timers` shows the backup timer's next run on Monday 04:00 Europe/Helsinki, and one run started by hand leaves a backup that opens read-only as root with the same number of players and practices as the live database. The first real Monday backup is checked after the merge (Human-PM, 2026-11-22). Catching up after a power cut is checked in code review, by the timer's catch-up setting, not by switching the Pi off.
- S14: in slice 3, change the database path in `/etc/tiimit/tiimit.env` to a file that does not exist, deploy, and check the phone still shows the practices and no new database file appeared. Undo the change.

## Open questions
- None open. Q1 to Q11 are resolved in the decisions log.

## Decisions log

Questions resolved by Human-PM (2026-11-21):
- **Q1 One-time setup.** Human-PM picks the time, since they set up the practices. The old backup binaries are deleted, together with the leftovers from R12.
- **Q2 The "does it answer" check** asks for the home page at the public address: closest to what users experience, and the home page reads the database, so if it does not answer nothing does.
- **Q3 Restore rehearsal off the live data.** Refined so no real data leaves the Pi: the full steps on the laptop with made-up names, and a read-only check of a real backup on the Pi (see S6 test notes).
- **Q4 The environment file** is `/etc/tiimit/tiimit.env`, readable by root only. Human-PM sets it up from the example file and the setup steps. Today's unit holds `DB_PATH`, `POSTHOG_API_KEY` and `ENV`, with no first-user password. The PostHog key reaches the browser anyway, so it is not a secret, but the file stays root-only for whatever is added later.
- **Q5 Weekly backup** on Monday at 04:00 Europe/Helsinki.

Review by the main session (2026-11-21), with the Pi's current state read by Human-PM. Not independent, since the main session wrote the draft; an independent red-team pass follows.
- **R1 Accepted. Branch deploys only without new migrations.** A migration from a branch stays in the live database even if the branch changes or is discarded, and deploying `main` again does not undo it.
- **R2 Accepted. Backups are a working login.** Login sessions and invitation tokens are stored as plain text, so a backup lets its holder log in, and restoring one brings back logged-out sessions and used invitations. Backups are readable by `tiimit` only, and the restore steps end by deleting all sessions and checking open invitations. Hashing the tokens is I30.
- **R3 Deferred.** The login cookie has no `Secure` flag, so a first plain-http request sends it unencrypted before nginx redirects. Fixed right after this milestone as its own small change (I22).
- **R4 Accepted. Settings in a root-only environment file** instead of `Environment=` lines that every account on the Pi can read.
- **R5 Accepted. Rules for the root script.** It takes no paths or other input from the caller, installs from a fixed upload folder owned by the deploy user (not `/tmp`), refuses anything but a regular file, and never runs the uploaded program as root. Backups and pruning run as `tiimit`, so root never deletes files in a folder another account can write to.
- **R6 Accepted. Port Roundtrip's systemd hardening** to the app and the weekly backup job: no gaining rights, a read-only system, write access to `/var/lib/tiimit` only. It limits what someone who breaks in through the public side can touch.
- **R7 Settle with a test at refinement.** Unknown whether an older binary starts against a database with newer migrations. The answer decides whether the rollback steps put back the binary alone or the binary and its backup.
- **R8 Accepted risks.** A leaked deploy key lets someone replace the app and so reach the data, as with any deploy user. All backups stay on the SD card (I32).
- **R9 Accepted. The app's folder is owned by `tiimit` and the program by the deploy user today**, not by root as the guide says. Either can replace the app, the app's own account included, so a break-in through the public side could survive a restart. Setup sets both to root (S9).
- **R10 Accepted. The database can be read by every account on the Pi today** (folder `755`, file `644`), Roundtrip's account included. Setup makes the folder `700` and the file `600`, and the backup folders the same (S9).
- **R11 Accepted. Password login over SSH is on today**, and the firewall allows SSH from anywhere. Turned off whatever the router forwards, after checking the key login works; the deploy user has a key and root has none (S10). Which ports the router forwards was not checked.
- **R12 Accepted. Leftovers in the app's folder:** the CLI removed in milestone 03, and a folder of 3 attendance reports, readable by every account, with real names in them. Both are deleted at setup.
- **R13 Confirmed.** Today's sudo rule allows running as any user, a move from `/tmp` and a wildcard copy (gap 4). Replaced by the one-script rule.
- **R14 Confirmed.** The unit sets `ENV=production`, but the app ignores it and runs in development mode (gap 8, S7).

Facts read from the Pi (2026-11-21): Ubuntu 24.04 on `aarch64`, 20 GB free, `sqlite3` 3.45.1, time zone Europe/Helsinki, automatic updates on, no global IPv6 address. The database is 90 KB in journal mode `delete`, not WAL: there are no side files to handle, but a restore removes a `tiimit.db-journal` file if one exists. nginx proxies `tiumo.fi` to `localhost:8080` and sends no security headers.


Red-team findings (2026-11-22, an independent sub-agent on the revised spec), with Human-PM's decisions. All accepted unless noted.
- **R15 Blocking. The new sudo rule could lock Human-PM out of root.** Root has no SSH key, so a broken sudo means taking out the SD card. Setup checks `sudo -v` with the password first, installs the rule only with `visudo`, and keeps a root shell open in a second terminal throughout. The deploy user is in the `sudo` group and Human-PM knows its password (Q6).
- **R16 Blocking. The root script is only safe once `/opt/tiimit` belongs to root.** Before that, either account could plant a link there and have root overwrite a system file. Ownership changes first in the setup, and the script refuses to run otherwise.
- **R17 Blocking. The check against the public address can be wrong both ways:** a redirect from nginx passes while the app is down, and the laptop may not reach the public address from home. The main check runs on the Pi and needs a normal page; the public check is a warning only. Changes Q2.
- **R18 A wrong database path starts the app on an empty database, and the check passes.** The script stops when the database file is missing; S2 breaks the listen address instead.
- **R19 Backups where the app's account could delete them.** Protected in a root-only folder (Q11).
- **R20 The app reads a `.env` file from its working folder**, which would let it add settings to itself. The working folder stays root-owned `/opt/tiimit`; no `.env` is there today (Q8).
- **R21 Settled from the Goose source (v3.26.0, `up.go`), replacing R7.** An older program starts on a newer database without complaint, so a rollback after a migration also puts back the backup. Goose refuses to start on a migration dated before one already run, hence the dating rule.
- **R22 An interrupted deploy** could leave an untested program installed and then copy it over the previous one on the next deploy. The script keeps going when the connection drops.
- **R23 The restore runs as the app's account**, so the database keeps its owner, and deletes sessions and unused invitations in the same step.
- **R24 A laptop rehearsal skips the steps most likely to fail**, the Pi's account, unit and permissions. Rehearsed on the Pi instead (Q9, changes Q3).
- **R25 The one-time setup has a fixed order:** root shell in a second terminal; manual backup; owners and permissions; environment file; unit, restart, check the phone; script and sudo rule; first `make deploy`; SSH passwords off; leftovers deleted.
- **R26 Ubuntu's `50-cloud-init.conf` would keep password login on.** The setting goes in a file that sorts first, checked with `sshd -T`. The file exists on the Pi (Q8).
- **R27 A check for the hardening**, in S8's test notes.
- **R28 Scenario rules.** S3 and S7 split, S4 folded into S1, S8 to S10 merged into S8.

Review of the conditions by Human-PM (2026-11-22):
- S13 added by Human-PM: nothing proved that the weekly backup happens on its own, only that 3 are kept. Its first real Monday is checked after the merge, so acceptance does not wait up to a week.
- S5 now covers the deploy backups' pruning too, which had dropped out of S1 at refinement.
- S6 also checks that invitations stop working after a restore (R23).
- S14 added: R18's "stop when the database is missing" is a data-safety rule, so it gets a scenario.
- Catching up after a power cut stays a Decided item checked in code review: switching the Pi off across 04:00 is fiddly and adds little.
- The public address warning (R17) gets no scenario: it only warns, and breaking nginx to test it is not worth the risk.
- Eleven scenarios, two of them `go test`s. Above the rough ceiling of 8, but most are short checks on one setup, and splitting would mean setting the Pi up twice.
- **R29 Keep more backups: declined.** 3 and 3 is enough for now (Q10).
- **R30 The branch check refuses any change to a migration file**, not only new ones.
- **R31 What goes public.** Push only after the setup is done; the example settings keep `127.0.0.1`. The username written to the log on a failed login (`cmd/web/handlers_auth.go:30`) is fixed together with I22.

Questions resolved by Human-PM (2026-11-22):
- **Q6** The deploy user is in the `sudo` and `adm` groups, and Human-PM knows its password.
- **Q7** `https://tiumo.fi` loads from the laptop on home Wi-Fi, so the public warning check works from there.
- **Q8** No `.env` in `/opt/tiimit`. `/etc/ssh/sshd_config.d/` holds only `50-cloud-init.conf`.
- **Q9** Rehearse the restore on the Pi with a backup taken a minute earlier.
- **Q10** Keep 3 deploy and 3 weekly backups.
- **Q11** Protect the backups from every account but root. In practice: they live in `/var/backups/tiimit`; the app's account makes the copy and root moves it in; the weekly job runs as root through the same script; a restore has one extra step where root hands the backup to the app's account; looking at a backup needs `sudo`.

Decisions made while building:
- Slice 2: the unit also sets `UMask=0077`, so files the app creates itself, such as SQLite's temporary journal, are readable by `tiimit` only, like the database. One line, and it keeps R10 true for files created later.
- Slice 2: the guide's setup step 3 first checks that the program and data folders hold no links, since root changes owners there (R16).
- Slice 2, accepted 2026-11-28: the security score went from 9.2 to 8.3. More hardening options were offered and skipped for now (I33).
- Slice 3: root copies a backup into the root-only folder by reading it through the app's account, as it reads the upload through the deploy user, instead of moving it (Q11). A move could carry a link the app planted in its place; a copy carries only the bytes. The app's account checks the copy with SQLite's `integrity_check` first, and a backup never overwrites an older one.
- Slice 3: the script marks a version that passed its check (`web.checked`). Only a checked version becomes `web.previous`, and after an unchecked one the deploy backups are not pruned (R22). Replaces Roundtrip's "is the service running" guess, which passes an app that runs but does not answer.
- Slice 3: the script reads `DB_PATH` and `LISTEN_ADDRESS` from the root-only settings file, so S14 tests the path the app really uses, and the check asks for the address the app really listens on. It ignores hang-ups and writes every message to the Pi's log too (`journalctl -t tiimit-pi`), so a deploy that finishes after the connection drops leaves a record (R22).
- Slice 3: the script's refusals and its keep-the-last-good-version logic were tested on the laptop against a copy with fake system commands (no root): missing or unset database, folder writable by others, an upload that is missing, a link or a folder, a corrupt database, extra arguments, a broken version and the deploy after it, pruning to 3.
- Slice 3: recovering onto an empty database now puts `FIRST_USER_*` in the root-only settings file instead of a drop-in file every account can read.
- Slice 3, accepted 2026-11-28: S1, S2, S3, S12, S14 and all of S8 passed on the Pi, setup steps 6 to 9 included. The key login was checked from the laptop before passwords went off.
- Slice 4: the weekly backup is `tiimit-pi backup weekly`, started by `tiimit-backup.timer` through a root service that may write only the data and backup folders and has no network (R6). The sudo rule does not allow it. The setup gains step 10 for it, after the leftovers, so R25's order is unchanged up to there.
- Slice 4: the script takes a lock in the backup folder, waiting up to a minute, so a deploy and the weekly backup never run at the same time: they share the temporary backup file in the data folder.
- Slice 4: the restore uses SQLite's `.restore` into the existing database file, as the app's account (R23), so the owner and the journal are handled by SQLite. Found while testing: `.restore` from a file that does not exist restores an empty database, so the source is opened read-only (`?mode=ro`), which makes a missing file an error. The steps first keep a copy of the data as it is (`before-restore.db`), so restoring the wrong backup can be undone. Rehearsed on the laptop with made-up data built from the real migrations before the rehearsal on the Pi.
- Slice 4 on the Pi, 2026-11-28: S13 and S5 passed (the timer's next run Monday 04:00 EET, a backup by hand with the same counts as the live database, 3 of each kind after 4 rounds). The weekly job's messages showed twice in the log, since systemd also logs what the script prints; the service now drops the printed copy. The S6 rehearsal is moved to Human-PM's next session, after the code review, so it runs the steps as they will merge.
- Code review (high, one round, 2026-11-28): 8 findings. Fixed: a half-written backup counted among the newest 3 (the only important one), the check mark removed before the upload was copied, a failed restart without the rollback steps, the weekly backup's refusals saying "nothing was installed", and `ARCHITECTURE.md`. Skipped: quoted values in the settings file (the example says not to quote, and the failure is a false alarm, never a false pass); `ENV` also being a shell variable (none of the machines sets it, and renaming means editing the Pi's settings); the S7 test not covering `main.go`'s wiring (the Pi showed it in S2 and in every deploy's check). Each fix has a fake-command test.
- Done, 2026-11-28: Human-PM installed the fixed script and closed the milestone without the S6 rehearsal on the Pi. The restore steps passed on the laptop against made-up data built from the real migrations, but have not run against the Pi's account, unit and permissions, the parts R24 calls most likely to fail. The rehearsal on the Pi is I34.

## Slices
Each slice is one commit and ends with Human-PM's test checklist before the next starts. Slices 2 and 3 follow the setup order of R25.

1. **The app reads its settings and names its version.** Listen address and environment from settings with today's defaults; the commit written to the log at startup; the unit and example environment file in the repository. Code only; nothing changes on the Pi yet. Covers S7, S11, and the log line of S1.
2. **The Pi set up safely.** Guide steps for the first half of the setup: root shell, manual backup, owners and permissions, environment file, hardened unit. Human-PM runs them; the app runs from the new unit. Do not deploy with the old `./deploy.sh` after this slice: it would make the program the deploy user's again. Covers the file and hardening parts of S8.
3. **Safe deploys.** The root script (protected backup, install, check, rollback steps), the sudoers file, `make deploy` with its checks, the guide's deploy and rollback steps, then the rest of the setup: sudo rule, first `make deploy`, SSH passwords off, leftovers deleted. Covers S1, S2, S3, S12, S14, and the rest of S8.
4. **Weekly backups and the restore.** The weekly timer, pruning, the restore steps, and the rehearsal on the Pi. Covers S5, S6, S13.
