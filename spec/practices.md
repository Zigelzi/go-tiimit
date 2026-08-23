# Displaying teams in the practice in their own tabs

Handing out vests is a one-handed job. The coach is standing on the pitch with a bag of vests, reading names off their phone and ticking each player as they give one out. Today the practice page stacks both teams in one long list, so the team's vest count scrolls out of view exactly when it's needed — the coach has to scroll back up to answer the only question they actually have: *how many are left in the bag?*

Showing both teams at once also means the coach scrolls past eight names that aren't relevant to what they're doing right now. They hand out one team's vests, then the other's.

## Objective

Put each team on its own tab and keep that team's details pinned to the top of the screen, so the coach can work through one team at a time and always see how far along they are. Build it in thin slices so each one is usable and deployable on its own.

### Team details in a sticky header (Slice 1)

The practice date, the team's total score and the team's vest progress stay visible while the coach scrolls the player list.

- [x] Keep the practice date and total number of players visible while scrolling.
- [x] Show the selected team's total score in its own component
- [x] Show the selected team's vest count in the header

**Expected behaviour**
2. The vest count goes up immediately when you tick a player and down when you untick one.
3. The total score shown belongs to the team you are currently looking at.

### A vest count you can trust (Slice 2)

The counter is the reason the header exists, so it has to be right even when the coach is tapping quickly. This slice is about correctness rather than new UI.

- [x] Toggling a vest is a single atomic database statement rather than a read followed by a write.
- [x] Ticking several players in quick succession never unticks a player who was already ticked.
- [x] The count always reflects the most recent tick, not an earlier one that happened to arrive late.
- [ ] A tick that fails to save is visible on screen rather than silently lost.

**Expected behaviour**
1. You can tick five players as fast as you can tap, and all five stay ticked.
2. After ticking five players quickly, the header reads `5 / 8` — not `3 / 8` or `4 / 8`.
3. If a tick fails to save, the number of ticked names and the count in the header disagree, so you can see something went wrong.
4. Reloading the page shows exactly the players you ticked.

**Open decision — what happens when a save fails**
A ticked box that failed to save currently stays ticked, and the mismatch with the header count is the only signal.
- Leave it, and rely on the mismatch. The coach can reload to resync. *(Recommended — a wrong count for a few seconds is a small problem, and the mismatch surfaces it rather than hiding it.)*
- Untick the box automatically so the screen is always internally consistent, at the cost of the coach not knowing why their tap "didn't take".

Decision: *(not yet made)*

### One team per tab (Slice 3)

The coach sees one team at a time and switches between them with a tab.

- [x] Show a tab per team, with the number of players in each.
- [x] Show only the selected team's players and details.
- [x] The selected tab stays selected after ticking a vest.
- [x] Active tab is highlighed

**Expected behaviour**
1. You can see which team you are looking at without reading the player names.
2. You can switch between the two teams with one tap.
4. Ticking a vest does not throw you back to the other team's tab.

**Open decision — how the tabs are labelled**
The coach recognises teams by the colour of the vests they are handing out, not by a number.
- Label the tabs by vest colour ("Orange" / "Dark") so they match what is physically in the coach's hand.
- Label them "Team 1" / "Team 2" to match the rest of the app and the database.
- Show both — a colour dot plus the team number, which is what the tab bar does today.

Decision: *(not yet made — worth checking against how the coach actually refers to the teams during a session.)*

### Moving a player keeps your place (Slice 4)

Moving a player between teams currently reloads the whole page, which would throw the coach back to the first tab every time.

- [ ] Moving a player updates both teams without a full page reload.
- [ ] The selected tab stays selected after moving a player.
- [ ] Both teams' scores, player counts and vest counts are correct after a move.

**Expected behaviour**
1. You can move a player while looking at either team's tab and stay on that tab.
2. The moved player is gone from one team and present in the other, in the right position by score.
3. Both tabs show the correct number of players immediately after the move.
4. If the moved player had been given a vest, both teams' vest counts update to match.

**Open decision — feedback when a player moves to the hidden tab**
With both teams visible, moving a player is self-explanatory — you watch them jump to the other list. With tabs, they simply vanish, and the only signal is the other tab's player count changing.
- Ship it as-is and see whether it actually feels wrong in use. *(Recommended — it may be a non-issue, and it is cheap to add feedback later.)*
- Briefly highlight the other tab when it receives a player.
- Show a short confirmation of who moved where.

Decision: *(not yet made)*

### Out of scope for now (later slices)

- **Score difference between teams.** Showing how balanced the two teams are (e.g. `Δ 0.3`) is a separate idea — useful, but it answers "are these teams fair?", not "how many vests are left?".
- **Re-rolling the distribution** from the practice page.
- **More than two teams.** Two teams is assumed throughout the app, including the database queries.
- **Keeping two devices in sync.** If an assistant coach opens the same practice on their own phone, they will not see the head coach's ticks. This is already true today and is not made worse by tabs. It only becomes worth solving if two people start handing out vests at the same time.
