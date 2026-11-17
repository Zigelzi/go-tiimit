# User flows

Human-PM's file. Agents read it and edit it only to keep references to flows and steps in sync (see [`README.md`](README.md#flows-and-steps)); the ratings are Human-PM's alone.

## Who uses it

**Organizers**
The football team members who run the practices and the team. Currently there is single primary organizer (PO) and then 0-n assisting organizers (AO), who stand in for the primary organizer when they aren't present.

Organizer or assistant organizer uses Tiimit on a phone on the practice day. 30-60 minutes before the practices PO exports the attendees and distributes the teams, also when PO isn't at the practice. Then in the beginning of the practice the organizer hands out the vests to one team by using Tiimit.

Team has at least one organizer, who is also the admin user for Tiimit. Primary organizer is the only one who has access to MyClub admin UI to export the players and attendee lists to set up the practices.

**Players**

Players don't use Tiimit directly, but by signing up to practices in separate app called MyClub. The organizer exports the attendees of a practice and uses that to distribute the teams.

## The flow

### Team player management

**Season changes**

1. Season starts
2. Players sign up to the season in MyClub
3. PO reviews which players attend to the season and if there's new players
4. PO adds new players
5. PO marks the dropped players as inactive

**Player details**

1. PO notices that some players preferences or score doesn't match
2. PO update the goalie status or score
3. PO reviews that the players score makes sense as whole

**User management**

1. PO asks a player to help with the practices
2. PO creates new credentials to the assisting organizer.
3. AO logs in

### Setting up and running practices

**Setting up the practice**

1. Players sign up to the practice in MyClub: each starts at "no answer" and moves to "attends" or "doesn't attend"
2. 30-60 minutes before the practice, PO exports the attendance report from MyClub, from a phone if needed. PO does this even when they can't come to the practice
3. PO imports the attendance report to Tiimit, which splits the players who attend into two balanced teams with a goalie on each side
4. PO reviews the teams and moves players to rebalance them
5. If PO can't come to the practice, PO sends the link to the practice to the AO who runs it

**Running the practice**

1. At the start of the practice, the organizer (PO, or AO from the link PO sent) opens today's practice on a phone
2. Organizer goes through one team and hands out vests, ticking each player who got one and seeing how many are left
3. A player who didn't answer in MyClub turns up anyway, and the organizer puts them on a team by hand, since Tiimit can't add a player to a practice that is already set up (I4)

## Success measures and ratings

<!-- Draft written by an agent from the legacy specs on 2026-11-17. Human-PM: correct the measures and fill in the ratings (1 to 5). -->

| Flow                          | Success looks like                                                      | Rating (1–5) |
| ----------------------------- | ----------------------------------------------------------------------- | ------------ |
| Team player management        | A new player is playable at the next practice without SSH or the CLI.   | ?            |
| Setting up the practice       | Teams are ready in under a minute and feel even on the pitch.           | ?            |
| Running the practice          | Vests handed out one-handed, without scrolling back to check the count. | ?            |
