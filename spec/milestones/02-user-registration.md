# Milestone 02: User registration

**Status:** done · **Branch:** feature/user-registration
<!-- Status: draft → red-teamed → in-progress → done -->

## Why
Serves **User management, steps 2 and 3**. Takes on **I16** (sign-up that an organizer initiates). Removing the CLI (I17) follows in milestone 03 (`03-remove-cli.md`), once this one is merged (Q5).

Today a new organizer gets an account only when the primary organizer opens a terminal on the Raspberry Pi over SSH and runs the CLI, then passes the password on to them. That fails the flow's success measure ("set up their account and log in securely and in less than 5 mins", rated 1). The CLI also stores passwords wrongly: it lowercases them and keeps only the first word (tech review #4, `docs/tech-review-2026-11.md`), so a chosen password like `Kesä Treeni 2026` becomes `kesä`.

Assumed impact: adding an assistant organizer takes a couple of minutes from a phone, and nobody but the new organizer ever knows their password.

## Decided
<!-- Choices Human-PM made at kickoff, step 3 and refinement, written down so
     they live in the spec rather than in one session's memory. -->
- **Invite links.** A logged-in user creates a one-time invite link in the web app and sends it themselves (for example on WhatsApp). The person who opens it chooses their own username and password. No email is sent by Tiimit.
- **Any logged-in user can invite.** Roles (I21) are not built; with three trusted users that is acceptable. The invite records who created it and which account it created; neither is shown (R7).
- **Invite links last 2 days** and work for one account (Q1). At exactly 2 days a link has expired (R9). Creating a new link does not end older ones, and a link cannot be seen again after leaving the page (R14).
- **The full link is built by the browser** from the page's own address; the server only knows the path (R1, option C). The page has a Share button that opens the phone's share sheet, and a read-only text field with the full link as the fallback.
- **The invite pages do not load PostHog**, so the token in their address never reaches a third party (R3).
- **Normalising a username**, for creating an account and for logging in: surrounding spaces are removed and letters are lowercased. All existing usernames are already lowercase, so nobody's login changes.
- **Username rule** for a new account, checked after normalising: letters only (a to z, plus å, ä, ö), 5 to 20 letters. Login does not check it, so existing usernames keep working (R5). The form shows the hint "5 to 20 letters. No numbers or spaces." under the username, as it does for the password.
- **Password rule** for a new account: 10 to 64 characters (counted as letters, not bytes), at most 72 bytes, and not the same as the normalised username. No composition rules (digits, capitals, symbols). Spaces are allowed and the password is kept exactly as typed. The form always shows the hint "At least 10 characters. A few words separated with dashes (-) works well." Existing passwords are not checked against it (Q4, Q6, R10).
- **Order of checks** when the form is sent: the link first, then the username and password rules, then whether the username is taken, and the password is hashed last (R6). Someone without a working link learns nothing about usernames and cannot make the Pi hash passwords.
- **Opening a link never uses it up**; only creating the account does (R6). WhatsApp opens links to build a preview.
- **Using up the link and creating the account happen together**: either both or neither (R2).
- **Opening an invite link while logged in** shows that you are already logged in, with a log-out button, instead of the form (Q2). Logging out goes to the front page, as it does everywhere (R8).
- **The Users page** is reached from a front-page card like Players: "Users", with the number of users as its subtext. It lists usernames only, in alphabetical order (R7).
- **The expiry is shown as a Finnish time**, for example "Works once, until Sun 23 Nov 18:40" (R14).
- **Fix #4 of the tech review is skipped**: the CLI is deleted in milestone 03 instead of fixed. Existing users keep their current passwords; a page to change one is a later idea (I26).
- Kickoff context: roles for a future multi-club Tiimit (Tiimit admin, primary organizer of a team, assistant organizers) are recorded as I27 and stay out.

## Scope
- **In:**
  - A Users card on the front page and a Users page listing who has an account, with a button that creates an invite link to share.
  - An invite link that works for one account, for 2 days.
  - Opening the link shows a form to choose a username and password; sending it creates the account and logs the new user in.
  - A message for every way the link or the form can fail (S4 to S7, S11 to S13).
  - Normalising the username in one place, used by both sign-up and login.
  - Form settings that let phone password managers suggest and save the new password (R4).
  - Updating `spec/domain-model.md` and `ARCHITECTURE.md` where they say users come only from the CLI (R13).
- **Out:**
  - Removing the CLI and creating the first user of an empty database (milestone 03, Q5). Until then the CLI keeps working as it does.
  - Roles and permissions (I21, I27). Every user keeps the same rights.
  - Changing or resetting a password (I15, I26), and a show-password toggle (I29).
  - Listing, revoking or re-sending open invites.
  - Deleting or disabling users.
  - Tiimit sending email or any other message.
  - A navigation bar (R7).
  - CSRF protection and cookie flags (I22, I23). Creating an invite is a new form with the same exposure as the existing ones.
  - Rewording `user-flows.md` step 2 ("PO creates new credentials"); that file is Human-PM's.

## Acceptance conditions (BDD)

### S1: Organizer views existing users
- Given users `testcoach` and `ainvir` exist
- And `testcoach` is logged in
- When they open the front page
- Then they see a Users card that says there are 2 users
- And the card leads to the Users page, which lists `ainvir` and `testcoach` in alphabetical order
- And the Users page has a button to create an invite
- And the Users page has a "Home" link back to the front page

### S2: Organizer creates an invite link
- Given an organizer is on the Users page
- When they create an invite
- Then they see the full invite link in a field they can copy from
- And they see a Share button
- And the page says the link works once, until a time 2 days from now
- Verification: manual for the Share button opening the phone's share sheet and for copying the link on a phone. The server's reply is tested automatically.

### S3: New organizer joins with an invite link
- Given an organizer has created an invite link
- And someone opened the link once before without sending the form
- And nobody is logged in
- When someone opens the link
- And chooses the username ` AinVir ` and the password `Kesä Treeni 2026`
- Then they are logged in as `ainvir`
- And they land on the front page
- And later they can log in as `ainvir` with the password `Kesä Treeni 2026`
- Verification: manual for the phone offering to suggest and save the password.

### S4: Someone opens an invite link that was already used
- Given someone has already joined with an invite link
- When anyone opens the same link again
- Then they see a message that the link was already used and to ask the organizer for a new one
- And no form to create an account is shown

### S5: Someone opens an expired invite link
- Given an invite link was created 2 days ago or earlier
- And nobody has joined with it
- When someone opens it
- Then they see a message that the link has expired and to ask the organizer for a new one
- And no form to create an account is shown

### S6: New organizer chooses a username that is taken
- Given a user `ainvir` exists
- And someone has opened a working invite link
- When they choose the username ` AINVIR ` and a password that follows the password rule
- Then they see a message that the username is taken
- And the username they typed is still in the form
- And no account is created
- And the same link still works for another try

### S7: New organizer sends the form with a field that breaks the rules
- Given someone has opened a working invite link
- When they send the form with any one of these:
  - an empty username
  - a username shorter than 5 letters or longer than 20
  - a username with something other than letters, such as `ain-vir` or `ainvir2`
  - a password shorter than 10 characters
  - a password longer than 64 characters
  - a password of 64 characters or fewer that is too long to store because of letters like `ä`
  - a password the same as the username, such as `AinoVirtanen` for the username `ainovirtanen`
- Then they see a message saying which field to fix and why
- And the username they typed is still in the form
- And the form still shows the password hint
- And no account is created
- And the same link still works for another try

### S8: Organizer logs in with different capitals or spaces in the username
- Given a user `ainvir` exists
- When they log in as ` AinVir ` with the right password
- Then they are logged in

### S9: Someone who is not logged in tries to see users or create an invite
- Given nobody is logged in
- When they open the Users page, or try to create an invite from a page left open after logging out
- Then they are sent to the login page
- And no invite is created

### S11: Someone who is logged in opens an invite link
- Given an organizer is logged in
- And a working invite link exists
- When they open the link
- Then they see that they are already logged in
- And they see their own username
- And they see a button to log out
- And no form to create an account is shown
- And the link still works for someone who is not logged in

### S12: Two people send the form from the same link
- Given two people have opened the same working invite link
- When the first sends the form with a username and password that follow the rules
- And then the second sends the form with a different username and password that follow the rules
- Then the first is logged in
- And the second sees the message that the link was already used
- And only the first person's account exists

### S13: Someone opens a mistyped or cut-off invite link
- Given an invite link exists
- When someone opens a link with part of its end missing, or with a made-up code
- Then they see a message that the link does not work and to ask the organizer for a new one
- And no form to create an account is shown
- And sending the form to that link creates no account

## Test notes
- All scenarios are HTTP acceptance tests with the milestone 00 harness (`cmd/web/helpers_test.go`). Made-up names only. `ainvir`, `ainovirtanen` and `testcoach` are made up.
- S1: `ainvir` is created directly through the queries, as `createCoach` does.
- S2: the test reads the invite path from the reply and checks the expiry text against a fixed "now" passed to the handler. The full link is built in the browser, so the test checks only that the page carries the path. The text field and Share button are checked by their presence.
- S3: the test opens the link once with a plain GET first (the WhatsApp preview), then signs up with a second client.
- S5: the invite is created through the same query with an expiry time in the past; no test clock (R9). One row at exactly 2 days, one at more.
- S7: one table-driven test, one row per case. The byte case uses a 64-character password made mostly of `ä`. The username rule and the password rule are core rules, so each also gets a table-driven unit test of its own (constitution, Who implements).
- S9: one row for opening the Users page, one for a create-invite request with no cookie, one with an htmx header (expects `HX-Redirect` to the login page, as `requireAuth` does today).
- S12: two clients, one request after the other; no goroutines. Truly simultaneous requests are covered by doing both steps in one transaction, not by a test (R2).
- S13: one row with a shortened real token, one with a made-up one, each for opening and for sending.
- PostHog: one assertion in S3 and S13 that the invite page's reply does not include the analytics script (R3).

## Open questions
<!-- For the red-team pass (workflow step 2) and Human-PM to resolve (step 3).
     A resolved question keeps its number and moves to the decisions log. -->

## Decisions log
<!-- Red-team findings (R1, R2, …) with Human-PM's decision on each, resolved
     questions (Q…), and decisions made while building. One line of why each. -->
- Scenarios renumbered S1 to S10 in reading order after Human-PM added the Users page scenario first (2026-11-21). Safe because nothing referred to them yet; from here on the numbers are stable.
- **Q1** Invite links last 2 days, not 7. Human-PM can always send a new link, and a shorter life limits the harm if one leaks from a chat.
- **Q2** A logged-in user who opens an invite link sees that they are already logged in, with a log-out button, and no form (S11). Switching accounts silently would be surprising, and the link stays unused for the person it was meant for.
- **Q3** Moved with the first-user settings to milestone 03, where Human-PM accepted the plain-text password as a risk (03/Q3 there).
- **Q4** Password rule: 10 to 64 characters, at most 72 bytes, not the same as the username, no composition rules, with a hint on the form (Decided). 10 is a deliberate step below NIST SP 800-63B rev 4, which asks for 15 when the password is the only login factor (https://pages.nist.gov/800-63-4/sp800-63b.html#password): Human-PM weighed typing on a phone against the risk for an app with a handful of users. 64 is NIST's lower bound for the maximum; 72 bytes is bcrypt's hard limit, which `ä` and similar letters reach sooner because each takes 2 bytes. A check against breached passwords is I28.
- **Q5** Split: this milestone keeps the invite flow; the first-user settings and removing the CLI are milestone 03 (`chore/remove-cli`). S10 moved there as 03/S1 and is a gap here.
- **Q6** The hint suggests dashes between words, but spaces stay allowed in passwords. Rejecting them would be a composition rule, which the password rule leaves out, and NIST asks verifiers to accept spaces; the hint alone steers people away from phone keyboards turning a double space into ". ".
- Red-team pass 2026-11-21 (R1 to R14), all decided.
- **R1** The browser builds the full link from the page's own address (option C), for both the Share button and the text field. No new setting and no nginx change. A configured base address (option A) becomes the right choice once Tiimit sends links outside a browser, such as the email notifications in I6. Reading the address from the request (option B) was rejected: nginx sends its own upstream address by default, and trusting that header unchecked is a known weakness.
- **R2** Accepted. New scenario S12 for two submits from the same link: only one account is created. The link is used up and the account created in one database transaction, and the submit button is disabled while sending.
- **R3** The invite pages do not load PostHog, so the token in their address is never sent to a third party. PostHog's own masking was weighed: session replay masking (`maskCapturedNetworkRequestFn`, https://posthog.com/docs/session-replay/privacy) covers recordings only, and masking ordinary events needs a `before_send` hook that rewrites every address property (https://posthog.com/docs/libraries/js/usage#amending-or-sampling-events). Skipping analytics on two pages is simpler and certain; following the sign-up in PostHog is worth nothing with a handful of users.
- **R4** Accepted in part. The form fields get the settings that let phone password managers suggest and save the password. A show-password toggle on login and sign-up is I29, out of this milestone to keep it small.
- **R5** Username rule for new accounts: letters only (a to z plus å, ä, ö), 5 to 20 letters, after normalising (spaces trimmed, lowercased). Today's usernames follow a "3 letters of first name + 3 of last name" pattern (`miisav`); it fits the rule but is not enforced. Login does not check the rule, so no existing username is affected.
- **R6** Accepted. A mistyped or cut-off link gets the message "this link does not work, ask for a new one" (S13). The link is checked first, on opening and on sending, before the username and before the password is hashed. Opening a link never uses it up (WhatsApp opens links to build a preview).
- **R7** The Users page is reached from a front-page card like Players: "Users" with a subtext showing how many users there are. Each row shows the username only. The invite records who created it and which account it created; neither is shown. A navigation bar is a separate milestone.
- **R8** Logging out on an invite page goes to the front page, as logout does everywhere. Clear and consistent; no return address to validate.
- **R9** Accepted. The invite stores its expiry time, computed in Go in UTC, and the query is given the current time, so S5's test creates an already-expired invite without a test clock. At exactly 48 hours a link has expired.
- **R10** Accepted. Password length counts letters, not bytes; the 72-byte check runs before hashing. S3's password has a capital and a space. The typed username stays in the form after an error.
- **R11** Accepted. The token is stored as is, like session ids, which live longer (14 days against 2); hashing it was considered and dropped as the less simple option. The token is never logged (log the invite id). Creating an invite is a POST. Sign-up starts the session with the same code as login.
- **R12** Accepted. S1, S3, S9 and S11 reworded at refinement, and sharing the link and the phone keyboard behaviour marked `Verification: manual`.
- **R13** Accepted. Before merge: `spec/domain-model.md` gets an Invite entry and loses "created only with the CLI; there is no sign-up"; the "Users only from the CLI" decision in `ARCHITECTURE.md` is updated. Code review runs at `high` (new public page).
- **R14** Accepted. Older links keep working until they expire and cannot be seen again; the spec says so. The page shows when the link stops working as a time, not as "2 days".
- Refinement 2026-11-21: example usernames changed from `aino` to `ainvir` and `ainovirtanen`, because `aino` is shorter than the 5-letter minimum (R5). The users list is alphabetical. The expiry time is shown in Finnish time (`Europe/Helsinki`); the binary embeds the time zone data (`time/tzdata`) so it does not depend on what is installed on the Pi.
- Refinement 2026-11-21: 12 scenarios is over the rough ceiling of 8. Kept as one milestone: they are one user goal (join with an invite), and most of the extra scenarios are the failure messages the red-team asked for. Split into three slices instead.
- Slice 2: invite times are stored and compared in whole UTC seconds. SQLite compares the stored times as text, and a time with no fraction sorts after one with a fraction in the same second; whole seconds keep the order right. Session expiry keeps the real clock, because sessions compare against SQLite's own time.
- Slice 2: the sign-up fields have no `required` attribute, so the server's messages always show instead of the browser's own pop-up. A used, expired or unknown link already shows its message instead of the form; slice 3 adds the tests and the logged-in case.
- Slice 3 checklist: logging out did nothing after signing up with a link. The session cookie had no path, so the one set by `POST /invites/<token>` was scoped to `/invites` and never reached `/logout`. The cookie now always has path `/`, and logging out goes to the front page even without a session. Found by Human-PM; the bug predates the milestone but only showed once a cookie was set from a deeper address.
- 2026-11-21: slices 1 to 3 built and accepted by Human-PM on their checklists. The existing "Logout" button now says "Log out", matching the new one. Next: code review (`/code-review high`, constitution step 6), then the R13 doc updates, then acceptance and merge.
- Wording: the app says "invitation" for the thing the organizer creates and the new user opens ("Create invitation", "Invitation expired", "This invitation was already used"). "Link" is used only for the address itself ("Invitation link" above the field). Human-PM's call after review; the specs and code keep "invite" as the short name.
- Code review round 1 (`/code-review high`, 7 findings):
  - Fixed: a built `web` binary had been committed (the `.gitignore` line `./web` matched nothing; now `/web`).
  - Fixed: the sign-up and login forms had no `method`, so before htmx loaded a tap sent the password in the address. Both now post natively too.
  - Fixed: sending the sign-up form while logged in created a second account and switched the session (S11 covered only opening the page). It now shows "You are already logged in" and leaves the link working.
  - Fixed: when the account was created but starting its session failed, the message said "try again", which would only meet a used link. It now says to log in.
  - Fixed: invite pages set `<meta name="referrer" content="same-origin">`, so the token does not reach the htmx CDN in browsers with a looser default.
  - Declined: a `busy_timeout` for truly simultaneous submits. It belongs to tech review #7 with the other connection settings; only one account is created either way, and the code comment now says so.
  - Declined: counting users for logged-out visitors on the front page. It follows the existing player-count pattern, and a failure there fails the same way as today.

## Slices
1. **Users page and invite link.** Front-page Users card, Users page, creating an invite (new table, so a migration), the page with the link, text field and Share button. The link leads nowhere yet. Covers S1, S2, S9.
   - Checklist: (1) the front page shows a Users card with the right count; (2) the Users page lists every username alphabetically; (3) creating an invite shows a full `https://` link and the time it stops working; (4) on your phone, Share opens the share sheet and the link can be sent on WhatsApp.
2. **Join with an invite.** The sign-up page and form, the username and password rules, creating the account and logging in, normalising the username at login, no PostHog on the invite page. Covers S3, S6, S7, S8, S12.
   - Checklist: (1) open a new link in a private window and sign up with a mixed-case username and a password with a space; (2) log out and log in again with the username in different capitals; (3) try a short password and a taken username, and check that the link still works afterwards; (4) on your phone, check that it offers to save the password.
3. **Link states.** Used, expired, mistyped and logged-in cases. Covers S4, S5, S11, S13.
   - Checklist: (1) open a used link; (2) open a link with its end cut off; (3) open a working link while logged in, then log out.
