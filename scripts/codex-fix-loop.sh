#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

BOT_HEADER='### 🤖 Gemini PR Review'
CONTEXT_MARKER='<!-- gemini-review-context -->'

usage() {
  cat <<'USAGE'
Usage:
  scripts/codex-fix-loop.sh fetch [PR_NUMBER]
  scripts/codex-fix-loop.sh latest [PR_NUMBER]
  scripts/codex-fix-loop.sh comment-human PR_NUMBER COMMENT_FILE
  scripts/codex-fix-loop.sh comment-gemini PR_NUMBER COMMENT_FILE
  scripts/codex-fix-loop.sh comment-human-text PR_NUMBER COMMENT_TEXT
  scripts/codex-fix-loop.sh comment-gemini-text PR_NUMBER COMMENT_TEXT
  scripts/codex-fix-loop.sh rerun PR_NUMBER

fetch/latest detect PR from the current branch when PR_NUMBER is omitted.
comment-gemini prepends the Gemini context marker and /gemini-review trigger.
Use "-" as COMMENT_TEXT to read the body from stdin.
USAGE
}

detect_pr() {
  local branch
  branch="$(git branch --show-current)"
  gh pr list --head "$branch" --json number --jq '.[0].number'
}

pr_number_or_detect() {
  local pr="${1:-}"
  if [ -z "$pr" ]; then
    pr="$(detect_pr)"
  fi
  if [ -z "$pr" ] || [ "$pr" = "null" ]; then
    printf '[STOP] PR number not found. Pass PR_NUMBER explicitly.\n' >&2
    exit 2
  fi
  printf '%s\n' "$pr"
}

fetch_reviews() {
  local pr
  pr="$(pr_number_or_detect "${1:-}")"
  gh pr view "$pr" --json comments --jq \
    '.comments[] | select(.body | startswith("'"$BOT_HEADER"'")) | {createdAt, url, body}'
}

latest_review() {
  local pr
  pr="$(pr_number_or_detect "${1:-}")"
  gh pr view "$pr" --json comments --jq \
    '[.comments[] | select(.body | startswith("'"$BOT_HEADER"'"))] | last | {createdAt, url, body}'
}

comment_human() {
  local pr="${1:-}"
  local file="${2:-}"
  if [ -z "$pr" ] || [ -z "$file" ]; then
    usage >&2
    exit 2
  fi
  gh pr comment "$pr" --body-file "$file"
}

read_body_arg_or_stdin() {
  local body="${1:-}"
  if [ "$body" = "-" ]; then
    cat
    return
  fi
  printf '%s\n' "$body"
}

comment_human_text() {
  local pr="${1:-}"
  local body_arg="${2:-}"
  if [ -z "$pr" ] || [ -z "$body_arg" ]; then
    usage >&2
    exit 2
  fi

  local body
  body="$(read_body_arg_or_stdin "$body_arg")"
  gh pr comment "$pr" --body "$body"
}

comment_gemini() {
  local pr="${1:-}"
  local file="${2:-}"
  if [ -z "$pr" ] || [ -z "$file" ]; then
    usage >&2
    exit 2
  fi

  local tmp
  tmp="$(mktemp)"
  {
    printf '%s\n' "$CONTEXT_MARKER"
    printf '/gemini-review\n\n'
    cat "$file"
  } > "$tmp"
  gh pr comment "$pr" --body-file "$tmp"
  rm -f "$tmp"
}

comment_gemini_text() {
  local pr="${1:-}"
  local body_arg="${2:-}"
  if [ -z "$pr" ] || [ -z "$body_arg" ]; then
    usage >&2
    exit 2
  fi

  local body
  body="$(read_body_arg_or_stdin "$body_arg")"
  gh pr comment "$pr" --body "$(printf '%s\n/gemini-review\n\n%s' "$CONTEXT_MARKER" "$body")"
}

rerun() {
  local pr="${1:-}"
  if [ -z "$pr" ]; then
    usage >&2
    exit 2
  fi
  gh pr comment "$pr" --body '/gemini-review'
}

case "${1:-}" in
  fetch) shift; fetch_reviews "${1:-}" ;;
  latest) shift; latest_review "${1:-}" ;;
  comment-human) shift; comment_human "$@" ;;
  comment-gemini) shift; comment_gemini "$@" ;;
  comment-human-text) shift; comment_human_text "$@" ;;
  comment-gemini-text) shift; comment_gemini_text "$@" ;;
  rerun) shift; rerun "$@" ;;
  -h|--help|help) usage ;;
  *) usage >&2; exit 2 ;;
esac
