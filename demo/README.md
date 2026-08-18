# Live demo runbook — a 3-layer stack

Builds `api → webhook → controller` as a real stack of PRs on
`ahmad-ibra/gh-stack-demo`, shows the stack map, does a cascading merge,
and (optionally) resolves an engineered conflict.

## Pre-flight (before the talk)
- `gh auth status` → logged in, `repo` + `workflow` scopes
- `gh extension list` → `github/gh-stack` present; `gh stack alias` run (so `gs` works)
- `git switch main && git pull` — clean tree
- No leftover demo branches: `git branch --list 'feature/backup-*'` is empty
  (clean up with `gh stack unstack` + `git branch -D`, and delete remote branches/PRs)

## Core demo (~5–7 min)

```bash
# 0. start on a clean main
git switch main && git pull

# 1. start the stack
gs init

# 2. Layer 1 — API
#    edit demo/api/v1/backupschedule_types.go (add a field under the TODO)
gs add -Am "api: add retention to BackupSchedule" feature/backup-api

# 3. Layer 2 — webhook  (also append to the shared registry -> sets up the conflict)
#    edit demo/internal/webhook/backupschedule_webhook.go (fill in a check)
#    edit demo/internal/features/registry.go: add "backup-webhook" to registeredFeatures
gs add -Am "webhook: validate cron schedule" feature/backup-webhook

# 4. Layer 3 — controller  (also append to the shared registry -> the conflict)
#    edit demo/internal/controller/backupschedule_controller.go (fill in reconcile)
#    edit demo/internal/features/registry.go: add "backup-controller" to registeredFeatures
gs add -Am "controller: reconcile backup CronJob" feature/backup-controller

# 5. ship it
gs submit
gs view
```

Then in the browser: open the bottom PR, **show the stack map**, walk up the layers.

```bash
# 6. cascading merge: merge the bottom PR (UI or CLI), then restack locally
gs merge          # or click Merge on PR #1 in the UI
gs sync           # pulls the cascade; upper layers rebase onto main
```

## Optional: conflict beat (~2 min, only if time)
Because layers 2 and 3 both appended to `registeredFeatures`, a `gs sync` /
`gs rebase` after editing the base will stop on a conflict in
`demo/internal/features/registry.go`. Resolve it live:

```bash
# git shows the conflict in registry.go — keep both feature keys
gs rebase --continue   # (or: git add + git rebase --continue, then gs sync)
```

## Fallback (if live push/auth fails on stage)
A pre-created stack of PRs is left open on the repo (label `demo-fallback`).
Switch to the browser and walk that stack + stack map instead of building live.

## Teardown (after the talk)
```bash
gs unstack                      # remove local stack tracking
git switch main
git branch -D feature/backup-api feature/backup-webhook feature/backup-controller
# close/delete the demo PRs + remote branches on github.com (or `gh stack unstack` handles remote)
git switch main && git pull
```
