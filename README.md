# gh-stack-demo

Presentation + live demo introducing **GitHub Stacked PRs** (`gh-stack`) to engineering.

## Present

```bash
slides slides/stacked-prs.md
```

~30 min, 31 slides. The live demo runs after the "Let's see it live" slide.

## Demo

The `demo/` directory is a throwaway Kubernetes-operator scaffold (never
compiles) used to build a real 3-layer stack live. See
[`demo/README.md`](demo/README.md) for the runbook.

## Requirements

- [`slides`](https://github.com/maaslalani/slides) — terminal presentation tool
- `gh` + the stacked-PRs extension: `gh extension install github/gh-stack`
