# LLM Constitution

Tiimit started as a learning project, built by hand to learn programming. From 2026-11 it is built with an agentic workflow instead, because Human-PM no longer has the time to write most of the code themselves. The workflow is the one that proved itself in Roundtrip (`../roundtrip/spec/constitution.md`), adapted here. Human-PM still wants to be able to read and understand the code, so interpretability stays a first-class goal.

## Agentic workflow

Two actors: Human-PM (me, making the product decisions) and the main session (the agent driving the work and, where noted, spawning sub-agents).

Goals:

1. Reduce time to feedback by right-sizing the work and scoping it correctly.
2. Maximize code quality and maintainability, so the codebase can be easily extended by a human or an LLM, and is secure.
3. Optimize token usage, to learn how to run effective agentic workflows.

**When goals conflict, the earlier goal wins: fast feedback first, then quality, then tokens. None of them is served by elaborate multi-agent orchestration.** At this scale (one club, fewer than 10 users) the main session does most phases inline. Sub-agents are spawned only at the two points where independent context genuinely pays for itself: red-teaming the spec, and reviewing the diff. (A spawned agent starts cold and re-derives context, so each spawn has a real token cost.)

Version control is part of the workflow. Human-PM made the baseline commit that introduced this workflow; after that all commits are made by the main session, not the Human-PM; this is a standing authorization to commit without asking each time. Each milestone runs on its own branch off `main`; each vertical slice is one focused commit on that branch once its acceptance test is green. The branch diff is what the code-review loop reviews, and the branch merges to `main` only after Human-PM accepts it. The branch is created at the start of the milestone (step 1) and the spec is authored on it (branch-first), so a discarded milestone leaves nothing on `main`.

Conventions:

- Commit messages follow Conventional Commits: `type(scope): summary` (types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`). Per-slice commits on the branch are working checkpoints. Commits before the switch to this workflow do not follow it; leave them as they are.
- Branch names follow [Conventional Branch](https://conventionalbranch.org/): `<type>/<description>`, lowercase `a-z0-9` + hyphens, no leading/trailing/consecutive hyphens (e.g. `feature/practice-tabs`). Use the purpose prefix matching the milestone (`feature`/`fix`/`chore`/`hotfix`/`release`).
- Squash before merge. On acceptance, squash the branch's slice commits into a single Conventional Commit so `main` stays a clean, linear log of one commit per milestone.
- Remote: `origin` is github.com/Zigelzi/go-tiimit. Merges to `main` are local (`git merge --squash`) and `main` is pushed once Human-PM accepts the milestone. No pull requests: the code review in step 6 happens on the local branch diff.

### Workflow:

1. In a new session, Human-PM describes what to build. The main session scaffolds the milestone: it creates the `feature/<name>` branch off `main`, copies `_template.md` to `NN-<name>.md` with Status and Branch filled, and commits it. Then Human-PM writes the minimal spec (what and why) into the file and the main session commits the draft.
2. Spawn a red-team agent on the spec: weak points, clarifying questions, usability/feasibility risks. It also checks the acceptance conditions against the rules in [`README.md`](README.md#milestones): plain language readable without code knowledge, and Given, When and Then in every scenario. The red-team runs before refinement (step 4), so it cannot catch mistakes the refinement introduces; Human-PM asks for a second pass when a refined spec looks like it is drifting. Human-PM may skip the red-team for a small spec.
3. Human-PM decides which findings to accept. _(checkpoint)_
4. Main session refines this milestone's acceptance conditions into precise Behaviour-Driven Development (BDD) scenarios that deliver the feature end to end, and writes its Slices. Those conditions become acceptance tests before implementation, written independently of the code that will satisfy them.
5. Main session implements milestone by milestone as vertical slices, sequentially, each slice satisfying its BDD acceptance test. A task is right-sized if one agent can implement and self-verify it in one context without re-loading the whole app. Each slice is committed to the milestone branch once its test is green.
   - **Stop after every slice.** Commit it, hand Human-PM a 2–4 step test checklist for that slice alone, and wait. This is the default, not a question to ask. Chaining slices and delivering one long checklist at the end lets a defect from slice 1 get built on.
   - **Before the first slice, raise the model choice** (see [Model & effort](#model--effort)). The main session says in one line what kind of work the slices actually are (mechanical, or a genuine design problem), and Human-PM decides.
   - **Tests first, implementation inline.** See [Who implements](#who-implements).
   - **Visual bugs: measure before the second fix.** When a defect from the checklist is visual and the first fix does not resolve it, stop fixing by reasoning. Reproduce it in a headless browser, with screenshots at phone pixel density, before changing more code. The main session does this inline; it is not a new spawn point.
6. Code-review loop (bounded).
   a. Spawn a code-review agent (`/code-review`) on the finished feature's diff. It tags each finding blocking (correctness / bug / security / a failing BDD condition) or non-blocking (style / preference). A core domain rule added or changed without unit tests (see [Who implements](#who-implements)) is blocking. The review reads the sources, not generated code: it skips every file whose first line is `// Code generated ... DO NOT EDIT.` (sqlc's output in `internal/db`; `internal/db/config.go` is not generated and is reviewed) and reviews the `.templ` files and `sql/` they come from. Say so in the review request. b. Implementer addresses each: fix it, or decline a non-blocking one with a one-line reason. c. Disagreements resolve by type: factual (is there a bug?): settle with a test, the test decides; subjective (too complex?): default to the simpler option; product/scope in disguise: escalate to Human-PM as a plain-language trade-off, never as a technical vote. d. Max 2 rounds. Blocking findings must be fixed or proven-not-a-bug before the feature is done; anything unresolved escalates to Human-PM. e. Changes after review. Small fixes found while testing (wording, styling, a tweak to an existing scenario) go straight in with their tests. A new scenario either gets a short second review round on its own diff or moves to the next milestone; Human-PM picks.
7. Report to Human-PM: what's done, what to test/verify, and request feedback → continue / complete / discard. The report includes a test checklist: a few concrete steps to try in a browser, each with what should happen, drawn from the milestone's scenarios, always covering the ones marked `Verification: manual`. It also says in one line whether the milestone adds a migration, because **a deploy is also a migration** of the club's live data (see [`../docs/deployment-guide.md`](../docs/deployment-guide.md)). Before the merge, the main session checks [`../ARCHITECTURE.md`](../ARCHITECTURE.md) against the branch diff and updates it on the branch only if the change altered a layer, a rule between layers, a package's job or a decision the doc records. Moving code between files or adding a function does not count. The update goes into the same squashed commit, and the report says in one line whether the doc changed. A `chore/` or `refactor/` branch with no spec follows the same rule. On accept, the main session sets the spec's Status to `done`, merges the milestone branch to `main` and pushes. The spec's Status line is the only place the milestone's state is recorded, so that edit is not optional bookkeeping. Deploying to the Pi (`make deploy`) is a separate step Human-PM asks for; the main session does not deploy on its own. _(checkpoint)_

Sessions. Start each milestone in a fresh session to keep context and token cost small; the spec, `CLAUDE.md`, and memory carry what a new session needs to pick up. The spawned sub-agents (steps 2 and 6) run in isolated context, so they never require restarting the main session.

Skills vs sub-agents. A skill encodes how to do a repeatable task; a sub-agent provides isolated fresh context. Reuse the built-in skills that already fit: `/code-review` (step 6), `/run` (step 7). Do not build custom skills upfront; extract one only after running the procedure manually a couple of times. Refine / implement stay as prose here, not skills.

Spec organization and how specs are written: see [`README.md`](README.md).

## Who implements

Roundtrip ran an experiment (2026-09 to 2026-11) comparing implementation in the main session against implementation in sub-agents. It settled in practice on the main session, and Tiimit adopts that result instead of re-running it:

- **The main session writes the BDD acceptance tests first**, from the milestone's scenarios, before any implementation exists, and sees them fail. The tests are a contract written first, even though the same session then writes the code.
- **Core rules also get unit tests.** Acceptance tests prove a scenario end to end, but they do not pin down a rule's edge cases. When a slice adds or changes a core domain rule (the logic the app exists for, such as distributing players into teams in `practice.Distribute`, scoring players or picking goalies), the main session also writes table-driven unit tests for that rule. They are written before the implementation, together with the acceptance tests. They cover the cases no coach would write a scenario for: no players, an odd count, no goalie or only one, equal scores. Keep it proportional: a handler, a query or code that only passes data along is covered by the acceptance tests and gets no unit test of its own. Where these tests live: [`../ARCHITECTURE.md`](../ARCHITECTURE.md#conventions).
- **The main session implements inline.** A spawn's fixed start-up cost (about 40k tokens in Roundtrip) is several times a typical slice of a few dozen lines in two or three files.
- **A sub-agent only for a big slice**, roughly 100+ hand-written lines. Its brief gets the red tests, the milestone spec, the slice's scope and a pointer to `ARCHITECTURE.md`. Decisions must live in the spec, not the conversation, since the sub-agent has no memory of it.
- **A test that looks wrong stops the implementer.** Some tests written before the code will be wrong in ways only the implementation shows. The implementer stops and says which test and why; it does not bend the code to pass it. The test is fixed keeping what it checks, and Human-PM is told. A sub-agent's brief says this explicitly, since it cannot ask mid-run.

## Model & effort

Two dials trade cost for quality: model tier and reasoning effort. Match both to two questions: how hard is the reasoning, and how costly is a mistake here? Light default: don't over-tune it. Context size × number of turns costs more than model tier, so fresh sessions and targeted reads matter most. Switching model mid-session re-sends the whole context uncached on the new model, so switch rarely.

- The two spawned critics (red-team in step 2, code review in step 6): strong model at high effort. They read small inputs, so capability is cheap there and high-leverage.
- Mechanical stretches (scaffolding, commits, the step-7 report): a fast, cheap model at low effort is fine.
- Code-review effort is a skill argument: `/code-review medium` by default, `high` for changes touching auth, sessions, personal data (player names) or data integrity, and avoid `ultra` unless a milestone is large or risky.

## Tech stack

Chosen to maximize interpretability (Human-PM wants to read what is happening), minimize dependencies and token cost, and deploy simply to the Raspberry Pi.

- Go: backend language; compiles to a single static binary for trivial Pi deployment.
- templ: type-safe HTML templating in Go (https://templ.guide/llms.md).
- HTMX: server-driven interactions; responses are literal HTML, the most transparent thing to debug.
- TailwindCSS via the standalone CLI binary (no Node/npm/bun toolchain).
- Alpine.js: the chosen tool for client-side reactive interactions, added only when a milestone first needs one. Prefer native HTML (`<details>`, `<dialog>`) and small plain scripts (like today's `static/tabs.js`) for simple toggles.
- SQLite (`modernc.org/sqlite`, pure Go): a single file with no server process.
- Goose: versioned migrations, embedded and run on every startup.
- sqlc: generates type-safe Go from hand-written SQL.
- BDD conditions map to standard `go test` (table-driven where it fits), not a Gherkin framework like `godog`. Acceptance tests drive the real HTTP handlers through `httptest` against a fresh SQLite database (built in milestone 00).

## UI principles

Tiimit is used on both desktop and phone. About 90 % of production use is on a phone, often standing on the pitch, often one-handed (handing out vests); Human-PM tests mostly on a laptop browser. Mobile is the primary target.

- Mobile first: build and test each feature at phone width first; a desktop layout is a bonus, never the driver.
- Function before polish: ship basic styling and keep moving; visual passes come later.
- Usability over style on mobile: large tap targets and simple interactions win over visual refinement.
- The line between the two: mobile usability essentials (tap-target size, readable text, reachable controls, obvious state) count as function, not polish.

## Data and privacy

The app is reachable from the public web and stores real people's names (PII). Every page that shows player data requires login. Real MyClub exports (`*.xlsx`, `bruno/fixtures/`) never get committed; tests use made-up names. A milestone that touches auth, sessions or the public surface gets `/code-review high`.
