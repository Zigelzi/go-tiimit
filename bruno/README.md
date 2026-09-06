# Bruno collection

Manual testing for the non-page endpoints — the ones HTMX calls for
fragments and actions. The full-page routes (`GET /`, `GET /players`,
`GET /practices/new`, `GET /practices/{id}`, `GET /login`) are left out:
you can already exercise those in a browser.

**One request per endpoint.** The edge cases are written up in each
request's **Docs** tab, as a table of "change this field, expect that" —
so you tweak the request in place rather than hunting through a tree of
near-identical saved requests.

## Getting started

1. Install Bruno (https://www.usebruno.com/) and open this `bruno/`
   folder as a collection.
2. Start the app: `make dev`.
3. Select the **local** environment and fill in `username` / `password`
   for a user you created with `go run ./cmd/cli`.
4. Run **auth > Login**. Bruno keeps the `user_session_id` cookie and
   sends it with everything else.
5. Set `playerId`, `practiceId` and `myClubId` to rows that exist in your
   local DB.

Handy for step 5:

```bash
sqlite3 tiimit.db "select id, name, myclub_id from players limit 5;"
sqlite3 tiimit.db "select id, date from practices order by id desc limit 5;"
```

## Fixtures

**Create practice** needs a real Excel file at
`bruno/fixtures/attendance-2026-01-15.xlsx`. That folder is gitignored,
since MyClub exports contain real names. Its edge cases need two more
files, listed in that request's docs.

## Reading the responses

Almost nothing here returns JSON — the responses are HTML fragments.
Use Bruno's **Raw** response tab rather than Preview, and watch the
**Headers** tab, because the interesting signals are headers:

| Header | Meaning |
|---|---|
| `HX-Redirect` | the action succeeded; HTMX navigates to that URL |
| `HX-Reswap: none` | `renderError` — don't swap the target, just show the banner |
| neither | a fragment meant to replace the HTMX target |

### renderError vs renderOK

Both return `200 OK`. The status code is not the signal.

| | `renderError` | `renderOK` |
|---|---|---|
| `HX-Reswap: none` | yes | no |
| Body | only `#app-error`, with a message | the fragment, then an empty `#app-error` |
| Effect in the browser | banner appears, page unchanged | page updates, any old banner cleared |

The empty banner in `renderOK` is what clears a stale error — that's why
success has to render it too, rather than just rendering the fragment.

Only `handleSetPlayerInactive` uses these helpers so far.
**players > Set player inactive** covers both paths and how to assert on
each.

## How the handlers currently signal failure

Aside from `handleSetPlayerInactive`, the handlers use three patterns:

| Pattern | Status | What the browser does | Where |
|---|---|---|---|
| Silent return | `200`, empty body | swaps the target with nothing | Create practice with no file; Get player row with a bad id |
| Bare status code | `400` / `500`, empty body | nothing — HTMX ignores 4xx/5xx bodies by default | Get team panel with team 3; Toggle vest / Move player for someone not in the practice |
| Form re-render | `200`, form fragment | inline field errors | Add player; Save player row |

## Behaviour the requests document

Each of these is written up on the request that produces it:

- **Add player** shows "must be a number" rather than "can't be empty" for
  an empty MyClub ID — both checks write the same `FieldErrors` key and
  the second overwrites the first. Its non-numeric branch also has a
  broken log format verb (`%\n`) that prints `%!(NOVERB)`.
- **Add player** stores negative MyClub IDs; the only range limit is the
  browser's `min="0"`.
- **Set player inactive** for a nonexistent player takes the success path —
  an UPDATE matching zero rows isn't an SQL error, and the query doesn't
  return rows affected.
- **Get team panel** returns an empty team with `200` for a practice that
  doesn't exist; the practice *page* 404s for the same id.
- **Toggle vest** / **Move player** return `500` when the practice/player
  pair doesn't exist, because `sql.ErrNoRows` from a `RETURNING` clause
  reaches the same branch as a database failure.
- **Create practice** distributes teams before parsing the date from the
  filename, so the bad-filename case does that work and discards it.
- **Save player row** doesn't return after setting `GeneralError` on a
  failed update; it falls through and renders a row built from a
  zero-valued player.
