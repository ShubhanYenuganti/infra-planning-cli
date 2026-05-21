# Aws App Runner CLI

<fullname>App Runner</fullname> <p>App Runner is an application service that provides a fast, simple, and cost-effective way to go directly from an existing container image or source code to a running service in the Amazon Web Services Cloud in seconds. You don't need to learn new technologies, decide which compute service to use, or understand how to provision and configure Amazon Web Services resources.</p> <p>App Runner connects directly to your container registry or source code repository. It provides an automatic delivery pipeline with fully managed operations, high performance, scalability, and security.</p> <p>For more information about App Runner, see the <a href="https://docs.aws.amazon.com/apprunner/latest/dg/">App Runner Developer Guide</a>. For release information, see the <a href="https://docs.aws.amazon.com/apprunner/latest/relnotes/">App Runner Release Notes</a>.</p> <p> To install the Software Development Kits (SDKs), Integrated Development Environment (IDE) Toolkits, and command line tools that you can use to access the API, see <a href="http://aws.amazon.com/tools/">Tools for Amazon Web Services</a>.</p> <p> <b>Endpoints</b> </p> <p>For a list of Region-specific endpoints that App Runner supports, see <a href="https://docs.aws.amazon.com/general/latest/gr/apprunner.html">App Runner endpoints and quotas</a> in the <i>Amazon Web Services General Reference</i>.</p>

Learn more at [Aws App Runner](https://github.com/mermade/aws2openapi).

## Install

### Go

```
go install github.com/mvanhorn/printing-press-library/library/other/aws-app-runner-pp-cli/cmd/aws-app-runner-pp-cli@latest
```

### Binary

Download from [Releases](https://github.com/mvanhorn/printing-press-library/releases).

## Quick Start

### 1. Install

See [Install](#install) above.

### 2. Set Up Credentials

Get your API key from your API provider's developer portal. The key typically looks like a long alphanumeric string.

```bash
export AWS_APP_RUNNER_HMAC="<paste-your-key>"
```

You can also persist this in your config file at `~/.config/aws-app-runner-pp-cli/config.toml`.

### 3. Verify Setup

```bash
aws-app-runner-pp-cli doctor
```

This checks your configuration and credentials.

### 4. Try Your First Command

```bash
aws-app-runner-pp-cli x-amz-target-app-runner-associate-custom-domain list
```

## Usage

<!-- HELP_OUTPUT -->

## Commands

### x-amz-target-app-runner-associate-custom-domain

Manage x amz target app runner associate custom domain

- **`aws-app-runner-pp-cli x-amz-target-app-runner-associate-custom-domain associate-custom-domain`** - <p>Associate your own domain name with the App Runner subdomain URL of your App Runner service.</p> <p>After you call <code>AssociateCustomDomain</code> and receive a successful response, use the information in the <a>CustomDomain</a> record that's returned to add CNAME records to your Domain Name System (DNS). For each mapped domain name, add a mapping to the target App Runner subdomain and one or more certificate validation records. App Runner then performs DNS validation to verify that you own or control the domain name that you associated. App Runner tracks domain validity in a certificate stored in <a href="https://docs.aws.amazon.com/acm/latest/userguide">AWS Certificate Manager (ACM)</a>.</p>

### x-amz-target-app-runner-create-auto-scaling-configuration

Manage x amz target app runner create auto scaling configuration

- **`aws-app-runner-pp-cli x-amz-target-app-runner-create-auto-scaling-configuration create-auto-scaling-configuration`** - <p>Create an App Runner automatic scaling configuration resource. App Runner requires this resource when you create or update App Runner services and you require non-default auto scaling settings. You can share an auto scaling configuration across multiple services.</p> <p>Create multiple revisions of a configuration by calling this action multiple times using the same <code>AutoScalingConfigurationName</code>. The call returns incremental <code>AutoScalingConfigurationRevision</code> values. When you create a service and configure an auto scaling configuration resource, the service uses the latest active revision of the auto scaling configuration by default. You can optionally configure the service to use a specific revision.</p> <p>Configure a higher <code>MinSize</code> to increase the spread of your App Runner service over more Availability Zones in the Amazon Web Services Region. The tradeoff is a higher minimal cost.</p> <p>Configure a lower <code>MaxSize</code> to control your cost. The tradeoff is lower responsiveness during peak demand.</p>

### x-amz-target-app-runner-create-connection

Manage x amz target app runner create connection

- **`aws-app-runner-pp-cli x-amz-target-app-runner-create-connection create-connection`** - <p>Create an App Runner connection resource. App Runner requires a connection resource when you create App Runner services that access private repositories from certain third-party providers. You can share a connection across multiple services.</p> <p>A connection resource is needed to access GitHub repositories. GitHub requires a user interface approval process through the App Runner console before you can use the connection.</p>

### x-amz-target-app-runner-create-observability-configuration

Manage x amz target app runner create observability configuration

- **`aws-app-runner-pp-cli x-amz-target-app-runner-create-observability-configuration create-observability-configuration`** - <p>Create an App Runner observability configuration resource. App Runner requires this resource when you create or update App Runner services and you want to enable non-default observability features. You can share an observability configuration across multiple services.</p> <p>Create multiple revisions of a configuration by calling this action multiple times using the same <code>ObservabilityConfigurationName</code>. The call returns incremental <code>ObservabilityConfigurationRevision</code> values. When you create a service and configure an observability configuration resource, the service uses the latest active revision of the observability configuration by default. You can optionally configure the service to use a specific revision.</p> <p>The observability configuration resource is designed to configure multiple features (currently one feature, tracing). This action takes optional parameters that describe the configuration of these features (currently one parameter, <code>TraceConfiguration</code>). If you don't specify a feature parameter, App Runner doesn't enable the feature.</p>

### x-amz-target-app-runner-create-service

Manage x amz target app runner create service

- **`aws-app-runner-pp-cli x-amz-target-app-runner-create-service create-service`** - <p>Create an App Runner service. After the service is created, the action also automatically starts a deployment.</p> <p>This is an asynchronous operation. On a successful call, you can use the returned <code>OperationId</code> and the <a href="https://docs.aws.amazon.com/apprunner/latest/api/API_ListOperations.html">ListOperations</a> call to track the operation's progress.</p>

### x-amz-target-app-runner-create-vpc-connector

Manage x amz target app runner create vpc connector

- **`aws-app-runner-pp-cli x-amz-target-app-runner-create-vpc-connector create-vpc-connector`** - Create an App Runner VPC connector resource. App Runner requires this resource when you want to associate your App Runner service to a custom Amazon Virtual Private Cloud (Amazon VPC).

### x-amz-target-app-runner-create-vpc-ingress-connection

Manage x amz target app runner create vpc ingress connection

- **`aws-app-runner-pp-cli x-amz-target-app-runner-create-vpc-ingress-connection create-vpc-ingress-connection`** - Create an App Runner VPC Ingress Connection resource. App Runner requires this resource when you want to associate your App Runner service with an Amazon VPC endpoint.

### x-amz-target-app-runner-delete-auto-scaling-configuration

Manage x amz target app runner delete auto scaling configuration

- **`aws-app-runner-pp-cli x-amz-target-app-runner-delete-auto-scaling-configuration delete-auto-scaling-configuration`** - Delete an App Runner automatic scaling configuration resource. You can delete a specific revision or the latest active revision. You can't delete a configuration that's used by one or more App Runner services.

### x-amz-target-app-runner-delete-connection

Manage x amz target app runner delete connection

- **`aws-app-runner-pp-cli x-amz-target-app-runner-delete-connection delete-connection`** - Delete an App Runner connection. You must first ensure that there are no running App Runner services that use this connection. If there are any, the <code>DeleteConnection</code> action fails.

### x-amz-target-app-runner-delete-observability-configuration

Manage x amz target app runner delete observability configuration

- **`aws-app-runner-pp-cli x-amz-target-app-runner-delete-observability-configuration delete-observability-configuration`** - Delete an App Runner observability configuration resource. You can delete a specific revision or the latest active revision. You can't delete a configuration that's used by one or more App Runner services.

### x-amz-target-app-runner-delete-service

Manage x amz target app runner delete service

- **`aws-app-runner-pp-cli x-amz-target-app-runner-delete-service delete-service`** - <p>Delete an App Runner service.</p> <p>This is an asynchronous operation. On a successful call, you can use the returned <code>OperationId</code> and the <a>ListOperations</a> call to track the operation's progress.</p> <note> <p>Make sure that you don't have any active VPCIngressConnections associated with the service you want to delete. </p> </note>

### x-amz-target-app-runner-delete-vpc-connector

Manage x amz target app runner delete vpc connector

- **`aws-app-runner-pp-cli x-amz-target-app-runner-delete-vpc-connector delete-vpc-connector`** - Delete an App Runner VPC connector resource. You can't delete a connector that's used by one or more App Runner services.

### x-amz-target-app-runner-delete-vpc-ingress-connection

Manage x amz target app runner delete vpc ingress connection

- **`aws-app-runner-pp-cli x-amz-target-app-runner-delete-vpc-ingress-connection delete-vpc-ingress-connection`** - <p>Delete an App Runner VPC Ingress Connection resource that's associated with an App Runner service. The VPC Ingress Connection must be in one of the following states to be deleted: </p> <ul> <li> <p> <code>AVAILABLE</code> </p> </li> <li> <p> <code>FAILED_CREATION</code> </p> </li> <li> <p> <code>FAILED_UPDATE</code> </p> </li> <li> <p> <code>FAILED_DELETION</code> </p> </li> </ul>

### x-amz-target-app-runner-describe-auto-scaling-configuration

Manage x amz target app runner describe auto scaling configuration

- **`aws-app-runner-pp-cli x-amz-target-app-runner-describe-auto-scaling-configuration describe-auto-scaling-configuration`** - Return a full description of an App Runner automatic scaling configuration resource.

### x-amz-target-app-runner-describe-custom-domains

Manage x amz target app runner describe custom domains

- **`aws-app-runner-pp-cli x-amz-target-app-runner-describe-custom-domains describe-custom-domains`** - Return a description of custom domain names that are associated with an App Runner service.

### x-amz-target-app-runner-describe-observability-configuration

Manage x amz target app runner describe observability configuration

- **`aws-app-runner-pp-cli x-amz-target-app-runner-describe-observability-configuration describe-observability-configuration`** - Return a full description of an App Runner observability configuration resource.

### x-amz-target-app-runner-describe-service

Manage x amz target app runner describe service

- **`aws-app-runner-pp-cli x-amz-target-app-runner-describe-service describe-service`** - Return a full description of an App Runner service.

### x-amz-target-app-runner-describe-vpc-connector

Manage x amz target app runner describe vpc connector

- **`aws-app-runner-pp-cli x-amz-target-app-runner-describe-vpc-connector describe-vpc-connector`** - Return a description of an App Runner VPC connector resource.

### x-amz-target-app-runner-describe-vpc-ingress-connection

Manage x amz target app runner describe vpc ingress connection

- **`aws-app-runner-pp-cli x-amz-target-app-runner-describe-vpc-ingress-connection describe-vpc-ingress-connection`** - Return a full description of an App Runner VPC Ingress Connection resource.

### x-amz-target-app-runner-disassociate-custom-domain

Manage x amz target app runner disassociate custom domain

- **`aws-app-runner-pp-cli x-amz-target-app-runner-disassociate-custom-domain disassociate-custom-domain`** - <p>Disassociate a custom domain name from an App Runner service.</p> <p>Certificates tracking domain validity are associated with a custom domain and are stored in <a href="https://docs.aws.amazon.com/acm/latest/userguide">AWS Certificate Manager (ACM)</a>. These certificates aren't deleted as part of this action. App Runner delays certificate deletion for 30 days after a domain is disassociated from your service.</p>

### x-amz-target-app-runner-list-auto-scaling-configurations

Manage x amz target app runner list auto scaling configurations

- **`aws-app-runner-pp-cli x-amz-target-app-runner-list-auto-scaling-configurations list-auto-scaling-configurations`** - <p>Returns a list of active App Runner automatic scaling configurations in your Amazon Web Services account. You can query the revisions for a specific configuration name or the revisions for all active configurations in your account. You can optionally query only the latest revision of each requested name.</p> <p>To retrieve a full description of a particular configuration revision, call and provide one of the ARNs returned by <code>ListAutoScalingConfigurations</code>.</p>

### x-amz-target-app-runner-list-connections

Manage x amz target app runner list connections

- **`aws-app-runner-pp-cli x-amz-target-app-runner-list-connections list-connections`** - Returns a list of App Runner connections that are associated with your Amazon Web Services account.

### x-amz-target-app-runner-list-observability-configurations

Manage x amz target app runner list observability configurations

- **`aws-app-runner-pp-cli x-amz-target-app-runner-list-observability-configurations list-observability-configurations`** - <p>Returns a list of active App Runner observability configurations in your Amazon Web Services account. You can query the revisions for a specific configuration name or the revisions for all active configurations in your account. You can optionally query only the latest revision of each requested name.</p> <p>To retrieve a full description of a particular configuration revision, call and provide one of the ARNs returned by <code>ListObservabilityConfigurations</code>.</p>

### x-amz-target-app-runner-list-operations

Manage x amz target app runner list operations

- **`aws-app-runner-pp-cli x-amz-target-app-runner-list-operations list-operations`** - <p>Return a list of operations that occurred on an App Runner service.</p> <p>The resulting list of <a>OperationSummary</a> objects is sorted in reverse chronological order. The first object on the list represents the last started operation.</p>

### x-amz-target-app-runner-list-services

Manage x amz target app runner list services

- **`aws-app-runner-pp-cli x-amz-target-app-runner-list-services list-services`** - Returns a list of running App Runner services in your Amazon Web Services account.

### x-amz-target-app-runner-list-tags-for-resource

Manage x amz target app runner list tags for resource

- **`aws-app-runner-pp-cli x-amz-target-app-runner-list-tags-for-resource list-tags-for-resource`** - List tags that are associated with for an App Runner resource. The response contains a list of tag key-value pairs.

### x-amz-target-app-runner-list-vpc-connectors

Manage x amz target app runner list vpc connectors

- **`aws-app-runner-pp-cli x-amz-target-app-runner-list-vpc-connectors list-vpc-connectors`** - Returns a list of App Runner VPC connectors in your Amazon Web Services account.

### x-amz-target-app-runner-list-vpc-ingress-connections

Manage x amz target app runner list vpc ingress connections

- **`aws-app-runner-pp-cli x-amz-target-app-runner-list-vpc-ingress-connections list-vpc-ingress-connections`** - Return a list of App Runner VPC Ingress Connections in your Amazon Web Services account.

### x-amz-target-app-runner-pause-service

Manage x amz target app runner pause service

- **`aws-app-runner-pp-cli x-amz-target-app-runner-pause-service pause-service`** - <p>Pause an active App Runner service. App Runner reduces compute capacity for the service to zero and loses state (for example, ephemeral storage is removed).</p> <p>This is an asynchronous operation. On a successful call, you can use the returned <code>OperationId</code> and the <a>ListOperations</a> call to track the operation's progress.</p>

### x-amz-target-app-runner-resume-service

Manage x amz target app runner resume service

- **`aws-app-runner-pp-cli x-amz-target-app-runner-resume-service resume-service`** - <p>Resume an active App Runner service. App Runner provisions compute capacity for the service.</p> <p>This is an asynchronous operation. On a successful call, you can use the returned <code>OperationId</code> and the <a>ListOperations</a> call to track the operation's progress.</p>

### x-amz-target-app-runner-start-deployment

Manage x amz target app runner start deployment

- **`aws-app-runner-pp-cli x-amz-target-app-runner-start-deployment start-deployment`** - <p>Initiate a manual deployment of the latest commit in a source code repository or the latest image in a source image repository to an App Runner service.</p> <p>For a source code repository, App Runner retrieves the commit and builds a Docker image. For a source image repository, App Runner retrieves the latest Docker image. In both cases, App Runner then deploys the new image to your service and starts a new container instance.</p> <p>This is an asynchronous operation. On a successful call, you can use the returned <code>OperationId</code> and the <a>ListOperations</a> call to track the operation's progress.</p>

### x-amz-target-app-runner-tag-resource

Manage x amz target app runner tag resource

- **`aws-app-runner-pp-cli x-amz-target-app-runner-tag-resource tag-resource`** - Add tags to, or update the tag values of, an App Runner resource. A tag is a key-value pair.

### x-amz-target-app-runner-untag-resource

Manage x amz target app runner untag resource

- **`aws-app-runner-pp-cli x-amz-target-app-runner-untag-resource untag-resource`** - Remove tags from an App Runner resource.

### x-amz-target-app-runner-update-service

Manage x amz target app runner update service

- **`aws-app-runner-pp-cli x-amz-target-app-runner-update-service update-service`** - <p>Update an App Runner service. You can update the source configuration and instance configuration of the service. You can also update the ARN of the auto scaling configuration resource that's associated with the service. However, you can't change the name or the encryption configuration of the service. These can be set only when you create the service.</p> <p>To update the tags applied to your service, use the separate actions <a>TagResource</a> and <a>UntagResource</a>.</p> <p>This is an asynchronous operation. On a successful call, you can use the returned <code>OperationId</code> and the <a>ListOperations</a> call to track the operation's progress.</p>

### x-amz-target-app-runner-update-vpc-ingress-connection

Manage x amz target app runner update vpc ingress connection

- **`aws-app-runner-pp-cli x-amz-target-app-runner-update-vpc-ingress-connection update-vpc-ingress-connection`** - <p>Update an existing App Runner VPC Ingress Connection resource. The VPC Ingress Connection must be in one of the following states to be updated:</p> <ul> <li> <p> AVAILABLE </p> </li> <li> <p> FAILED_CREATION </p> </li> <li> <p> FAILED_UPDATE </p> </li> </ul>


## Output Formats

```bash
# Human-readable table (default in terminal, JSON when piped)
aws-app-runner-pp-cli x-amz-target-app-runner-associate-custom-domain list

# JSON for scripting and agents
aws-app-runner-pp-cli x-amz-target-app-runner-associate-custom-domain list --json

# Filter to specific fields
aws-app-runner-pp-cli x-amz-target-app-runner-associate-custom-domain list --json --select id,name,status

# Dry run — show the request without sending
aws-app-runner-pp-cli x-amz-target-app-runner-associate-custom-domain list --dry-run

# Agent mode — JSON + compact + no prompts in one flag
aws-app-runner-pp-cli x-amz-target-app-runner-associate-custom-domain list --agent
```

## Agent Usage

This CLI is designed for AI agent consumption:

- **Non-interactive** - never prompts, every input is a flag
- **Pipeable** - `--json` output to stdout, errors to stderr
- **Filterable** - `--select id,name` returns only fields you need
- **Previewable** - `--dry-run` shows the request without sending
- **Retryable** - creates return "already exists" on retry, deletes return "already deleted"
- **Confirmable** - `--yes` for explicit confirmation of destructive actions
- **Piped input** - `echo '{"key":"value"}' | aws-app-runner-pp-cli <resource> create --stdin`
- **Cacheable** - GET responses cached for 5 minutes, bypass with `--no-cache`
- **Agent-safe by default** - no colors or formatting unless `--human-friendly` is set
- **Progress events** - paginated commands emit NDJSON events to stderr in default mode

Exit codes: `0` success, `2` usage error, `3` not found, `4` auth error, `5` API error, `7` rate limited, `10` config error.

## Use as MCP Server

This CLI ships a companion MCP server for use with Claude Desktop, Cursor, and other MCP-compatible tools.

### Claude Code

```bash
claude mcp add aws-app-runner aws-app-runner-pp-mcp -e AWS_APP_RUNNER_HMAC=<your-key>
```

### Claude Desktop

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "aws-app-runner": {
      "command": "aws-app-runner-pp-mcp",
      "env": {
        "AWS_APP_RUNNER_HMAC": "<your-key>"
      }
    }
  }
}
```

## Cookbook

Common workflows and recipes:

```bash
# List resources as JSON for scripting
aws-app-runner-pp-cli x-amz-target-app-runner-associate-custom-domain list --json

# Filter to specific fields
aws-app-runner-pp-cli x-amz-target-app-runner-associate-custom-domain list --json --select id,name,status

# Dry run to preview the request
aws-app-runner-pp-cli x-amz-target-app-runner-associate-custom-domain list --dry-run

# Sync data locally for offline search
aws-app-runner-pp-cli sync

# Search synced data
aws-app-runner-pp-cli search "query"

# Export for backup
aws-app-runner-pp-cli export --format jsonl > backup.jsonl
```

## Health Check

```bash
aws-app-runner-pp-cli doctor
```

<!-- DOCTOR_OUTPUT -->

## Configuration

Config file: `~/.config/aws-app-runner-pp-cli/config.toml`

Environment variables:
- `AWS_APP_RUNNER_HMAC`

## Troubleshooting

**Authentication errors (exit code 4)**
- Run `aws-app-runner-pp-cli doctor` to check credentials
- Verify the environment variable is set: `echo $AWS_APP_RUNNER_HMAC`

**Not found errors (exit code 3)**
- Check the resource ID is correct
- Run the `list` command to see available items

**Rate limit errors (exit code 7)**
- The CLI auto-retries with exponential backoff
- If persistent, wait a few minutes and try again

---

Generated by [CLI Printing Press](https://github.com/mvanhorn/cli-printing-press)
