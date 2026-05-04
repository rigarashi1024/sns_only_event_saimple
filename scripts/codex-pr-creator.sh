#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'USAGE'
Usage:
  scripts/codex-pr-creator.sh inspect
  scripts/codex-pr-creator.sh create "PR title" PR_BODY_FILE [COMMIT_MESSAGE] [--draft]

Environment fallback for create:
  PR_TITLE
  PR_BODY_FILE
  COMMIT_MESSAGE
  DRAFT=1
USAGE
}

current_branch() {
  git branch --show-current
}

inspect() {
  local branch
  branch="$(current_branch)"

  printf '== PR Creator Inspect ==\n'
  printf 'branch: %s\n\n' "$branch"

  printf '== Status ==\n'
  git status --short --branch

  printf '\n== Diff Stat ==\n'
  git diff --stat || true

  printf '\n== Commits ahead of main ==\n'
  git log --oneline main..HEAD || true

  printf '\n== Existing PRs for branch ==\n'
  gh pr list --head "$branch" --json number,title,url,state,isDraft
}

create() {
  local branch
  branch="$(current_branch)"
  local title="${PR_TITLE:-${1:-}}"
  local body_file="${PR_BODY_FILE:-${2:-}}"
  local commit_message="${COMMIT_MESSAGE:-${3:-}}"
  local draft="${DRAFT:-0}"

  for arg in "$@"; do
    if [ "$arg" = "--draft" ]; then
      draft=1
    fi
  done

  if [ "$branch" = "main" ]; then
    printf '[STOP] Refusing to create a PR directly from main.\n' >&2
    exit 2
  fi

  if [ -n "$commit_message" ] && [ -n "$(git status --porcelain)" ]; then
    git add .
    git commit -m "$commit_message"
  fi

  if [ -n "$(git status --porcelain)" ]; then
    printf '[STOP] Worktree has uncommitted changes. Set COMMIT_MESSAGE or commit manually.\n' >&2
    exit 2
  fi

  if [ -z "$title" ] || [ -z "$body_file" ]; then
    usage >&2
    exit 2
  fi

  git push -u origin "$branch"

  if gh pr list --head "$branch" --json number --jq 'length' | grep -qx '0'; then
    args=(--base main --head "$branch" --title "$title" --body-file "$body_file")
    if [ "$draft" = "1" ]; then
      args+=(--draft)
    fi
    gh pr create "${args[@]}"
  else
    printf 'Existing PR found for %s:\n' "$branch"
    gh pr list --head "$branch" --json number,title,url,state,isDraft
  fi
}

case "${1:-}" in
  inspect) inspect ;;
  create) shift; create "$@" ;;
  -h|--help|help) usage ;;
  *) usage >&2; exit 2 ;;
esac
