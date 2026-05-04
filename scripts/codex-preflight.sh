#!/usr/bin/env bash
set -u

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="$ROOT_DIR/apps/frontend"
BACKEND_DIR="$ROOT_DIR/apps/backend"

overall=0

section() {
  printf '\n== %s ==\n' "$1"
}

has_script() {
  local dir="$1"
  local name="$2"
  node -e "const p=require(process.argv[1]); process.exit(p.scripts && p.scripts[process.argv[2]] ? 0 : 1)" "$dir/package.json" "$name" >/dev/null 2>&1
}

run_step() {
  local label="$1"
  shift
  printf '\n$ %s\n' "$*"
  if "$@"; then
    printf '[OK] %s\n' "$label"
    return 0
  else
    local code=$?
    printf '[NG] %s (exit %s)\n' "$label" "$code"
    return "$code"
  fi
}

run_optional_script() {
  local dir="$1"
  local label="$2"
  local script_name="$3"
  if has_script "$dir" "$script_name"; then
    (cd "$dir" && run_step "$label" npm run "$script_name") || overall=1
    return
  fi
  printf '[SKIP] %s: npm script "%s" not found\n' "$label" "$script_name"
}

section "Preflight"
printf 'root: %s\n' "$ROOT_DIR"

if [ -f "$FRONTEND_DIR/package.json" ]; then
  section "Frontend"

  if has_script "$FRONTEND_DIR" "format:check"; then
    if ! (cd "$FRONTEND_DIR" && run_step "frontend format:check" npm run format:check); then
      if has_script "$FRONTEND_DIR" "format"; then
        (cd "$FRONTEND_DIR" && run_step "frontend format" npm run format) || overall=1
        (cd "$FRONTEND_DIR" && run_step "frontend format:check after format" npm run format:check) || overall=1
      else
        overall=1
      fi
    fi
  else
    printf '[SKIP] frontend format:check: npm script not found\n'
  fi

  if has_script "$FRONTEND_DIR" "lint"; then
    if ! (cd "$FRONTEND_DIR" && run_step "frontend lint" npm run lint); then
      if has_script "$FRONTEND_DIR" "lint:fix"; then
        (cd "$FRONTEND_DIR" && run_step "frontend lint:fix" npm run lint:fix) || overall=1
        (cd "$FRONTEND_DIR" && run_step "frontend lint after lint:fix" npm run lint) || overall=1
      else
        overall=1
      fi
    fi
  else
    printf '[SKIP] frontend lint: npm script not found\n'
  fi

  if has_script "$FRONTEND_DIR" "typecheck"; then
    (cd "$FRONTEND_DIR" && run_step "frontend typecheck" npm run typecheck) || overall=1
  elif has_script "$FRONTEND_DIR" "type-check"; then
    (cd "$FRONTEND_DIR" && run_step "frontend type-check" npm run type-check) || overall=1
  else
    (cd "$FRONTEND_DIR" && run_step "frontend nuxi typecheck" npx nuxi typecheck) || overall=1
  fi

  run_optional_script "$FRONTEND_DIR" "frontend test" "test"
else
  printf '[SKIP] frontend: %s not found\n' "$FRONTEND_DIR/package.json"
fi

if [ -f "$BACKEND_DIR/go.mod" ]; then
  section "Backend"
  (
    cd "$BACKEND_DIR" &&
      run_step "backend go test" env GOCACHE=/private/tmp/sns-only-event-go-cache GOMODCACHE=/Users/igarashiryuuta/go/pkg/mod go test ./...
  ) || overall=1
else
  printf '[SKIP] backend: %s not found\n' "$BACKEND_DIR/go.mod"
fi

section "Result"
if [ "$overall" -eq 0 ]; then
  printf 'Preflight passed.\n'
else
  printf 'Preflight failed. See [NG] entries above.\n'
fi

exit "$overall"
