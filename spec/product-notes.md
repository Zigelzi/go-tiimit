# Product notes

Reference for agents: dated observations from use, and the ideas list. Agents keep it up to date. The headings are the flows in [`user-flows.md`](user-flows.md).

## Observations

- 2026-11-17: Switched to the agentic workflow. Everything shipped before this date is described in [`legacy/`](legacy/); this file's ideas were extracted from it.

## Ideas

Numbers are stable: a done or dropped idea keeps its number and is marked, never removed or renumbered. A milestone that takes on an idea names it in its Why section and marks it here when it is done.

### Team player management

- **I1** Where the coach finds a player's MyClub ID when adding one by hand. Assumed they look it up in MyClub; revisit if it proves to be friction. *(legacy/players.md)*
- **I2** ~~Mark a player as goalie or not; archive and unarchive players.~~ Done before the switch to this workflow: goalie editing and setting players inactive or active again. *(legacy/improvements.md)*

### Setting up the practice

- **I3** Show the attendees who couldn't be distributed (for example, not in the club's roster).
- **I5** Create and distribute a guest or tryout player.
- **I6** Notify the other coaches by email about a new practice.
- **I7** Moving a player updates both teams without a page reload and keeps the coach on the tab they were looking at; both teams' scores, player counts and vest counts are right after the move. *(legacy/practices.md, slice 4, not built)*
- **I8** Feedback when a player moves to the other, hidden tab: ship without it and see, highlight the other tab, or show who moved where. *(legacy/practices.md, open decision; recommendation there: ship without)*
- **I9** Show the score difference between the teams.
- **I10** Re-roll the distribution from the practice page.

### Running the practice

- **I4** Distribute players who showed up without signing up (Running the practice, step 3).
- **I11** Show the number of players on the practice card.
- **I12** A vest tick that fails to save is visible on screen. Open decision: rely on the mismatch with the header count (recommended there), or untick the box automatically. *(legacy/practices.md, slice 2)*
- **I13** Label the team tabs by vest colour, by team number, or both. Check how the coach actually refers to the teams during a session. *(legacy/practices.md, open decision)*
- **I14** Keep two devices in sync when two coaches hand out vests at once. Only worth it if that starts happening.
- **I24** Inactive players are still distributed into teams. Found before milestone 00; left alone there because it changes what a coach sees.

### User management

- **I15** Reset password.
- **I16** Sign-up that an admin initiates or verifies.
- **I17** Deprecate the CLI; needs a safe way to create users first.
- **I21** Organizer-only rights: the organizer (an admin user) manages users and player details; assistant organizers (basic users) can view practices and tick vests, nothing else.

### Technical

- **I18** Tell "database unavailable" apart from "not logged in": the auth middleware treats every session lookup error as logged out. *(legacy/improvements.md)*
- **I22** The login cookie lacks the `Secure` flag and sets no explicit `SameSite` (`cmd/web/handlers_auth.go`). *(00/R12)*
- **I23** Forms have no CSRF protection. *(00/R12)*
- **I20** Deploy like Roundtrip. Roundtrip's milestone 18 (`../roundtrip/spec/milestones/18-pi-deploy.md`, files in `../roundtrip/deploy/`) solved most open items in [`../docs/deployment-guide.md`](../docs/deployment-guide.md#known-gaps--improvement-backlog), on the same Pi. Port it rather than redesign it:
  - A database snapshot with SQLite's `.backup` before every deploy, plus a weekly systemd timer; each keeps its newest 3. This is the most valuable part, because every deploy migrates the club's only copy of the data. (gaps 6, 12)
  - One root script kept in the repo (`roundtrip-pi`: snapshot, install, restart) as the only command the deploy user may run with sudo. That removes the wildcard sudo rule and puts the sudoers file in the repo. (gaps 4, 7)
  - Install with `install -o root`, keep one `.previous` binary instead of timestamped backups that pile up. (gaps 3, 4)
  - After the restart, check that the app answers. If it doesn't, say so and print the rollback steps. (gap 5, without the automatic rollback, which Roundtrip left out on purpose)
  - Stamp the build with its commit and log it at startup. Deploy only from a clean `main`. (gap 11)
  - The systemd unit and an env file in the repo, with the address and environment read from it instead of hardcoded in `main.go`. (gap 8)
  - `make deploy` instead of `./deploy.sh`, with the steps a person runs in one doc.

  Different from Roundtrip: Tiimit is public through nginx with login, so the firewall and nginx parts stay as they are. Rehearse a restore once before relying on the snapshots. Use `/code-review high`: this runs as root on a public server.
