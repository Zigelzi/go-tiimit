# Improvements

Possible improvement ideas categorized based on the user flow.

## Finding the practice

- [ ] Add number of players to the practice "card".

## Setting up practice

- [x] Move player to other team to balance teams.
- [x] View total score of the team to balance teams.
- [ ] View unknown players that weren't possible to distribute.
- [ ] Display players who weren't possible to distribute to practices (e.g doesn't exist in club)
- [ ] Distribute the players who haven't signed up, but showed up to the practice (unknown players).
- [ ] Create and distribute a guest/tryout player.
- [ ] Notify (other) coaches about new practice via email.

## Managing players of the club

- [x] Import players from MyClub via web UI.
- [x] Create new player to club.
- [ ] Mark player as goalie or not.
- [x] Adjust the scores of the player.
- [ ] Archive and unarchive (soft delete) players.

## Sign up and registering

- [ ] Reset password.
- [ ] Sign up with when initiated / verified by admin.

## Technical

- [ ] Deprecate CLI
  1.  Need to figure out how to create users safely in this case.
- [ ] Distinguish "database unavailable" from "not logged in".
  `authMiddleware` treats every `GetActiveSession` error as not-logged-in, so a
  DB failure silently redirects everyone to `/login` instead of showing an
  error. Users see "you're logged out" when the truth is "the app is broken".
  Separate `sql.ErrNoRows` (genuinely no session) from other errors (system
  failure). See `internal-docs/handler-error-handling.md` §7.

## Error handling cleanups

Found while writing `internal-docs/handler-error-handling.md` — worth fixing
during the `renderError` rollout. Section references point back to that doc.

- [ ] `handlers_players.go:255-298` — the same message,
  "Couldn't set player to inactive", is sent from six sites, including both
  branches that set a player **active**. Branch the message per direction. (§3)
- [ ] `handlers_auth.go:82` — logs the raw session token
  (`sessionCookie.Value`). Anyone with log access can replay that session. Log
  the user id or a prefix instead. (§2)
- [ ] `handlers_players.go:85` — malformed format verb: `%` followed by a
  newline instead of `%v`. The line prints
  `MyClubId isn't a number: &{%!\n(string=...)}`. Note `go vet` does **not**
  catch this (verified: `go vet -printf ./cmd/web/` is clean) — it catches
  wrong-type and wrong-count verbs, not every garbled format string. (§2)
- [ ] `handlers_players.go:25` and `:33` — two different queries log the
  identical line `failed to get players: %v`, so a grep hit can't tell you which
  one broke. Their *user* messages already distinguish active from inactive —
  it's only the logs that are ambiguous. (§2, §4)
- [ ] `handlers_players.go:25` and `:33` — user messages open with "Failed",
  which reads as system-speak. Prefer "Couldn't load the … right now." (§3)
- [ ] `handlers_practice.go:319` and `:364` — `"failed toget"` typo. Harmless,
  but it breaks a grep for `failed to get`. (§2)
- [ ] `handlers_practice.go:75` and `:102` — bare `log.Println(err)` with no
  context: no operation, no ids. (§2)
