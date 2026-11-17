# Milestone NN: <name>

**Status:** draft · **Branch:** <type>/<description>
<!-- Status: draft → red-teamed → in-progress → done -->

## Why
The problem and the user value, in 1–3 sentences.
Which `../user-flows.md` flow or step this serves, by flow name and step number (`Running the practice, step 2`), and which ideas (`I…`) it takes on.

## Decided
<!-- Optional. Choices Human-PM made at kickoff, written down so they live
     in the spec rather than in one session's memory. Delete if empty. -->
- ...

## Scope
- **In:** what this milestone delivers end-to-end
- **Out:** explicit non-goals, what we are deliberately NOT doing yet

## Acceptance conditions (BDD)
<!-- Group conditions into named scenarios. Name each from the user's
     perspective: a capability they gain or something they see.
     Exception: a scenario with no real user (e.g. an infrastructure or
     build check) is named by the behaviour being proven, not forced into
     a "coach" frame. Each scenario becomes a `go test` case, written
     before implementation. Plain language only, readable without code
     knowledge; every scenario has Given, When and Then, one condition per
     line. Mark one a Go test can't observe with `Verification: manual`. -->

### S1: <coach sees / can …>
- Given <state>
- When <action>
- Then <observable result>

### S2: <next scenario>
- ...

## Test notes
<!-- How each scenario is tested, by number: fixtures, what the test drives,
     how a manual one is checked. Mechanics live here, not in the scenarios. -->
- S1: ...

## Open questions
<!-- For the red-team pass (workflow step 2) and Human-PM to resolve (step 3).
     A resolved question keeps its number and moves to the decisions log. -->
- **Q1** ...

## Decisions log
<!-- Red-team findings (R1, R2, …) with Human-PM's decision on each, resolved
     questions (Q…), and decisions made while building. One line of why each. -->
- **R1** ...

## Slices
<!-- Written at refinement (workflow step 4). Each slice is one commit and
     ends with Human-PM's test checklist before the next starts. -->
1. **<slice name>.** What it builds. Covers S1, S2.
