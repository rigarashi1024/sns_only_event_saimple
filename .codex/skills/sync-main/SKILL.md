---
name: sync-main
description: Sync the local repository with the latest merged main branch. Use after a PR is merged, when the user says to reflect the latest main, update from main, or switch back to main after work is done.
---

# Sync Main

Update local `main` after a PR is merged.

## Workflow

1. Check current state:
   - `git status --short --branch`
   - `git branch --show-current`
2. If there are uncommitted changes:
   - Stop and ask how to handle them.
   - Do not stash, commit, discard, or switch branches without user approval.
3. Fetch remote:
   - `git fetch origin`
4. Switch to main:
   - `git switch main`
5. Fast-forward only:
   - `git pull --ff-only origin main`
6. Verify:
   - `git status --short --branch`
   - `git rev-parse HEAD origin/main main`
   - `git log --oneline --decorate --max-count=5`
7. Report:
   - current branch,
   - current HEAD,
   - whether `main` equals `origin/main`,
   - whether the worktree is clean.

## Safety Notes

- Prefer `git switch` over `git checkout`.
- Use `--ff-only`; do not create merge commits during sync.
- Never run `git reset --hard` unless explicitly requested.
- If a merged feature branch remains locally, leave it alone unless the user asks to delete it.
