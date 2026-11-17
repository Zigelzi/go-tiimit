# Spec structure & how to write specs

Read this first before adding or editing anything under `spec/`. It defines where specs live, what each file is responsible for, and how a spec gets written. The workflow that consumes these specs (red-team → refine → implement → review) lives in [`constitution.md`](constitution.md).

## Guiding rule

One fact lives in exactly one file. When two files would state the same thing, one links to the other instead of repeating it. Duplicated docs drift apart, and that is the failure mode this structure exists to prevent.

## File responsibilities

| File | Responsibility |
|---|---|
| `constitution.md` | How we work: agentic workflow, tech stack, git conventions, UI principles. Governance. |
| `user-flows.md` | Human-PM's file: who uses Tiimit, the flows they go through, and each flow's success measure and current rating (1 to 5). Agents read it, and edit it only to keep references to flows and steps in sync; the ratings are Human-PM's alone. |
| `product-notes.md` | Reference for agents: dated observations from use, and the ideas list (`I1`, `I2`, ...). Agents keep it up to date. |
| `domain-model.md` | Shared entities (Player, Practice, …). Each entity and how it is stored. |
| [`../ARCHITECTURE.md`](../ARCHITECTURE.md) | How the code is organised (at the repository root): the layers, the rules between them, and where new code goes. Updated before a merge when a layer, rule, package job or decision changed; see step 7 of `constitution.md`. |
| `milestones/NN-name.md` | One unit of work. See below. |
| `milestones/_template.md` | The template new milestone specs are copied from. |
| `legacy/` | The checklist specs from before this workflow (auth, players, practices, the improvements backlog). Read-only history: what was built and why. Their open items now live as ideas in `product-notes.md`; don't add to these files. |

## Milestones

```
spec/milestones/
  _template.md
  00-test-harness.md
  01-<name>.md
  ...
```

- One milestone = one file = one branch = one squashed commit. The filename mirrors the branch description (`01-practice-tabs.md` ↔ `feature/practice-tabs`), and reviewing a milestone means reading one spec against one diff.
- Number by build order (`00`, `01`, …). Copy `_template.md` to start a new one.
- Each spec has these sections, in the order of `_template.md` (kept deliberately minimal):
  - Status + Branch: the spec is the status tracker; there is no separate tool. `draft → red-teamed → in-progress → done`.
  - Why: the problem and user value, and which `user-flows.md` flow or step it serves.
  - Decided (optional): choices Human-PM made at kickoff, so they live in the spec and not in one session's memory.
  - Scope (In / Out): the Out list is the main lever against scope creep: state what this milestone deliberately does not do.
  - Acceptance conditions (BDD): Given/When/Then, grouped into named scenarios. Write them so someone with no code knowledge (another coach) can read and challenge them: no table, query or attribute names. Every scenario has all three of Given, When and Then, and each condition is its own line: a second outcome is a new `And` line, not joined to the first with a comma or "and". How a scenario is tested goes in a separate Test notes section below the scenarios. Name each from the user's perspective (a capability or something they see); name a scenario with no real user (an infrastructure/build check) by the behaviour proven instead. Each scenario becomes a `go test` case, written before implementation. This is the section that matters most.
  - Number scenarios `S1`, `S2`, … in the heading (`### S3: Coach moves a player to the other team`) so they can be referenced in conversation, open items and code review. Numbers are stable: a new scenario takes the next free number, a removed one leaves a gap; never renumber. The matching test carries the number in its comment (`// S3: Coach moves a player to the other team`). Outside the spec, prefix the milestone: `01/S3`.
  - Verify each scenario the most feasible way. A `go test` is the default; behaviour a Go test can't observe (e.g. how a tab feels to tap on a phone) is marked `Verification: manual` under the scenario and checked by Human-PM at step 7. Still test the part that is observable (e.g. the server's reply) automatically.
  - Open questions: where a draft hands off to the red-team pass and back to Human-PM. Number them `Q1`, `Q2`, …. Same stability rule as scenarios: never renumber, a resolved one keeps its number and moves to the decisions log.
  - Decisions log: each red-team finding (`R1`, `R2`, …) with Human-PM's decision, each resolved question, and decisions made while building, with one line of why. A later session reads this instead of the conversation that made the decision.
  - Slices: written at refinement. Each slice names what it builds and which scenarios it covers, and ends with Human-PM's test checklist.

## Referencing things by number

One letter per kind of thing, so a number is never ambiguous across documents:

| Prefix | What | Lives in |
|---|---|---|
| `S` | A BDD scenario | A milestone spec. Prefix the milestone outside it: `03/S8`. |
| `Q` | An open question about that milestone | The same milestone spec. |
| `R` | A red-team finding and the decision on it | The milestone spec's decisions log. Prefix the milestone outside it: `04/R2`. |
| `I` | An idea for improving the product | [`product-notes.md`](product-notes.md), because ideas outlive any one milestone. |

Code comments are always outside the spec, so they always carry the milestone: `04/R2`, never a bare `R2`, since every milestone numbers from 1.

All four follow the scenario rule: numbers are stable, a removed one leaves a gap, never renumber.

### Flows and steps

[`user-flows.md`](user-flows.md) owns the flows (the bold headings, such as **Running the practice**) and their numbered steps, and is the single source of truth for them. Each flow numbers its steps from 1, so other files refer to a flow by its name and to a step by flow name plus number: `Running the practice, step 2`. A bare `step 2` is ambiguous. The section headings of `product-notes.md` are flow names, and so are the rows of the success table.

Unlike S, Q and I numbers, flows may be renamed and steps renumbered. When that happens, the same commit updates every reference, including those in finished milestone specs:

1. Find them: `grep -rn "user-flows.md\|, step [0-9]" spec CLAUDE.md README.md`, plus a grep for the old flow name. Ignore the workflow steps of `constitution.md`, which are a different numbering.
2. Update each reference, including inside `user-flows.md` itself and in `product-notes.md`.

## Starting a milestone

Open a fresh session and describe what you want to build; the session creates the `feature/<name>` branch and scaffolds the spec file from the template, then you write the spec into it.

A good kickoff description covers:

- The user goal: what the coach is trying to do and why, from their point of view. "Let me see how many vests are left while I hand them out," not "add a sticky header."
- The thin slice: the smallest end-to-end path that delivers that goal. If it doesn't fit in a sentence, it is probably two milestones.
- One user goal per milestone. A kickoff joining two goals with "and" is two milestones. Rough ceiling: about 8 scenarios.
- In / out: what this milestone deliberately leaves for later. This is the strongest guard against scope creep.
- Which `user-flows.md` flow or step it serves, and which ideas (`I…`) it takes on.
- What you are unsure about: these seed the red-team pass and the open questions.

Leave out the how: no database schema, no page layout, no test wording. Rough is fine; red-team and refine sharpen it.

Example kickoff:

> Milestone 1: when I move a player to the other team, keep me on the tab I was looking at, so I don't lose my place while balancing teams (I7). In: moving updates both teams' counts, scores and vests without a reload. Out: highlighting the other tab, more than two teams. Serves Setting up the practice, step 4. Unsure: is it confusing when the moved player just disappears?

## How a spec gets written

Steps 1–4 of the workflow in [`constitution.md`](constitution.md#workflow): Human-PM drafts intent (not perfect BDD), a red-team pass adds to Open questions, Human-PM resolves them, and the main session refines the scenarios into the tests.
