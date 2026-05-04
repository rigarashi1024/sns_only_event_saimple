---
name: fix-loop
description: Handle Gemini PR review feedback for this repo. Use when the user asks to fetch Gemini review comments, fix important findings, add non-urgent findings to TODO, respond to review comments, rerun Gemini review, or loop until critical issues are resolved.
---

# Fix Loop

Resolve only review feedback that should be fixed now. Track deferred feedback in `docs/GEMINI_REVIEW_TODO.md` or explain it in a PR comment.

## Inputs

- PR number from the user, or detect from current branch:
  `gh pr list --head "$(git branch --show-current)" --json number,title,url`
- Latest Gemini review comments:
  `gh pr view <PR> --json comments --jq '.comments[] | select(.body | startswith("### 🤖 Gemini PR Review")) | {createdAt, url, body}'`

## Classification

Classify each finding before editing:

- `must`: Clear bug/security/logic issue in the current patch that can break behavior, leak data, corrupt data, or block the intended workflow. Fix now.
- `todo`: Plausible concern, but not urgent for this PR or depends on future production/runtime assumptions. Add to `docs/GEMINI_REVIEW_TODO.md`.
- `comment`: False positive, already handled, or intentionally out of scope. Leave a PR comment explaining why no change is needed.

Do not fix style, naming, speculative, or TODO-managed issues just because Gemini says "security".

## Fix Workflow

1. Fetch the latest Gemini review.
2. Summarize findings and classification.
3. For `must` items:
   - Inspect the referenced code.
   - Make the smallest safe change.
   - Add or update focused tests when practical.
4. For `todo` items:
   - Append an entry to `docs/GEMINI_REVIEW_TODO.md` using its template.
   - Include PR number/comment URL, reason to defer, and revisit timing.
5. For `comment` items:
   - Post a concise PR comment explaining why no code change is needed.
6. Run relevant validation:
   - Backend Go: `env GOCACHE=/private/tmp/sns-only-event-go-cache GOMODCACHE=/Users/igarashiryuuta/go/pkg/mod go test ./...` from `apps/backend`.
   - Frontend: `npm run build` from `apps/frontend` when frontend changed.
   - Script-only JS: `node --check <file>`.
7. Commit and push only when the user asked to carry the loop through or when operating under this skill for an existing PR.
8. Trigger Gemini rerun if needed:
   - `gh pr comment <PR> --body "/gemini-review"`
9. Repeat at most 3 iterations. Stop earlier if:
   - validation fails and needs human judgment,
   - Gemini repeats a TODO/comment-only finding,
   - the remaining feedback is non-critical.

## PR Comment Guidance

Use comments to preserve the review flow. Do not edit old Gemini comments.

Comment format:

```md
確認しました。

- Finding ...: 今回は修正不要です。理由: ...
- Finding ...: `docs/GEMINI_REVIEW_TODO.md` に TODO として記録しました。理由: ...
```

## Safety Notes

- Do not merge PRs.
- Do not force-push unless the user explicitly asks.
- Do not delete branches.
- Do not revert unrelated user changes.
