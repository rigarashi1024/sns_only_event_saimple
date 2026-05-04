---
name: pr-creator
description: Create a GitHub pull request for this repo after preflight. Use when the user asks to create/open a PR, prepare a PR body, push a branch, or continue the standard flow after preflight.
---

# PR Creator

Create a PR with the repo's current review workflow in mind.

## Preconditions

- Prefer the flow: `preflight -> pr-creator -> fix-loop -> sync-main`.
- Check worktree first: `git status --short --branch`.
- Do not overwrite or revert user changes.
- If preflight was not run in this turn or there is no clear evidence it passed, either run the `preflight` skill first or create a Draft PR and say why.

## Workflow

1. Inspect current state:
   - `git branch --show-current`
   - `git status --short --branch`
   - `git diff --stat`
   - `git log --oneline main..HEAD`

2. Branch check:
   - Prefer `ai/<action>-<topic>` for new Codex-created work.
   - Accept existing user/feature branch names when already pushed or used by an open PR.
   - Do not rename a pushed branch without explicit user approval.

3. Commit if needed:
   - If there are uncommitted changes and the user asked to create the PR, prepare a focused commit.
   - Use concise commit messages such as `docs: update Gemini review workflow`.
   - Do not add Claude Code signatures or Anthropic co-author lines.
   - Do not run interactive rebase. If history cleanup is needed, explain options and ask.

4. Push branch:
   - Use `git push -u origin <branch>` for first push.
   - If network or permission fails, request escalation normally.

5. Build PR body:
   - Keep it concise and factual.
   - Include summary, changes, validation, and review points.
   - Mention Draft status when preflight failed, was skipped, or work is intentionally incomplete.

6. Create PR:
   - Use `gh pr create --base main --head <branch> --title <title> --body-file <tmpfile>`.
   - Add `--draft` when preflight is missing/failed or the work is WIP.
   - Report the PR URL and whether Gemini review should run automatically.

## PR Body Template

```md
## Summary

- ...

## Changes

- ...

## Validation

- [ ] `...`

## Review Notes

- ...
```

## Safety Notes

- Never push directly to `main`.
- Never merge the PR.
- Never delete branches unless the user explicitly asks.
- If an open PR already exists for the branch, update it instead of creating a duplicate.
