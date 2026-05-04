---
name: fix-loop
description: Handle Gemini PR review feedback for this repo. Use when the user asks to fetch Gemini review comments, fix important findings, add non-urgent findings to TODO, respond to review comments, rerun Gemini review, or loop until critical issues are resolved.
---

# Fix Loop

Use scripts for GitHub fetch/comment mechanics and keep Codex reasoning only for classification.

Execution policy:

- Run `bash scripts/codex-fix-loop.sh ...` with escalation from the first attempt.
- The script calls GitHub APIs and posts PR comments, so a normal sandbox run may fail.
- The prefix is approved in Codex rules.

1. Fetch the latest Gemini review for the current branch PR:

```bash
bash scripts/codex-fix-loop.sh latest
```

Use `fetch` instead of `latest` only when the older review history is needed.

2. Classify each finding:

- `must`: clear bug/security/logic issue in the current patch; fix now.
- `todo`: plausible but not urgent; append to `docs/GEMINI_REVIEW_TODO.md`.
- `comment`: false positive, already handled, intentionally out of scope, or repeated comment-only finding.

3. Validate only what changed:

- backend: `bash scripts/codex-preflight.sh` when broad validation is useful, or targeted Go tests if the failure is narrow.
- frontend: `bash scripts/codex-preflight.sh` when frontend changed.
- script-only JS: `node --check <file>`.

4. Post comments through the script:

Human-facing comment, does not rerun Gemini and is not passed to Gemini:

```bash
bash scripts/codex-fix-loop.sh comment-human-text <PR> "確認しました。..."
```

Gemini-context comment, passes the explanation to the next Gemini prompt and triggers rerun:

```bash
bash scripts/codex-fix-loop.sh comment-gemini-text <PR> "確認しました。..."
```

Plain rerun without extra context:

```bash
bash scripts/codex-fix-loop.sh rerun <PR>
```

Use `comment-human` / `comment-gemini` with a file only when a multiline hand-written file is specifically needed. For generated comments, prefer the `*-text` commands so the script owns posting without extra temporary files.

Stop after at most 3 iterations, or earlier if Gemini repeats a TODO/comment-only finding.
