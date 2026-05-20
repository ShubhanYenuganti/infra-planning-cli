# Cross-CLI Compound Recipes

SQL snippets and shell pipelines that join data across multiple `infra-press`
SQLite stores. These are queries no individual vendor CLI can answer.

> **Note:** These recipes assume you've run `sync` on every relevant CLI first.
> Stores live at `~/.infra-press/<cli>.db`.

## Recipe template

```sql
-- Recipe name (one line)
-- What it answers (one sentence)
-- Required syncs: <list the CLIs whose stores must be synced>

ATTACH '~/.infra-press/cloud-run-admin-pp-cli.db' AS cr;
-- ... query ...
```

## Recipes (added as sprints ship)

### Stale services across all synced stores

_To be added once 2+ CLIs ship. See sprint v1 retro for the first instance._

---

(More recipes accumulate as sprints add CLIs.)
