---
author: Ahmad Ibrahim
date: August 2026
paging: Slide %d / %d
---

# Stacked PRs with gh-stack

Breaking big changes into small, reviewable, unblocked pieces.

_Ahmad Ibrahim · Engineering_

---

# Agenda

1. The problem: a day in the life
2. The mental model: what "stacking" is
3. Why it helps: review, flow, throughput, safety
4. How **gh-stack** works
5. **Live demo** 🔴
6. Gotchas & reporting
7. The landscape & how to start

---

# Meet the cast

- **Alice**, shipping a feature: *scheduled cluster backups*
  - API type · validating webhook · controller · UI
- **Bob** will review it

A normal week. Alice hits a fork: **how does she ship this?**

Two options. Both are bad.

---

# Option A: the mega-PR

Everything in one PR: **+20,000 lines**.

- The whole feature is gated behind **one enormous review**
- Nothing lands until Bob clears all 20k
- The branch **rots** against `main`; conflicts accrue

---

# Bob vs the mega-PR

Bob gets pinged. Opens +20,000 lines.

- Burn **hours** on a real review, or **LGTM** and move on
- Focus fades over a giant diff → rubber-stamp → **bugs slip through**
- Alice got **no early feedback**: the design was locked in 20k lines ago

---

# Option B: smaller PRs

The *right* instinct: split it up.

`api → webhook → controller → ui`

But each layer **depends on the one below**, and GitHub PRs target `main`.

---

# Option B's trap

To keep moving, Alice must either:

- **Wait** for PR1 to merge before starting PR2 → serialized, idle, **blocked**
- Or **branch off the unmerged PR1** and hand-rebase the whole chain every time review churns it → **rebase hell she owns**

---

# The false dilemma

> Mega-PR (unreviewable) **or** small PRs (blocked / rebase hell).

Nobody was lazy. The **tooling forced the choice.**

---

# What if you didn't have to choose?

Small PRs (Option B), **without** the blocking.

That's a **stack.**

The PR stops being the unit of work. The **layer** is.

---

# The mental model

One big branch → a chain of small branches, each based on the one below:

```
   feat/ui          ●  PR #4  → base: feat/controller
                    │
   feat/controller  ●  PR #3  → base: feat/webhook
                    │
   feat/webhook     ●  PR #2  → base: feat/api
                    │
   feat/api         ●  PR #1  → base: main
                    │
   main             ┴
```

Each PR is small. Each targets the one below. Together they're the feature.

---

# The same feature, stacked

| # | Layer | Targets | Size |
|---|-------|---------|------|
| 1 | **api**: `BackupSchedule` type | `main` | small |
| 2 | **webhook**: validate schedule | api | small |
| 3 | **controller**: reconcile CronJob | webhook | small |
| 4 | **ui**: console form | controller | small |

Four focused reviews instead of one 20k-line slog.

---

# Smaller diffs = better review

Batch size drives review **quality and speed**.

- A 200-line PR gets a real review. A 20k-line PR gets an LGTM.
- Bugs are caught where they're introduced, not archaeology later.

---

# Unblock yourself

You don't wait on review to keep building.

- Open PR #1, then **stack PR #2 on top and keep going**
- Review happens on lower layers while you work up top
- No branching off in-review branches by hand

---

# Throughput

Small batches **ship faster**.

- Shorter review latency per PR → shorter cycle time
- Layers merge as they're approved, not all-or-nothing
- The "performance" story: flow, not heroics

---

# Continuous & agile

Small changes land **frequently**, not in one big drop.

- Tighter feedback loops → closer to true continuous integration
- **Small blast radius:** a regression is trivial to isolate
  - bisect / revert **one layer** vs digging through a 20k merge

---

# Alice's week, restacked

Same feature. Same four changes. As a stack:

- ✅ Early feedback on the API **before** the UI exists
- ✅ Bob reviews four focused PRs, actually catches bugs
- ✅ Alice keeps building up the stack, **never blocked**
- ✅ No hand-rebasing: the tool restacks for her

Same people. The **tooling** changed the outcome.

---

# OK, how do I actually do this?

Meet **`gh-stack`**, GitHub's native stacked-PR extension.
(Public preview. We have early access.)

---

# Install

```bash
gh extension install github/gh-stack
gh stack alias          # optional: adds `gs`
```

Now `gh stack …` (or `gs …`) manages your stack locally, then pushes to GitHub.

---

# The five beats

```
create  →  push  →  submit  →  review  →  merge
```

Build layers locally, push them, open the PRs, get focused reviews, land them.

---

# Build the stack

```bash
gh stack init                       # start a stack on main
# ... edit the API layer ...
gh stack add -Am "api: BackupSchedule type" feature/backup-api
# ... edit the webhook layer ...
gh stack add -Am "webhook: validate schedule" feature/backup-webhook
```

`add` stages, commits, and stacks a new branch on top, all in one command.

---

# Ship the stack

```bash
gh stack submit     # push all branches + open/update every PR
gh stack view       # see the stack locally
```

`submit` sets each PR's base to the layer below, automatically.

---

# The stack map

What Bob sees on github.com. Navigate the layers:

```
  Stack: add-scheduled-backups
    ○ #4  ui: console form            (top)
    │
    ○ #3  controller: reconcile
    │
    ○ #2  webhook: validate schedule
    │
    ● #1  api: BackupSchedule type    ← reviewing
    │
    main
```

(We'll see the real one live in a minute.)

---

# Cascading merge

Merge any layer → everything below it lands, and the rest **auto-rebase**.

```bash
gh stack merge      # or click Merge on a PR in the UI
gh stack sync       # pull the cascade down locally
```

No manual rebasing of the branches above. The tool restacks them.

---

# git → gs cheat sheet

| Goal | Manual git/gh | gh-stack |
|------|---------------|----------|
| Start bottom layer | `git checkout -b feat-api main` | `gs init` |
| Add dependent layer | `git checkout -b feat-webhook feat-api` | `gs add -Am "…" feat-webhook` |
| Push branches | `git push -u origin <each>` | `gs push` |
| Open PRs w/ right bases | open each PR, hand-set base | `gs submit` |
| Restack after change/merge | `git rebase` each branch in order, fix conflicts | `gs sync` (or `gs rebase`) |

**`gh-stack` is Option B, minus the manual pain.**

---

# Let's see it live 🔴

A real 3-layer stack: **api → webhook → controller**.

---

# Review UX & CI

- CI runs on **each PR as if it targets `main`**: every layer is tested standalone
- Branch protection is enforced on the **final target branch**, not the intermediate layer branches

---

# Gotchas

- Keep the stack **in sync**: `gs sync` after merges / review churn
- **Don't stack everything**: a tiny one-shot change is just a PR
- Stacking rewards **discipline**: small, coherent layers with clean boundaries

---

# Stacks multiply CI load

Each layer runs CI → an N-layer stack ≈ **N× the runs** of one PR.
Restacks (review churn or a merge) **re-trigger** down the chain.

- Mitigate: skip redundant intermediate runs, path filters, concurrency-cancel, a merge queue, right-sized runners
- **Bigger picture:** AI-generated code is already inflating PR volume. CI must scale to changing throughput **regardless**. Stacking just surfaces the need sooner.

---

# It's public preview, tell them

We're early adopters. Our feedback shapes the tool.

- `gh stack feedback`  →  `gh.io/stacks-feedback`
- Bugs / requests: **github/gh-stack** issues
- Docs: `gh.io/stacks`
- Internal: **#eng-gh-stack** to flag problems & get help

---

# You're not the first

The stacking landscape:

- **Graphite**, **Aviator (`av`)**, **spr**, Meta's **ghstack / Sapling**, **git-town**

Stacking is a proven workflow. gh-stack makes it **native to GitHub**.

---

# Why gh-stack for us

- **Native GitHub**: no third-party app, no extra permissions
- The **stack map** is built into the PR UI
- **CLI is optional**: works alongside the normal PR flow
- Zero new SaaS in the review path

---

# Try it Monday

- A few of us are already stacking, and it works
- Pick **one** medium feature this week
- Split it into **2-3 layers**, `gs init` / `gs add` / `gs submit`
- Docs: `gh.io/stacks`

---

# Thanks

**Stacked PRs:** small, reviewable, unblocked.

- Escape the false dilemma
- Faster, safer, more continuous delivery
- `gh extension install github/gh-stack`

Questions?
