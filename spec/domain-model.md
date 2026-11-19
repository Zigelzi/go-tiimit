# Domain model

The entities and how each is stored. The schema itself is the goose migrations in `sql/schema/`; this file says what the tables mean. Code layout is in [`../ARCHITECTURE.md`](../ARCHITECTURE.md).

## The team in real life

Tiimit serves one team: VJS Kuntofutis, a hobby football team in the VJS club. Elsewhere in the specs and code, "the club" means this team. It is implicit: nothing is stored per club or team.

- **People.** One organizer (Human-PM) runs the practices and is the only one with MyClub admin access, so the only one who can download the player and attendance exports. Two assistant organizers, who are players too, hand out the vests. All three are users of Tiimit (see [User and session](#user-and-session)).
- **Seasons.** The year has two similar seasons, winter (Nov–Mar) and summer (Apr–Oct). Tiimit doesn't store seasons and doesn't need to. A season change shows up only as roster turnover: leavers are set inactive, joiners are added by hand.
- **Signing up.** Every player signs up for each practice in MyClub. A player starts at "no answer" and moves to "attends" or "doesn't attend". Only "attends" reaches the teams (see [Practice](#practice)). Some "no answer" players still turn up and need a place on a team that is already split; Tiimit can't do that yet (I4 in [`product-notes.md`](product-notes.md)).

## Player

A member of the club who can be put on a team. Table `players`.

- **Name** and **MyClub ID** (unique): the ID is how attendance rows are matched to players.
- **Scores**: run power and ball handling, whole numbers 0–10 in the UI, stored as `REAL`. A new player added by hand starts at 5/5 so an unknown player doesn't weaken one team (legacy/players.md). The total drives team balancing (`internal/player/scoring.go`).
- **Goalie**: whether the player is willing to play in goal. Each team should get one.
- **Inactive**: `inactivated_at` set means the player has left or paused; inactive players are listed separately on the players page. (Practice creation matches attendees by MyClub ID without checking this, so an inactive player who signs up is still distributed.)

## Practice

One session on one date, split into two teams. Table `practices`; who plays where is `practice_players`.

- A practice is created from an attendance report, dated from the report's file name. Rows marked "attends" whose MyClub ID matches a player become the practice's players. "No answer" and "doesn't attend" rows are ignored, and an "attends" row with no matching player is skipped and only logged (I3).
- The player list is fixed once the practice is created: players can move between teams but none can be added or removed (I4, I5).
- **Team**: each practice player is on team 1 or 2. Exactly two teams is assumed throughout, including the database check.
- **Vest**: whether the coach has handed this player a vest at this practice.
- Nothing can be deleted in the app yet. The schema declares that deleting a practice or player removes their practice rows (`ON DELETE CASCADE`), but SQLite enforces that only with foreign keys switched on, which `internal/db/config.go` doesn't do. A milestone that adds deleting switches them on first.

## User and session

Someone who can log in: the organizer or an assistant organizer. Tables `users` (username, hashed password) and `user_sessions` (cookie session id, expiry 14 days, soft-deleted on logout). Users are created only with the CLI (`cmd/cli`); there is no sign-up.

- **All users have the same rights.** Any logged-in user can add and edit players, create practices, move players and tick vests. Organizer-only rights are not built (I21).
- A user is not linked to a player, even when the person is one.
