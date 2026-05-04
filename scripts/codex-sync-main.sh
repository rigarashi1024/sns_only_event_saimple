#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

printf '== Sync Main ==\n'
git status --short --branch

if [ -n "$(git status --porcelain)" ]; then
  printf '\n[STOP] Worktree has uncommitted changes. Commit, stash, or discard them before syncing main.\n' >&2
  exit 2
fi

current_branch="$(git branch --show-current)"
printf 'current branch: %s\n' "$current_branch"

git fetch origin
git switch main
git pull --ff-only origin main

printf '\n== Verify ==\n'
git status --short --branch

head_sha="$(git rev-parse HEAD)"
origin_sha="$(git rev-parse origin/main)"
main_sha="$(git rev-parse main)"

printf 'HEAD:        %s\n' "$head_sha"
printf 'origin/main: %s\n' "$origin_sha"
printf 'main:        %s\n' "$main_sha"

if [ "$head_sha" = "$origin_sha" ] && [ "$head_sha" = "$main_sha" ]; then
  printf 'main is synced with origin/main.\n'
else
  printf '[NG] main is not synced with origin/main.\n' >&2
  exit 1
fi

printf '\n== Recent Log ==\n'
git log --oneline --decorate --max-count=5
