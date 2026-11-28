# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project status

Deliberately not restated here. This file is loaded every session, which makes a status summary in it go stale. `spec/README.md` makes the rule explicit ("one fact lives in exactly one file", "the spec is the status tracker; there is no separate tool"), so:

- **What is built and where each milestone stands:** `spec/milestones/NN-*.md`, each carrying a Status line (`draft -> red-teamed -> in-progress -> done`). Everything built before the switch to this workflow (2026-11) is described in `spec/legacy/`. `git log --oneline main` is the merged list, one squashed commit per milestone from then on.
- **The entities and how they are stored:** `spec/domain-model.md`.
- **What each spec file is for:** the table in `spec/README.md`.

Do not restate any of that here. Add to the file it belongs in and let this section keep pointing.

Code layout and the rules for where new code goes: `ARCHITECTURE.md` at the repository root. Read it before adding or moving code, and update it before a branch merges to `main` when the change altered a layer, a rule, a package's job or a recorded decision (step 7 of `spec/constitution.md`). Do not restate it here. `go doc ./internal/practice` (and the other packages) lists what exists.

## What Tiimit is

Tiimit is an app that football team organizer (or coach) use to distribute the players of a practice into two even teams.

Tiimit is a web app (templ + HTMX).

Currently it's used by a single team of a single football club and has fewer than 10 users, deployed to a Raspberry Pi and reachable from the public web behind login. See `README.md` and `spec/user-flows.md`.

## Commands

- `make generate`: sqlc, templ and Tailwind. `make run` and `make prod/build-arm64` run it first, so a fresh checkout works with one command. The first run downloads the pinned standalone Tailwind CLI into `./bin` (no Node/bun)
- `make test`: templ generate, then every test (it does not need sqlc or Tailwind). A single test: `go test ./internal/practice/ -run TestDistribute -v`
- `make run`: the server with `.env`'s settings, on `localhost:8080`
- `make dev`: live reload on `localhost:7331` (templ watch + proxy, Air, Tailwind watch)
- `make prod/build-arm64`: arm64 binary for the Pi into `./build`
- `make deploy`: `make test`, build, then back up the database, install and check on the Pi (needs `deploy/deploy.conf`, copied from `deploy/deploy.conf.example`). **A deploy is also a migration of the live database.** Only when Human-PM asks. Setup and recovery: `docs/deployment-guide.md`
- After editing SQL: `sqlc generate`. After editing a `.templ`: `templ generate` (or have `make dev` running) before `go build`/`go test`/`go vet`
- Migrations run automatically against `DB_PATH` on every server startup. Env vars: see `.env.example` (`POSTHOG_API_KEY` is optional)

Generated code: a file is generated if, and only if, its first line is `// Code generated ... DO NOT EDIT.` (templ's `*_templ.go`, sqlc's `internal/db/*.sql.go`, `db.go`, `models.go`). Change its source (`.templ`, `sql/`) instead. sqlc's output is committed; `*_templ.go` and the built `cmd/web/static/tailwind.css` are gitignored, so regenerate them on a fresh checkout. Everything else, including `internal/db/config.go`, is edited directly.

Stack (decided; see `spec/constitution.md`): Go + [templ](https://templ.guide/llms.md) + HTMX + TailwindCSS, SQLite (`modernc.org/sqlite`), Goose migrations, sqlc.

## How work is meant to be done here

`spec/constitution.md` is the source of truth for the workflow, and `spec/README.md` for how specs are structured and written. Read them, don't rely on this summary. The operative rules for a session working here:

- Two actors: Human-PM (the user, who makes the product decisions) and the main session (does most phases inline).
- Spawn sub-agents at only two points: red-teaming the spec, and code-reviewing the diff. Implementation runs inline in the main session; a sub-agent only for a slice of roughly 100+ hand-written lines.
- When goals conflict, the earlier goal wins (feedback, then quality, then tokens). Prefer targeted reads over broad codebase sweeps.
- BDD acceptance conditions become tests before implementation, written from the spec, not from the code.
- Vertical slices, sequentially. Stop after every slice: commit, give Human-PM a 2–4 step test checklist, wait.
- Code-review loop is bounded (max 2 rounds). Factual disputes → settle with a test; subjective → the simpler option; product/scope → escalate to Human-PM as a plain-language trade-off.
- Starting a milestone: Human-PM opens a fresh session and describes what to build; the main session creates the `feature/<name>` branch and scaffolds the spec from `spec/milestones/_template.md` on that branch. Ideas to pick from: `spec/product-notes.md`.
- Commits are the agent's job, with standing authorization to commit without asking. Conventional Commits (`type(scope): summary`), Conventional Branch names (`feature/practice-tabs`). Squash the branch into one commit, merge locally to `main` (`git merge --squash`) after Human-PM accepts, then push to `origin` (github.com/Zigelzi/go-tiimit). Deploying is a separate step Human-PM asks for.

The three goals shaping every decision, in priority order: (1) reduce time to feedback by right-sizing and scoping the work, (2) maximize code quality and maintainability, so a human or an LLM can extend the codebase easily and it is secure, (3) optimize token usage, to learn how to run effective agentic workflows.

## Scope constraints to respect

- One club, fewer than 10 users. Keep the code simple rather than scaled; don't build multi-club features until a milestone asks for them.
- Real player names are personal data: every page showing them requires login, and real MyClub exports (`*.xlsx`, `bruno/fixtures/`) are never committed. Tests use made-up names.
