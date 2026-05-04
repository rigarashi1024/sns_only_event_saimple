---
name: preflight
description: Run the repository preflight checks with minimal Codex reasoning. Use before PR creation or when the user asks for quality checks.
disable-model-invocation: false
project: true
---

# Preflight

Run the scripted preflight:

```bash
bash scripts/codex-preflight.sh
```

The script owns the command flow and report:

- frontend format/lint/typecheck/test when scripts exist
- fallback `npx nuxi typecheck`
- backend `go test ./...`

Codex should only summarize the final `[OK]`, `[NG]`, and `[SKIP]` lines. If the script fails, inspect only the failing section instead of re-running each command manually.
