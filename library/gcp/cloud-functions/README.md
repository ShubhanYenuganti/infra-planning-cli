# Cloud Functions CLI

Manages lightweight user-provided functions executed in response to events.

Learn more at [Cloud Functions](https://google.com).

## Install

### Go

```
go install github.com/mvanhorn/printing-press-library/library/other/cloud-functions-pp-cli/cmd/cloud-functions-pp-cli@latest
```

### Binary

Download from [Releases](https://github.com/mvanhorn/printing-press-library/releases).

## Quick Start

### 1. Install

See [Install](#install) above.

### 2. Set Up Credentials

Get your access token from your API provider's developer portal, then store it:

```bash
cloud-functions-pp-cli auth set-token YOUR_TOKEN_HERE
```

Or set it via environment variable:

```bash
export CLOUD_FUNCTIONS_TOKEN="your-token-here"
```

### 3. Verify Setup

```bash
cloud-functions-pp-cli doctor
```

This checks your configuration and credentials.

### 4. Try Your First Command

```bash
cloud-functions-pp-cli name-generate-download-url list
```

## Usage

<!-- HELP_OUTPUT -->

## Commands

### name-generate-download-url

Manage name generate download url

- **`cloud-functions-pp-cli name-generate-download-url cloudfunctions-projects-locations-functions-generate-download-url`** - Returns a signed URL for downloading deployed function source code. The URL is only valid for a limited period and should be used within 30 minutes of generation. For more information about the signed URL usage see: https://cloud.google.com/storage/docs/access-control/signed-urls

### resource-get-iam-policy

Manage resource get iam policy

- **`cloud-functions-pp-cli resource-get-iam-policy cloudfunctions-projects-locations-functions-get-iam-policy`** - Gets the access control policy for a resource. Returns an empty policy if the resource exists and does not have a policy set.

### resource-set-iam-policy

Manage resource set iam policy

- **`cloud-functions-pp-cli resource-set-iam-policy cloudfunctions-projects-locations-functions-set-iam-policy`** - Sets the access control policy on the specified resource. Replaces any existing policy. Can return `NOT_FOUND`, `INVALID_ARGUMENT`, and `PERMISSION_DENIED` errors.

### resource-test-iam-permissions

Manage resource test iam permissions

- **`cloud-functions-pp-cli resource-test-iam-permissions cloudfunctions-projects-locations-functions-test-iam-permissions`** - Returns permissions that a caller has on the specified resource. If the resource does not exist, this will return an empty set of permissions, not a `NOT_FOUND` error. Note: This operation is designed to be used for building permission-aware UIs and command-line tools, not for authorization checking. This operation may "fail open" without warning.


## Output Formats

```bash
# Human-readable table (default in terminal, JSON when piped)
cloud-functions-pp-cli name-generate-download-url list

# JSON for scripting and agents
cloud-functions-pp-cli name-generate-download-url list --json

# Filter to specific fields
cloud-functions-pp-cli name-generate-download-url list --json --select id,name,status

# Dry run — show the request without sending
cloud-functions-pp-cli name-generate-download-url list --dry-run

# Agent mode — JSON + compact + no prompts in one flag
cloud-functions-pp-cli name-generate-download-url list --agent
```

## Agent Usage

This CLI is designed for AI agent consumption:

- **Non-interactive** - never prompts, every input is a flag
- **Pipeable** - `--json` output to stdout, errors to stderr
- **Filterable** - `--select id,name` returns only fields you need
- **Previewable** - `--dry-run` shows the request without sending
- **Retryable** - creates return "already exists" on retry, deletes return "already deleted"
- **Confirmable** - `--yes` for explicit confirmation of destructive actions
- **Piped input** - `echo '{"key":"value"}' | cloud-functions-pp-cli <resource> create --stdin`
- **Cacheable** - GET responses cached for 5 minutes, bypass with `--no-cache`
- **Agent-safe by default** - no colors or formatting unless `--human-friendly` is set
- **Progress events** - paginated commands emit NDJSON events to stderr in default mode

Exit codes: `0` success, `2` usage error, `3` not found, `4` auth error, `5` API error, `7` rate limited, `10` config error.

## Use as MCP Server

This CLI ships a companion MCP server for use with Claude Desktop, Cursor, and other MCP-compatible tools.

### Claude Code

```bash
claude mcp add cloud-functions cloud-functions-pp-mcp -e CLOUD_FUNCTIONS_TOKEN=<your-token>
```

### Claude Desktop

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "cloud-functions": {
      "command": "cloud-functions-pp-mcp",
      "env": {
        "CLOUD_FUNCTIONS_TOKEN": "<your-key>"
      }
    }
  }
}
```

## Cookbook

Common workflows and recipes:

```bash
# List resources as JSON for scripting
cloud-functions-pp-cli name-generate-download-url list --json

# Filter to specific fields
cloud-functions-pp-cli name-generate-download-url list --json --select id,name,status

# Dry run to preview the request
cloud-functions-pp-cli name-generate-download-url list --dry-run

# Sync data locally for offline search
cloud-functions-pp-cli sync

# Search synced data
cloud-functions-pp-cli search "query"

# Export for backup
cloud-functions-pp-cli export --format jsonl > backup.jsonl
```

## Health Check

```bash
cloud-functions-pp-cli doctor
```

<!-- DOCTOR_OUTPUT -->

## Configuration

Config file: `~/.config/cloud-functions-pp-cli/config.toml`

Environment variables:
- `CLOUD_FUNCTIONS_TOKEN`

## Troubleshooting

**Authentication errors (exit code 4)**
- Run `cloud-functions-pp-cli doctor` to check credentials
- Verify the environment variable is set: `echo $CLOUD_FUNCTIONS_TOKEN`

**Not found errors (exit code 3)**
- Check the resource ID is correct
- Run the `list` command to see available items

**Rate limit errors (exit code 7)**
- The CLI auto-retries with exponential backoff
- If persistent, wait a few minutes and try again

---

Generated by [CLI Printing Press](https://github.com/mvanhorn/cli-printing-press)
