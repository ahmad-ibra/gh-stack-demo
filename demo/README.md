# Live demo runbook — a 3-layer stack

Builds `api → webhook → controller` as a real stack of PRs on
`ahmad-ibra/gh-stack-demo`, shows the stack map, does a cascading merge,
and (optionally) resolves an engineered conflict.

The `demo/` module is a real Kubernetes operator (a `BackupSchedule`
controller). `go build ./...` passes, so every layer is a believable diff.

## Pre-flight (before the talk)
- `gh auth status` → logged in, `repo` + `workflow` scopes
- `gh extension list` → `github/gh-stack` present; `gh stack alias` run (so `gs` works)
- `git switch main && git pull` — clean tree
- No leftover demo branches: `git branch --list 'feat/*'` is empty
  (clean up with `gh stack unstack` + `git branch -D`, and delete remote branches/PRs)

## Core demo (~5–7 min)

```bash
# 0. start on a clean main
git switch main && git pull

# 1. Layer 1 — API
gs init feat/api
#    edit demo/api/v1/backupschedule_types.go: add a TimeZone field to BackupScheduleSpec
git add . && git commit -m "api: add TimeZone to BackupSchedule"

# 2. Layer 2 — webhook
gs add feat/webhook
#    edit demo/internal/webhook/backupschedule_webhook.go: validate the TimeZone
#    edit demo/internal/features/registry.go: append a Feature to `registered`  (sets up the conflict)
git add . && git commit -m "webhook: validate TimeZone"

# 3. Layer 3 — controller
gs add feat/controller
#    edit demo/internal/controller/backupschedule_controller.go: set CronJob.Spec.TimeZone
#    edit demo/internal/features/registry.go: append a Feature to `registered`  (the conflict)
git add . && git commit -m "controller: set CronJob time zone"

# 4. ship it
gs submit
gs view
```

Then in the browser: open the bottom PR, **show the stack map**, walk up the layers.

```bash
# 5. cascading merge: merge the bottom PR (UI or CLI), then restack locally
gs merge          # or click Merge on PR #1 in the UI
gs sync           # pulls the cascade; upper layers rebase onto main
```

## Optional: conflict beat (~2 min, only if time)
Layers 2 and 3 both appended to the `registered` slice in
`demo/internal/features/registry.go`, so a `gs sync` / `gs rebase` after the
base changes stops on a conflict there. Resolve it live (keep both features):

```bash
# git shows the conflict in registry.go — keep both Feature entries
gs rebase --continue   # (or: git add + git rebase --continue, then gs sync)
```

## Fallback (if live push/auth fails on stage)
A pre-created stack of PRs is left open on the repo (label `demo-fallback`).
Switch to the browser and walk that stack + stack map instead of building live.

## Teardown (after the talk)
```bash
gs unstack                      # remove local stack tracking
git switch main
git branch -D feat/api feat/webhook feat/controller
# close/delete the demo PRs + remote branches on github.com (or `gh stack unstack` handles remote)
git switch main && git pull
```
