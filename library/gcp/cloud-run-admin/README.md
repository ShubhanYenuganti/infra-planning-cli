# Cloud Run Admin CLI

Deploy and manage user provided container images that scale automatically based on incoming requests. The Cloud Run Admin API v1 follows the Knative Serving API specification, while v2 is aligned with Google Cloud AIP-based API standards, as described in https://google.aip.dev/.

Learn more at [Cloud Run Admin](https://google.com).

## Install

### Go

```
go install github.com/mvanhorn/printing-press-library/library/other/cloud-run-admin-pp-cli/cmd/cloud-run-admin-pp-cli@latest
```

### Binary

Download from [Releases](https://github.com/mvanhorn/printing-press-library/releases).

## Quick Start

### 1. Install

See [Install](#install) above.

### 2. Set Up Credentials

Get your access token from your API provider's developer portal, then store it:

```bash
cloud-run-admin-pp-cli auth set-token YOUR_TOKEN_HERE
```

Or set it via environment variable:

```bash
export CLOUD_RUN_ADMIN_TOKEN="your-token-here"
```

### 3. Verify Setup

```bash
cloud-run-admin-pp-cli doctor
```

This checks your configuration and credentials.

### 4. Try Your First Command

```bash
cloud-run-admin-pp-cli name-run list
```

## Usage

<!-- HELP_OUTPUT -->

## Commands

### name-run

Manage name run

- **`cloud-run-admin-pp-cli name-run run-projects-locations-jobs-run`** - Triggers creation of a new Execution of this Job.

### name-wait

Manage name wait

- **`cloud-run-admin-pp-cli name-wait run-projects-locations-operations-wait`** - Waits until the specified long-running operation is done or reaches at most a specified timeout, returning the latest state. If the operation is already done, the latest state is immediately returned. If the timeout specified is greater than the default HTTP/RPC timeout, the HTTP/RPC timeout is used. If the server does not support this method, it returns `google.rpc.Code.UNIMPLEMENTED`. Note that this method is on a best-effort basis. It may return the latest state before the specified timeout (including immediately), meaning even an immediate response is no guarantee that the operation is done.

### resource-get-iam-policy

Manage resource get iam policy

- **`cloud-run-admin-pp-cli resource-get-iam-policy run-projects-locations-services-get-iam-policy`** - Gets the IAM Access Control policy currently in effect for the given Cloud Run Service. This result does not include any inherited policies.

### resource-set-iam-policy

Manage resource set iam policy

- **`cloud-run-admin-pp-cli resource-set-iam-policy run-projects-locations-services-set-iam-policy`** - Sets the IAM Access control policy for the specified Service. Overwrites any existing policy.

### resource-test-iam-permissions

Manage resource test iam permissions

- **`cloud-run-admin-pp-cli resource-test-iam-permissions run-projects-locations-services-test-iam-permissions`** - Returns permissions that a caller has on the specified Project. There are no permissions required for making this API call.


## Output Formats

```bash
# Human-readable table (default in terminal, JSON when piped)
cloud-run-admin-pp-cli name-run list

# JSON for scripting and agents
cloud-run-admin-pp-cli name-run list --json

# Filter to specific fields
cloud-run-admin-pp-cli name-run list --json --select id,name,status

# Dry run — show the request without sending
cloud-run-admin-pp-cli name-run list --dry-run

# Agent mode — JSON + compact + no prompts in one flag
cloud-run-admin-pp-cli name-run list --agent
```

## Agent Usage

This CLI is designed for AI agent consumption:

- **Non-interactive** - never prompts, every input is a flag
- **Pipeable** - `--json` output to stdout, errors to stderr
- **Filterable** - `--select id,name` returns only fields you need
- **Previewable** - `--dry-run` shows the request without sending
- **Retryable** - creates return "already exists" on retry, deletes return "already deleted"
- **Confirmable** - `--yes` for explicit confirmation of destructive actions
- **Piped input** - `echo '{"key":"value"}' | cloud-run-admin-pp-cli <resource> create --stdin`
- **Cacheable** - GET responses cached for 5 minutes, bypass with `--no-cache`
- **Agent-safe by default** - no colors or formatting unless `--human-friendly` is set
- **Progress events** - paginated commands emit NDJSON events to stderr in default mode

Exit codes: `0` success, `2` usage error, `3` not found, `4` auth error, `5` API error, `7` rate limited, `10` config error.

## Use as MCP Server

This CLI ships a companion MCP server for use with Claude Desktop, Cursor, and other MCP-compatible tools.

### Claude Code

```bash
claude mcp add cloud-run-admin cloud-run-admin-pp-mcp -e CLOUD_RUN_ADMIN_TOKEN=<your-token>
```

### Claude Desktop

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "cloud-run-admin": {
      "command": "cloud-run-admin-pp-mcp",
      "env": {
        "CLOUD_RUN_ADMIN_TOKEN": "<your-key>"
      }
    }
  }
}
```

## Cookbook

Common workflows and recipes:

```bash
# List resources as JSON for scripting
cloud-run-admin-pp-cli name-run list --json

# Filter to specific fields
cloud-run-admin-pp-cli name-run list --json --select id,name,status

# Dry run to preview the request
cloud-run-admin-pp-cli name-run list --dry-run

# Sync data locally for offline search
cloud-run-admin-pp-cli sync

# Search synced data
cloud-run-admin-pp-cli search "query"

# Export for backup
cloud-run-admin-pp-cli export --format jsonl > backup.jsonl
```

## Health Check

```bash
cloud-run-admin-pp-cli doctor
```

<!-- DOCTOR_OUTPUT -->

## Configuration

Config file: `~/.config/cloud-run-admin-pp-cli/config.toml`

Environment variables:
- `CLOUD_RUN_ADMIN_TOKEN`

## Troubleshooting

**Authentication errors (exit code 4)**
- Run `cloud-run-admin-pp-cli doctor` to check credentials
- Verify the environment variable is set: `echo $CLOUD_RUN_ADMIN_TOKEN`

**Not found errors (exit code 3)**
- Check the resource ID is correct
- Run the `list` command to see available items

**Rate limit errors (exit code 7)**
- The CLI auto-retries with exponential backoff
- If persistent, wait a few minutes and try again

---

Generated by [CLI Printing Press](https://github.com/mvanhorn/cli-printing-press)
