# Managing players

The coach needs to keep the club's roster up to date so that attendance imports and team distribution work. Today the only way to add a player is to export the full roster from MyClub, SSH into the Raspberry Pi, and run the CLI import. That's cumbersome enough that adding a mid-season joiner gets delayed, and until they're added they can't be distributed into a practice.

## Objective

Let the coach manage the roster from the web app, without the CLI or SSH. Build it in thin slices so each one is usable and deployable on its own.

### Viewing players

- [x] Display all players of the club to logged-in users.
- [x] Show each player's scores (run power, ball handling, total).

### Adding a player (Slice 1 — name + MyClub ID)

A coach can add a single new player from the Players page. The player is immediately part of the roster and matchable in the next attendance import.

- [x] Add a form on the Players page to create a player (name + MyClub ID).
- [x] Persist the new player and show them in the list without a full page reload.
- [x] Give feedback about what went wrong when trying to submit incorrect data or backend fails.

**Expected behaviour**
1. You can add a player by entering a name and a MyClub ID.
2. The new player appears in the players list immediately after adding.
3. When you distribute teams by using the attendance report, new players are distributed and show in the practice.
4. If you enter a MyClub ID that already exists, you get clear feedback and no duplicate is created.
5. If the name or MyClub ID is empty, you get feedback and no player is created.

**Open decision — starting scores**
Scores are required by the schema, but this slice doesn't collect them. A new player needs a default:
- Neutral (run power 5, ball handling 5) so an unknown player doesn't systematically weaken one team. *(Recommended.)*
- Zero (0/0) so it's obvious the player is unrated — but they'll always sort last and skew distribution.

Decision: New players use neutral default scores (5/5). Consider allowing custom scores based on feedback and experience.

### Out of scope for now (later slices)

- Editing a player's scores from the web (still done via CLI import today).
- Setting goalie status at creation or from the web (CLI `ToggleGoalieStatus` still works; new players default to not-goalie).
- Where the coach finds the MyClub ID — assumed they look it up in MyClub. Revisit if this proves to be friction.
