---
description: Check out a PR, run the v1 verification harness, and post a ranked markdown report as a PR comment.
argument-hint: <PR-number>
---

# /verify-pr — run the verification harness against a PR

Use this command when reviewing a PR locally or when CI's verify workflow has run and you want a fuller human-readable report.

## What this command does

1. **Check out the PR locally** via `gh pr checkout $1`.
2. **Run `make verify`** from the repo root. This invokes L1 through L8 (L8 auto-skips clouds without OIDC creds in the local environment).
3. **Capture stdout to `verify.log`** in the repo root.
4. **Parse the log** into per-layer pass/fail/info counts.
5. **Render a markdown report** with these sections:
   - Header: `## Verification report — PR #$1 — commit <sha>`
   - Table: one row per layer (L1-L8) with status + counts.
   - L7 advisory subsection: rank "high" findings first, "info" findings collapsed under a `<details>` block.
   - L8 subsection: per-cloud SKIP / PASS / FAIL breakdown.
6. **Post the report** as a PR comment using a stable marker so re-runs update in place:

```bash
gh pr comment "$1" --body-file /tmp/verify-pr-report.md \
  --edit-last  # if a prior comment with the marker exists, update it
```

Use the marker `<!-- verify-pr -->` at the top of every report so subsequent runs can detect and replace the previous comment.

## Important caveats

- Running this locally requires `make`, `go`, `yq`, and (for Azure CLIs to build) network access to apis.guru. Confirm via `which make go yq` before invoking.
- L8 requires real cloud credentials in the environment (`aws sts get-caller-identity`, `gcloud auth list`, `az account show`). Without them, L8 SKIPs every cloud — that's expected and not a failure.
- The command will checkout the PR branch — make sure your working tree is clean first (`git status` shows clean) or stash before invoking.

## Output format

The posted comment should look like:

```markdown
<!-- verify-pr -->
## Verification report — PR #123 — commit abc1234

| Layer | Status | Detail |
|---|---|---|
| L1 Lint | PASS 30/30 | |
| L2 Build | PASS 6/6 | |
| L3 Golden | PASS 24/24 | |
| L4 Smoke | PASS 30/30 | |
| L5 Recipe | PASS | doctor-all.sh: 6 CLIs reported |
| L6 Spec sync | PASS 6/6 | |
| L7 Upstream | INFO | 2 high-priority findings |
| L8 Live cloud | SKIP | OIDC not configured for AWS/GCP/Azure |

### L7 advisory findings (ranked)

**High priority:**
- `lambda-pp-cli`: operation removed upstream: GET /functions::ListFunctions
- `cloud-run-admin-pp-cli`: arg `pageSize` changed from optional to required

<details>
<summary>Informational (3 findings)</summary>

- `cloud-functions-pp-cli`: 5 new endpoints available upstream (not exposed)
- ...
</details>

**Result: PASS** (L1–L7 green; L8 not configured)
```
