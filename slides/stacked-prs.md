---
author: Ahmad Ibrahim · Stacked PRs with gh stack
paging: Slide %d / %d
---

# Stacked PRs with gh stack

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
  1 · problem         a typical dev's workflow
  2 · a better way    what "stacking" is
  3 · gh stack        the commands
  4 · live demo       stacking changes on each other
  5 · why it matters  throughput & traceability
  6 · reality         gotchas & CI cost
  7 · landscape       alternatives & why gh stack
```

---

# 1 · The problem
## The setup

Alice has a big feature to implement. It has multiple pieces that depend on each other:

```
  ┌───────┐   ┌─────────┐   ┌────────────┐   ┌──────┐
  │  api  │ ▶ │ webhook │ ▶ │ controller │ ▶ │  ui  │
  └───────┘   └─────────┘   └────────────┘   └──────┘
```

Alice hits a fork: **how does she ship this?**

Two options, neither are ideal.

---

# 1 · The problem
## Option A: the mega-PR

Everything crammed into one PR:

```
   ┌───────┬─────────┬────────────┬──────┐
   │  api  │ webhook │ controller │  ui  │
   └───────┴─────────┴────────────┴──────┘
                     ▼
          +20,000 lines, one giant PR to review
```

- **One giant diff**: difficult to give it a careful review
- **All-or-nothing**: nothing lands until the whole PR is approved
- **The branch rots**: while it waits for review, `main` moves on and conflicts build up

---

# 1 · The problem
## Option A: the mega-PR

Bob gets the review notification and opens a **+20,000-line** PR.

```
  defects caught in review, by lines changed
   1-100     ████████████████████  87%
   101-300   ██████████████████    78%
   301-600   ███████████████       65%
   601-1000  ██████████            42%
   1000+     ███████               28%   ◀ a mega-PR lands here
```

By end of day, Bob does one of two things:

- **Rubber-stamps it**: a quick LGTM, and bugs slip through
- **Grinds through it**: hours of review, but the size makes it shallow and inefficient anyway

Source: Propel Code, 50k+ PRs: https://www.propelcode.ai/blog/pr-size-impact-code-review-quality-data-study

---

# 1 · The problem
## Option B: smaller PRs

The right instinct: split it into small PRs. But each depends on the one
below, so they land one at a time:

```
  main
   ┃
   ┣━▶ branch → code → review ⌛ → merge ┐
   ┃ ◀──────────────────────────────────┘
   ┣━▶ branch → code → review ⌛ → merge ┐
   ┃ ◀──────────────────────────────────┘
   ┣━▶ branch → code → review ⌛ → merge ┐
   ┃ ◀──────────────────────────────────┘
   ▼
  time
```

Alice's idle time continues to stack as each change is reviewed. To keep
moving she could branch off her own unmerged branches, but then she owns
every rebase by hand.

---

# 1 · The problem
## The false dilemma

Two ways to ship the feature, and neither are great:

```
       Option A                      Option B
       the mega-PR                   smaller PRs
           │                             │
           ▼                             ▼
      hard to review,               lots of idle time,
      bugs slip through             or manual rebasing
           │                             │
           ▼                             ▼
           ✗                             ✗
```

Bad tooling is forcing this choice.

---

# 2 · A better way
## What if you didn't have to choose?

Take the good from each, leave the bad:

```
   keep building,               small, focused
   never blocked                reviews
          └───────────┬──────────┘
                      ▼
                  a  STACK
```

---

# 2 · A better way
## What a stack is

Two or more PRs: the bottom targets `main`, and each one above targets the
PR below it.

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

Each PR is a small, reviewable change. Anything an upper layer needs lives
in a lower one.

---

# 2 · A better way
## Meet gh stack

`gh stack` is a subcommand of the GitHub CLI (added by the `gh-stack`
extension) that builds and manages stacks for you:

- Create & manage stacks: add layers, reorder, navigate
- Rebasing & syncing handled: change a lower layer, it restacks the rest
- A stack map in the GitHub UI: every PR linked to the others
- Cascading merges: land the whole stack, or just the bottom part

The manual work from Option B, gone.

---

# 3 · gh stack
## Install

```bash
gh extension install github/gh-stack
gh stack alias          # optional: adds `gs`
```

Now `gs` manages your stack locally, then pushes to GitHub. (We'll use `gs` from here on.)

---

# 3 · gh stack
## Build the stack

```bash
gs init                          # start a stack on main

# edit the api layer, then add it as the bottom PR:
gs add -A -m "api: BackupSchedule type" feature/backup-api
#   -A stages your edits · -m is the commit message · then it stacks the branch

# edit the webhook layer, then stack it on top:
gs add -A -m "webhook: validate schedule" feature/backup-webhook
```

One `add` stages, commits, and creates the next layer.

---

# 3 · gh stack
## Ship the stack

```bash
gs submit     # push all branches + open/update every PR
gs view       # see the stack locally
```

`submit` sets each PR's base to the layer below, automatically.

---

# 3 · gh stack
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
gs merge      # or click Merge on a PR in the UI
gs sync       # pull the cascade down locally
```

---

# 3 · gh stack
## git → gs cheat sheet

| the git way | with gh stack |
|-------------|---------------|
| `git checkout -b feat-api main` | `gs init` |
| edit, `git commit`, `git checkout -b feat-webhook` | `gs add -A -m "…" feat-webhook` |
| `git push -u origin` each branch | `gs push` |
| open each PR, set its base by hand | `gs submit` |
| `git rebase` every branch in order, fix conflicts | `gs sync` |

**`gh stack` is Option B, minus the manual pain.**

---

# 4 · Live demo
## Let's see it live 🔴

```
    ╔══════════════════════════════════╗
    ║             LIVE  DEMO           ║
    ╚══════════════════════════════════╝
       api   ▶   webhook   ▶   controller
```

A real 3-layer stack.

---

# 5 · Why it matters
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

# 5 · Why it matters
## Throughput

```
  one big batch:
    ████████████████████  ──────────────▶  ships once, late

  small batches:
    ███ ▶ ███ ▶ ███ ▶ ███  ▶ ▶ ▶ ▶   ship continuously
```

Shorter review latency per PR → shorter cycle time. Flow, not heroics.

---

# 5 · Why it matters
## Continuous & agile

```
  a bug lands on main:
     ●──●──●──●──●──✗
  mega-merge :  one 20k commit   ▶  hunt inside it
  stacked    :  one small layer  ▶  revert / bisect fast  ✔
```

Small changes land **frequently**. Small blast radius = trivial to isolate.

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

- `gs feedback`  →  `gh.io/stacks-feedback`
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

Stacking is a proven workflow. gh stack makes it **native to GitHub**.

---

# 7 · Landscape
## Why gh stack for us

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
