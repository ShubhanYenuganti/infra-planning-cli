# Webapps Client CLI



## Install

### Go

```
go install github.com/mvanhorn/printing-press-library/library/other/functions-pp-cli/cmd/functions-pp-cli@latest
```

### Binary

Download from [Releases](https://github.com/mvanhorn/printing-press-library/releases).

## Quick Start

### 1. Install

See [Install](#install) above.

### 2. Set Up Credentials

Get your access token from your API provider's developer portal, then store it:

```bash
functions-pp-cli auth set-token YOUR_TOKEN_HERE
```

Or set it via environment variable:

```bash
export WEBAPPS_CLIENT_TOKEN="your-token-here"
```

### 3. Verify Setup

```bash
functions-pp-cli doctor
```

This checks your configuration and credentials.

### 4. Try Your First Command

```bash
functions-pp-cli analyze-custom-hostname list
```

## Usage

<!-- HELP_OUTPUT -->

## Commands

### analyze-custom-hostname

Manage analyze custom hostname

- **`functions-pp-cli analyze-custom-hostname web-apps`** - Analyze a custom hostname.

### apply-slot-config

Manage apply slot config

- **`functions-pp-cli apply-slot-config web-apps-to-production`** - Applies the configuration settings from the target slot onto the current slot.

### backup

Manage backup

- **`functions-pp-cli backup web-apps`** - Creates a backup of an app.

### backups

Manage backups

- **`functions-pp-cli backups web-apps-delete`** - Deletes a backup of an app by its ID.
- **`functions-pp-cli backups web-apps-get-status`** - Gets a backup of an app by its ID.
- **`functions-pp-cli backups web-apps-list`** - Gets existing backups of an app.

### config

Manage config

- **`functions-pp-cli config web-apps-create-or-update-configuration`** - Updates the configuration of an app.
- **`functions-pp-cli config web-apps-delete-backup-configuration`** - Deletes the backup configuration of an app.
- **`functions-pp-cli config web-apps-get-auth-settings`** - Gets the Authentication/Authorization settings of an app.
- **`functions-pp-cli config web-apps-get-backup-configuration`** - Gets the backup configuration of an app.
- **`functions-pp-cli config web-apps-get-configuration`** - Gets the configuration of an app, such as platform version and bitness, default documents, virtual applications, Always On, etc.
- **`functions-pp-cli config web-apps-get-configuration-snapshot`** - Gets a snapshot of the configuration of an app at a previous point in time.
- **`functions-pp-cli config web-apps-get-diagnostic-logs-configuration`** - Gets the logging configuration of an app.
- **`functions-pp-cli config web-apps-list-application-settings`** - Gets the application settings of an app.
- **`functions-pp-cli config web-apps-list-azure-storage-accounts`** - Gets the Azure storage account configurations of an app.
- **`functions-pp-cli config web-apps-list-configuration-snapshot-info`** - Gets a list of web app configuration snapshots identifiers. Each element of the list contains a timestamp and the ID of the snapshot.
- **`functions-pp-cli config web-apps-list-configurations`** - List the configurations of an app
- **`functions-pp-cli config web-apps-list-connection-strings`** - Gets the connection strings of an app.
- **`functions-pp-cli config web-apps-list-metadata`** - Gets the metadata of an app.
- **`functions-pp-cli config web-apps-list-publishing-credentials`** - Gets the Git/FTP publishing credentials of an app.
- **`functions-pp-cli config web-apps-list-site-push-settings`** - Gets the Push settings associated with web app.
- **`functions-pp-cli config web-apps-list-slot-configuration-names`** - Gets the names of app settings and connection strings that stick to the slot (not swapped).
- **`functions-pp-cli config web-apps-recover-site-configuration-snapshot`** - Reverts the configuration of an app to a previous snapshot.
- **`functions-pp-cli config web-apps-update-application-settings`** - Replaces the application settings of an app.
- **`functions-pp-cli config web-apps-update-auth-settings`** - Updates the Authentication / Authorization settings associated with web app.
- **`functions-pp-cli config web-apps-update-azure-storage-accounts`** - Updates the Azure storage account configurations of an app.
- **`functions-pp-cli config web-apps-update-backup-configuration`** - Updates the backup configuration of an app.
- **`functions-pp-cli config web-apps-update-configuration`** - Updates the configuration of an app.
- **`functions-pp-cli config web-apps-update-connection-strings`** - Replaces the connection strings of an app.
- **`functions-pp-cli config web-apps-update-diagnostic-logs`** - Updates the logging configuration of an app.
- **`functions-pp-cli config web-apps-update-metadata`** - Replaces the metadata of an app.
- **`functions-pp-cli config web-apps-update-site-push-settings`** - Updates the Push settings associated with web app.
- **`functions-pp-cli config web-apps-update-slot-configuration-names`** - Updates the names of application settings and connection string that remain with the slot during swap operation.

### containerlogs

Manage containerlogs

- **`functions-pp-cli containerlogs web-apps-get-container-logs-zip`** - Gets the ZIP archived docker log files for the given site
- **`functions-pp-cli containerlogs web-apps-get-web-site-container-logs`** - Gets the last lines of docker logs for the given site

### continuouswebjobs

Manage continuouswebjobs

- **`functions-pp-cli continuouswebjobs web-apps-delete-continuous-web-job`** - Delete a continuous web job by its ID for an app, or a deployment slot.
- **`functions-pp-cli continuouswebjobs web-apps-get-continuous-web-job`** - Gets a continuous web job by its ID for an app, or a deployment slot.
- **`functions-pp-cli continuouswebjobs web-apps-list-continuous-web-jobs`** - List continuous web jobs for an app, or a deployment slot.

### deployments

Manage deployments

- **`functions-pp-cli deployments web-apps-create`** - Create a deployment for an app, or a deployment slot.
- **`functions-pp-cli deployments web-apps-delete`** - Delete a deployment by its ID for an app, or a deployment slot.
- **`functions-pp-cli deployments web-apps-get`** - Get a deployment by its ID for an app, or a deployment slot.
- **`functions-pp-cli deployments web-apps-list`** - List deployments for an app, or a deployment slot.

### discoverbackup

Manage discoverbackup

- **`functions-pp-cli discoverbackup web-apps-discover-backup`** - Discovers an existing app backup that can be restored from a blob in Azure storage. Use this to get information about the databases stored in a backup.

### domain-ownership-identifiers

Manage domain ownership identifiers

- **`functions-pp-cli domain-ownership-identifiers web-apps-create-or-update`** - Creates a domain ownership identifier for web app, or updates an existing ownership identifier.
- **`functions-pp-cli domain-ownership-identifiers web-apps-delete`** - Deletes a domain ownership identifier for a web app.
- **`functions-pp-cli domain-ownership-identifiers web-apps-get`** - Get domain ownership identifier for web app.
- **`functions-pp-cli domain-ownership-identifiers web-apps-list`** - Lists ownership identifiers for domain associated with web app.
- **`functions-pp-cli domain-ownership-identifiers web-apps-update`** - Creates a domain ownership identifier for web app, or updates an existing ownership identifier.

### extensions

Manage extensions

- **`functions-pp-cli extensions web-apps-create-msdeploy-operation`** - Invoke the MSDeploy web app extension.
- **`functions-pp-cli extensions web-apps-get-msdeploy-log`** - Get the MSDeploy Log for the last MSDeploy operation.
- **`functions-pp-cli extensions web-apps-get-msdeploy-status`** - Get the status of the last MSDeploy operation.

### functions

Manage functions

- **`functions-pp-cli functions web-apps-create`** - Create function for web site, or a deployment slot.
- **`functions-pp-cli functions web-apps-delete`** - Delete a function for web site, or a deployment slot.
- **`functions-pp-cli functions web-apps-get`** - Get function information by its ID for web site, or a deployment slot.
- **`functions-pp-cli functions web-apps-get-admin-token`** - Fetch a short lived token that can be exchanged for a master key.
- **`functions-pp-cli functions web-apps-list`** - List the functions for a web site, or a deployment slot.

### host-name-bindings

Manage host name bindings

- **`functions-pp-cli host-name-bindings web-apps-create-or-update`** - Creates a hostname binding for an app.
- **`functions-pp-cli host-name-bindings web-apps-delete`** - Deletes a hostname binding for an app.
- **`functions-pp-cli host-name-bindings web-apps-get`** - Get the named hostname binding for an app (or deployment slot, if specified).
- **`functions-pp-cli host-name-bindings web-apps-list`** - Get hostname bindings for an app or a deployment slot.

### hybrid-connection-namespaces

Manage hybrid connection namespaces


### hybrid-connection-relays

Manage hybrid connection relays

- **`functions-pp-cli hybrid-connection-relays web-apps-list-hybrid-connections`** - Retrieves all Service Bus Hybrid Connections used by this Web App.

### hybridconnection

Manage hybridconnection

- **`functions-pp-cli hybridconnection web-apps-create-or-update-relay-service-connection`** - Creates a new hybrid connection configuration (PUT), or updates an existing one (PATCH).
- **`functions-pp-cli hybridconnection web-apps-delete-relay-service-connection`** - Deletes a relay service connection by its name.
- **`functions-pp-cli hybridconnection web-apps-get-relay-service-connection`** - Gets a hybrid connection configuration by its name.
- **`functions-pp-cli hybridconnection web-apps-list-relay-service-connections`** - Gets hybrid connections configured for an app (or deployment slot, if specified).
- **`functions-pp-cli hybridconnection web-apps-update-relay-service-connection`** - Creates a new hybrid connection configuration (PUT), or updates an existing one (PATCH).

### instances

Manage instances

- **`functions-pp-cli instances web-apps-list-identifiers`** - Gets all scale-out instances of an app.

### iscloneable

Manage iscloneable

- **`functions-pp-cli iscloneable web-apps-is-cloneable`** - Shows whether an app can be cloned to another resource group or subscription.

### listsyncfunctiontriggerstatus

Manage listsyncfunctiontriggerstatus

- **`functions-pp-cli listsyncfunctiontriggerstatus web-apps-list-sync-function-triggers`** - This is to allow calling via powershell and ARM template.

### metricdefinitions

Manage metricdefinitions

- **`functions-pp-cli metricdefinitions web-apps-list-metric-definitions`** - Gets all metric definitions of an app (or deployment slot, if specified).

### metrics

Manage metrics

- **`functions-pp-cli metrics web-apps-list`** - Gets performance metrics of an app (or deployment slot, if specified).

### migrate

Manage migrate

- **`functions-pp-cli migrate web-apps-storage`** - Restores a web app.

### migratemysql

Manage migratemysql

- **`functions-pp-cli migratemysql web-apps-get-migrate-my-sql-status`** - Returns the status of MySql in app migration, if one is active, and whether or not MySql in app is enabled
- **`functions-pp-cli migratemysql web-apps-migrate-my-sql`** - Migrates a local (in-app) MySql database to a remote MySql database.

### network-config

Manage network config

- **`functions-pp-cli network-config web-apps-create-or-update-swift-virtual-network-connection`** - Integrates this Web App with a Virtual Network. This requires that 1) "swiftSupported" is true when doing a GET against this resource, and 2) that the target Subnet has already been delegated, and is not
in use by another App Service Plan other than the one this App is in.
- **`functions-pp-cli network-config web-apps-delete-swift-virtual-network`** - Deletes a Swift Virtual Network connection from an app (or deployment slot).
- **`functions-pp-cli network-config web-apps-get-swift-virtual-network-connection`** - Gets a Swift Virtual Network connection.
- **`functions-pp-cli network-config web-apps-update-swift-virtual-network-connection`** - Integrates this Web App with a Virtual Network. This requires that 1) "swiftSupported" is true when doing a GET against this resource, and 2) that the target Subnet has already been delegated, and is not
in use by another App Service Plan other than the one this App is in.

### network-features

Manage network features

- **`functions-pp-cli network-features web-apps-list`** - Gets all network features used by the app (or deployment slot, if specified).

### network-trace

Manage network trace

- **`functions-pp-cli network-trace web-apps-get`** - Gets a named operation for a network trace capturing (or deployment slot, if specified).
- **`functions-pp-cli network-trace web-apps-get-operation`** - Gets a named operation for a network trace capturing (or deployment slot, if specified).
- **`functions-pp-cli network-trace web-apps-start-web-site`** - Start capturing network packets for the site (To be deprecated).
- **`functions-pp-cli network-trace web-apps-start-web-site-operation`** - Start capturing network packets for the site.
- **`functions-pp-cli network-trace web-apps-stop-web-site`** - Stop ongoing capturing network packets for the site.

### network-traces

Manage network traces

- **`functions-pp-cli network-traces web-apps-get-operation-v2`** - Gets a named operation for a network trace capturing (or deployment slot, if specified).
- **`functions-pp-cli network-traces web-apps-get-v2`** - Gets a named operation for a network trace capturing (or deployment slot, if specified).

### newpassword

Manage newpassword

- **`functions-pp-cli newpassword web-apps-generate-new-site-publishing-password`** - Generates a new publishing password for an app (or deployment slot, if specified).

### perfcounters

Manage perfcounters

- **`functions-pp-cli perfcounters web-apps-list-perf-mon-counters`** - Gets perfmon counters for web app.

### phplogging

Manage phplogging

- **`functions-pp-cli phplogging web-apps-get-site-php-error-log-flag`** - Gets web app's event logs.

### premieraddons

Manage premieraddons

- **`functions-pp-cli premieraddons web-apps-add-premier-add-on`** - Updates a named add-on of an app.
- **`functions-pp-cli premieraddons web-apps-delete-premier-add-on`** - Delete a premier add-on from an app.
- **`functions-pp-cli premieraddons web-apps-get-premier-add-on`** - Gets a named add-on of an app.
- **`functions-pp-cli premieraddons web-apps-list-premier-add-ons`** - Gets the premier add-ons of an app.
- **`functions-pp-cli premieraddons web-apps-update-premier-add-on`** - Updates a named add-on of an app.

### private-access

Manage private access

- **`functions-pp-cli private-access web-apps-get`** - Gets data around private site access enablement and authorized Virtual Networks that can access the site.
- **`functions-pp-cli private-access web-apps-put-vnet`** - Sets data around private site access enablement and authorized Virtual Networks that can access the site.

### processes

Manage processes

- **`functions-pp-cli processes web-apps-delete-process`** - Terminate a process by its ID for a web site, or a deployment slot, or specific scaled-out instance in a web site.
- **`functions-pp-cli processes web-apps-get-process`** - Get process information by its ID for a specific scaled-out instance in a web site.
- **`functions-pp-cli processes web-apps-list`** - Get list of processes for a web site, or a deployment slot, or for a specific scaled-out instance in a web site.

### public-certificates

Manage public certificates

- **`functions-pp-cli public-certificates web-apps-create-or-update`** - Creates a hostname binding for an app.
- **`functions-pp-cli public-certificates web-apps-delete`** - Deletes a hostname binding for an app.
- **`functions-pp-cli public-certificates web-apps-get`** - Get the named public certificate for an app (or deployment slot, if specified).
- **`functions-pp-cli public-certificates web-apps-list`** - Get public certificates for an app or a deployment slot.

### publishxml

Manage publishxml

- **`functions-pp-cli publishxml web-apps-list-publishing-profile-xml-with-secrets`** - Gets the publishing profile for an app (or deployment slot, if specified).

### reset-slot-config

Manage reset slot config

- **`functions-pp-cli reset-slot-config web-apps-reset-production-slot-config`** - Resets the configuration settings of the current slot if they were previously modified by calling the API with POST.

### restart

Manage restart

- **`functions-pp-cli restart web-apps`** - Restarts an app (or deployment slot, if specified).

### restore-from-backup-blob

Manage restore from backup blob

- **`functions-pp-cli restore-from-backup-blob web-apps`** - Restores an app from a backup blob in Azure Storage.

### restore-from-deleted-app

Manage restore from deleted app

- **`functions-pp-cli restore-from-deleted-app web-apps`** - Restores a deleted web app to this web app.

### restore-snapshot

Manage restore snapshot

- **`functions-pp-cli restore-snapshot web-apps`** - Restores a web app from a snapshot.

### siteextensions

Manage siteextensions

- **`functions-pp-cli siteextensions web-apps-delete-site-extension`** - Remove a site extension from a web site, or a deployment slot.
- **`functions-pp-cli siteextensions web-apps-get-site-extension`** - Get site extension information by its ID for a web site, or a deployment slot.
- **`functions-pp-cli siteextensions web-apps-install-site-extension`** - Install site extension on a web site, or a deployment slot.
- **`functions-pp-cli siteextensions web-apps-list-site-extensions`** - Get list of siteextensions for a web site, or a deployment slot.

### slots

Manage slots

- **`functions-pp-cli slots web-apps-create-or-update`** - Creates a new web, mobile, or API app in an existing resource group, or updates an existing app.
- **`functions-pp-cli slots web-apps-delete`** - Deletes a web, mobile, or API app, or one of the deployment slots.
- **`functions-pp-cli slots web-apps-get`** - Gets the details of a web, mobile, or API app.
- **`functions-pp-cli slots web-apps-list`** - Gets an app's deployment slots.
- **`functions-pp-cli slots web-apps-update`** - Creates a new web, mobile, or API app in an existing resource group, or updates an existing app.

### slotsdiffs

Manage slotsdiffs

- **`functions-pp-cli slotsdiffs web-apps-list-slot-differences-from-production`** - Get the difference in configuration settings between two web app slots.

### slotsswap

Manage slotsswap

- **`functions-pp-cli slotsswap web-apps-swap-slot-with-production`** - Swaps two deployment slots of an app.

### snapshots

Manage snapshots

- **`functions-pp-cli snapshots web-apps-list`** - Returns all Snapshots to the user.

### snapshotsdr

Manage snapshotsdr

- **`functions-pp-cli snapshotsdr web-apps-list-snapshots-from-drsecondary`** - Returns all Snapshots to the user from DRSecondary endpoint.

### sourcecontrols

Manage sourcecontrols

- **`functions-pp-cli sourcecontrols web-apps-create-or-update-source-control`** - Updates the source control configuration of an app.
- **`functions-pp-cli sourcecontrols web-apps-delete-source-control`** - Deletes the source control configuration of an app.
- **`functions-pp-cli sourcecontrols web-apps-get-source-control`** - Gets the source control configuration of an app.
- **`functions-pp-cli sourcecontrols web-apps-update-source-control`** - Updates the source control configuration of an app.

### start

Manage start

- **`functions-pp-cli start web-apps`** - Starts an app (or deployment slot, if specified).

### start-network-trace

Manage start network trace

- **`functions-pp-cli start-network-trace web-apps`** - Start capturing network packets for the site.

### stop

Manage stop

- **`functions-pp-cli stop web-apps`** - Stops an app (or deployment slot, if specified).

### stop-network-trace

Manage stop network trace

- **`functions-pp-cli stop-network-trace web-apps`** - Stop ongoing capturing network packets for the site.

### subscriptions

Manage subscriptions


### sync

Manage sync

- **`functions-pp-cli sync web-apps-repository`** - Sync web app repository.

### syncfunctiontriggers

Manage syncfunctiontriggers

- **`functions-pp-cli syncfunctiontriggers web-apps-sync-function-triggers`** - Syncs function trigger metadata to the scale controller

### triggeredwebjobs

Manage triggeredwebjobs

- **`functions-pp-cli triggeredwebjobs web-apps-delete-triggered-web-job`** - Delete a triggered web job by its ID for an app, or a deployment slot.
- **`functions-pp-cli triggeredwebjobs web-apps-get-triggered-web-job`** - Gets a triggered web job by its ID for an app, or a deployment slot.
- **`functions-pp-cli triggeredwebjobs web-apps-list-triggered-web-jobs`** - List triggered web jobs for an app, or a deployment slot.

### usages

Manage usages

- **`functions-pp-cli usages web-apps-list`** - Gets the quota usage information of an app (or deployment slot, if specified).

### virtual-network-connections

Manage virtual network connections

- **`functions-pp-cli virtual-network-connections web-apps-create-or-update-vnet-connection`** - Adds a Virtual Network connection to an app or slot (PUT) or updates the connection properties (PATCH).
- **`functions-pp-cli virtual-network-connections web-apps-delete-vnet-connection`** - Deletes a connection from an app (or deployment slot to a named virtual network.
- **`functions-pp-cli virtual-network-connections web-apps-get-vnet-connection`** - Gets a virtual network the app (or deployment slot) is connected to by name.
- **`functions-pp-cli virtual-network-connections web-apps-list-vnet-connections`** - Gets the virtual networks the app (or deployment slot) is connected to.
- **`functions-pp-cli virtual-network-connections web-apps-update-vnet-connection`** - Adds a Virtual Network connection to an app or slot (PUT) or updates the connection properties (PATCH).

### webjobs

Manage webjobs

- **`functions-pp-cli webjobs web-apps-get-web-job`** - Get webjob information for an app, or a deployment slot.
- **`functions-pp-cli webjobs web-apps-list-web-jobs`** - List webjobs for an app, or a deployment slot.


## Output Formats

```bash
# Human-readable table (default in terminal, JSON when piped)
functions-pp-cli analyze-custom-hostname list

# JSON for scripting and agents
functions-pp-cli analyze-custom-hostname list --json

# Filter to specific fields
functions-pp-cli analyze-custom-hostname list --json --select id,name,status

# Dry run — show the request without sending
functions-pp-cli analyze-custom-hostname list --dry-run

# Agent mode — JSON + compact + no prompts in one flag
functions-pp-cli analyze-custom-hostname list --agent
```

## Agent Usage

This CLI is designed for AI agent consumption:

- **Non-interactive** - never prompts, every input is a flag
- **Pipeable** - `--json` output to stdout, errors to stderr
- **Filterable** - `--select id,name` returns only fields you need
- **Previewable** - `--dry-run` shows the request without sending
- **Retryable** - creates return "already exists" on retry, deletes return "already deleted"
- **Confirmable** - `--yes` for explicit confirmation of destructive actions
- **Piped input** - `echo '{"key":"value"}' | functions-pp-cli <resource> create --stdin`
- **Cacheable** - GET responses cached for 5 minutes, bypass with `--no-cache`
- **Agent-safe by default** - no colors or formatting unless `--human-friendly` is set
- **Progress events** - paginated commands emit NDJSON events to stderr in default mode

Exit codes: `0` success, `2` usage error, `3` not found, `4` auth error, `5` API error, `7` rate limited, `10` config error.

## Use as MCP Server

This CLI ships a companion MCP server for use with Claude Desktop, Cursor, and other MCP-compatible tools.

### Claude Code

```bash
claude mcp add functions functions-pp-mcp -e WEBAPPS_CLIENT_TOKEN=<your-token>
```

### Claude Desktop

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "functions": {
      "command": "functions-pp-mcp",
      "env": {
        "WEBAPPS_CLIENT_TOKEN": "<your-key>"
      }
    }
  }
}
```

## Cookbook

Common workflows and recipes:

```bash
# List resources as JSON for scripting
functions-pp-cli analyze-custom-hostname list --json

# Filter to specific fields
functions-pp-cli analyze-custom-hostname list --json --select id,name,status

# Dry run to preview the request
functions-pp-cli analyze-custom-hostname list --dry-run

# Sync data locally for offline search
functions-pp-cli sync

# Search synced data
functions-pp-cli search "query"

# Export for backup
functions-pp-cli export --format jsonl > backup.jsonl
```

## Health Check

```bash
functions-pp-cli doctor
```

<!-- DOCTOR_OUTPUT -->

## Configuration

Config file: `~/.config/functions-pp-cli/config.toml`

Environment variables:
- `WEBAPPS_CLIENT_TOKEN`

## Troubleshooting

**Authentication errors (exit code 4)**
- Run `functions-pp-cli doctor` to check credentials
- Verify the environment variable is set: `echo $WEBAPPS_CLIENT_TOKEN`

**Not found errors (exit code 3)**
- Check the resource ID is correct
- Run the `list` command to see available items

**Rate limit errors (exit code 7)**
- The CLI auto-retries with exponential backoff
- If persistent, wait a few minutes and try again

---

Generated by [CLI Printing Press](https://github.com/mvanhorn/cli-printing-press)
