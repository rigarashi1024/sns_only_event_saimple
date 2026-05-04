---
name: pr-creator
description: Create a GitHub pull request for this repo after preflight. Use when the user asks to create/open a PR, prepare a PR body, push a branch, or continue the standard flow after preflight.
---

# PR Creator

Use scripts for the mechanical work.

Execution policy:

- Run `bash scripts/codex-pr-creator.sh ...` normally first.
- Escalate only if it fails due to sandbox, network, GitHub API, or `.git` write permissions.
- The prefix is approved in Codex rules for escalation when needed.

1. Inspect current branch and PR state:

```bash
bash scripts/codex-pr-creator.sh inspect
```

2. Run preflight when there is no clear evidence it passed:

```bash
bash scripts/codex-preflight.sh
```

3. Let Codex draft only the PR title/body and optional commit message.

4. Prefer the auto mode so the script owns commit, push, PR body temp file creation, cleanup, and existing PR detection:

```bash
bash scripts/codex-pr-creator.sh create-auto "PR title" "commit message"
```

Append `--draft` when preflight failed, was skipped, or the work is intentionally incomplete.

Use `create "PR title" PR_BODY_FILE "commit message"` only when a hand-written PR body file is specifically needed.

Safety:

- Never create a PR from `main`.
- Do not force-push, merge, or delete branches.
- If the script stops on a dirty worktree and no `COMMIT_MESSAGE` was provided, decide whether a commit is appropriate before rerunning.
