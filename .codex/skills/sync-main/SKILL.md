---
name: sync-main
description: Sync the local repository with the latest merged main branch. Use after a PR is merged, when the user says to reflect the latest main, update from main, or switch back to main after work is done.
---

# Sync Main

Run the scripted sync:

```bash
bash scripts/codex-sync-main.sh
```

The script checks for a dirty worktree, fetches `origin`, switches to `main`, fast-forwards with `--ff-only`, and prints verification.

Codex should not manually run the individual git commands unless the script fails and the failure needs diagnosis.
