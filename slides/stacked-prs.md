---
author: Ahmad Ibrahim · Stacked PRs with gh-stack
paging: Slide %d / %d
---

# Stacked PRs with gh-stack

```
   ┌──────────────┐        ┌────────────┐
   │              │        └────────────┘
   │              │        ┌────────────┐
   │              │        └────────────┘
   │              │  ──▶   ┌────────────┐
   │              │        └────────────┘
   │              │        ┌────────────┐
   └──────────────┘        └────────────┘
```

Breaking big changes into small, reviewable, unblocked pieces.

_Under the Hood · Ahmad Ibrahim · August 19th, 2026_

---

# Agenda

```
  1 · problem        a typical dev's workflow
  2 · mental model   what "stacking" is
  3 · why it helps   focused reviews · unblocking development
  4 · gh-stack       command overview
  5 · live demo      stacking changes on each other
  6 · reality        gotchas & CI cost
  7 · landscape      alternatives & why gh-stack
```

---

# 1 · The problem
## Meet the cast

```
  Alice builds:  scheduled cluster backups
  ┌────────────┬────────────┬────────────┬────────────┐
  │    api     │  webhook   │ controller │     ui     │
  │  new type  │  validate  │  reconcile │  console   │
  └────────────┴────────────┴────────────┴────────────┘
  Bob reviews.
```

Alice hits a fork: **how does she ship this?**

Two options. Both are bad.

---

# 1 · The problem
## Option A: the mega-PR

Everything in one PR:

```
  PR #1  "add scheduled backups"
  ──────────────────────────────────────────
   api  +  webhook  +  controller  +  ui
   +20,000   -0     across 73 files
  ──────────────────────────────────────────
              one giant review
```

- Gated behind **one enormous review**; nothing lands until it's all cleared
- The branch **rots** against `main`; conflicts accrue

---

# 1 · The problem
## Bob vs the mega-PR

```
  reviewing +20,000 lines...
  file  1/73   ██████████████   real focus
  file 20/73   ████████         skimming
  file 45/73   ████             "looks fine"
  file 73/73   █                LGTM  ▶  bugs slip
```

- A real review costs **hours**; the giant diff pushes toward **LGTM**
- Alice got **no early feedback**: the design was locked in 20k lines ago

---

# 1 · The problem
## Option B: smaller PRs

The *right* instinct: split it up.

```
   api  ─▶  webhook  ─▶  controller  ─▶  ui
             (each layer needs the one below)
```

But plain PRs all target `main`:

```
   webhook ─▶ main ?     controller ─▶ main ?
      wrong base · wrong diff · conflicts
```

---

# 1 · The problem
## Option B's trap

```
                 split into layers
                        │
            ┌───────────┴───────────┐
            ▼                       ▼
       WAIT for PR1            STACK by hand
       to merge first         off unmerged PR1
            │                       │
       idle · blocked          rebase every layer
                               on each review churn
            ▼                       ▼
         (stuck)                (rebase hell)
```

---

# 1 · The problem
## The false dilemma

```
      mega-PR                       small PRs
   unreviewable                blocked / rebase hell
       ✗                              ✗
        └──────────────┬──────────────┘
                       ▼
            the tooling forced the choice
```

Nobody was lazy.

---

# 2 · The mental model
## What if you didn't have to choose?

```
   Option B's small PRs   +   no blocking
   ───────────────────────────────────────
                  =   a STACK
```

The PR stops being the unit of work. The **layer** is.

---

# 2 · The mental model
## What a stack is

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

# 2 · The mental model
## The same feature, stacked

| # | Layer | Targets | Size |
|---|-------|---------|------|
| 1 | **api**: `BackupSchedule` type | `main` | small |
| 2 | **webhook**: validate schedule | api | small |
| 3 | **controller**: reconcile CronJob | webhook | small |
| 4 | **ui**: console form | controller | small |

```
   api         ■■
   webhook     ■■■
   controller  ■■■
   ui          ■■
   mega-PR     ■■■■■■■■■■   same code, one giant diff
```

Four focused reviews instead of one 20k-line slog.

---

# 3 · Why it helps
## Smaller diffs = better review

```
   200-line PR   ██████████   attention spent: high   ✔
    20k-line PR  █            attention spent: gone    ✗
```

- A 200-line PR gets a real review. A 20k-line PR gets an LGTM.
- Bugs are caught where they're introduced, not archaeology later.

---

# 3 · Why it helps
## Unblock yourself

```
  sequential (no stacking):
    api      |=build=|=review=|merge|
    webhook                          |=build=|=review=|...
             └─ Alice blocked, waiting on review ─┘

  stacked:
    api      |=build=|=review=|merge|
    webhook     |=build=|=review=|merge|
    ui             |=build=|=review=|merge|
             └─ keep building up the stack, never idle ─┘
```

---

# 3 · Why it helps
## Throughput

```
  one big batch:
    ████████████████████  ──────────────▶  ships once, late

  small batches:
    ███ ▶ ███ ▶ ███ ▶ ███  ▶ ▶ ▶ ▶   ship continuously
```

Shorter review latency per PR → shorter cycle time. Flow, not heroics.

---

# 3 · Why it helps
## Continuous & agile

```
  a bug lands on main:
     ●──●──●──●──●──✗
  mega-merge :  one 20k commit   ▶  hunt inside it
  stacked    :  one small layer  ▶  revert / bisect fast  ✔
```

Small changes land **frequently**. Small blast radius = trivial to isolate.

---

# 3 · Why it helps
## Alice's week, restacked

```
  before :  1 × 20k PR    ▶  slow · shallow · blocked
  after  :  4 small PRs   ▶  fast · focused · flowing
```

- ✅ Early feedback on the API **before** the UI exists
- ✅ Bob reviews four focused PRs, actually catches bugs
- ✅ Alice keeps building up the stack, **never blocked**
- ✅ No hand-rebasing: the tool restacks for her

Same people. The **tooling** changed the outcome.

---

# 4 · gh-stack
## How do I actually do this?

Meet **`gh-stack`**, GitHub's native stacked-PR extension.

```
   local branches  ─▶  gh stack  ─▶  a stack of PRs on GitHub
```

(Public preview. We have early access.)

---

# 4 · gh-stack
## Install

```bash
gh extension install github/gh-stack
gh stack alias          # optional: adds `gs`
```

Now `gh stack …` (or `gs …`) manages your stack locally, then pushes to GitHub.

---

# 4 · gh-stack
## The five beats

```
  ┌────────┐  ┌──────┐  ┌────────┐  ┌────────┐  ┌───────┐
  │ create │─▶│ push │─▶│ submit │─▶│ review │─▶│ merge │
  └────────┘  └──────┘  └────────┘  └────────┘  └───────┘
   gs init     gs push   gs submit   (on GH)     gs merge
   gs add
```

Build layers locally, push, open the PRs, get focused reviews, land them.

---

# 4 · gh-stack
## Build the stack

```bash
gh stack init                       # start a stack on main
# ... edit the API layer ...
gh stack add -Am "api: BackupSchedule type" feature/backup-api
# ... edit the webhook layer ...
gh stack add -Am "webhook: validate schedule" feature/backup-webhook
```

`add` stages, commits, and stacks a new branch on top, all in one command.

---

# 4 · gh-stack
## Ship the stack

```bash
gh stack submit     # push all branches + open/update every PR
gh stack view       # see the stack locally
```

`submit` sets each PR's base to the layer below, automatically.

---

# 4 · gh-stack
## The stack map

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

# 4 · gh-stack
## Cascading merge

```
  merge the bottom PR...        ...the rest auto-rebase:

    ui         ●                  ui         ●
    controller ●                  controller ●  rebased
    webhook    ●         ─▶       webhook    ●  rebased
    api        ● ▸ main           api    ✔ merged
    main       ┴                  main       ┴
```

```bash
gh stack merge      # or click Merge on a PR in the UI
gh stack sync       # pull the cascade down locally
```

---

# 4 · gh-stack
## git → gs cheat sheet

| Goal | Manual git/gh | gh-stack |
|------|---------------|----------|
| Start bottom layer | `git checkout -b feat-api main` | `gs init` |
| Add dependent layer | `git checkout -b feat-webhook feat-api` | `gs add -Am "…" feat-webhook` |
| Push branches | `git push -u origin <each>` | `gs push` |
| Open PRs w/ right bases | open each PR, hand-set base | `gs submit` |
| Restack after change/merge | `git rebase` each branch in order, fix conflicts | `gs sync` (or `gs rebase`) |

**`gh-stack` is Option B, minus the manual pain.**

---

# 5 · Live demo
## Let's see it live 🔴

```
    ╔══════════════════════════════════╗
    ║             LIVE  DEMO           ║
    ╚══════════════════════════════════╝
       api   ▶   webhook   ▶   controller
```

A real 3-layer stack.

---

# 6 · Reality
## Review UX & CI

```
  each PR is tested as if it targets main:
    PR #3 controller   ▶  CI ✔
    PR #2 webhook      ▶  CI ✔
    PR #1 api          ▶  CI ✔
```

Branch protection is enforced on the **final target branch**, not the intermediate layer branches.

---

# 6 · Reality
## Gotchas

- Keep the stack **in sync**: `gs sync` after merges / review churn
- **Don't stack everything**: a tiny one-shot change is just a PR
- Stacking rewards **discipline**: small, coherent layers with clean boundaries

---

# 6 · Reality
## Stacks multiply CI load

```
  1 mega-PR       ▶  CI runs ×1
  4-layer stack   ▶  CI runs ×4      (one per layer)
  + a merge/restack  ▶  re-runs down the chain
  now multiply across every AI-authored PR, too ...
```

- Mitigate: skip redundant intermediate runs, path filters, concurrency-cancel, a merge queue, right-sized runners
- **Bigger picture:** CI must scale to changing throughput **regardless**. Stacking just surfaces the need sooner.

---

# 6 · Reality
## It's public preview, tell them

We're early adopters. Our feedback shapes the tool.

- `gh stack feedback`  →  `gh.io/stacks-feedback`
- Bugs / requests: **github/gh-stack** issues
- Docs: `gh.io/stacks`
- Internal: **#eng-gh-stack** to flag problems & get help

---

# 7 · Landscape
## You're not the first

The stacking landscape:

```
   Graphite    Aviator (av)    spr
   ghstack / Sapling (Meta)    git-town
```

Stacking is a proven workflow. gh-stack makes it **native to GitHub**.

---

# 7 · Landscape
## Why gh-stack for us

- **Native GitHub**: no third-party app, no extra permissions
- The **stack map** is built into the PR UI
- **CLI is optional**: works alongside the normal PR flow
- Zero new SaaS in the review path

A few of us are already stacking, and it works.

---

# Thanks

```
     ┌──────────────┐  small
     ├──────────────┤  reviewable
     ├──────────────┤  unblocked
     ├──────────────┤  continuous
     └──────┬───────┘
           main
```

`gh extension install github/gh-stack`

Questions?
