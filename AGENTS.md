# infra-press — agent guide

A library of focused, agent-native CLIs for cloud-infra APIs. Stateless
services and scheduled jobs across AWS, GCP, Azure.

## When to reach for what

| Need | CLI | Cloud |
|---|---|---|
| Deploy/manage a containerized service | cloud-run-admin-pp-cli | GCP |
| Deploy/manage a containerized service | apprunner-pp-cli | AWS |
| Deploy/manage a containerized service | container-apps-pp-cli | Azure |
| Run a serverless function | cloud-functions-pp-cli | GCP |
| Run a serverless function | lambda-pp-cli | AWS |
| Run a serverless function | functions-pp-cli | Azure |

For storage, IAM, networking, databases, secrets, and Kubernetes: fall back to
`aws`, `gcloud`, or `az`. See the "Known fallbacks" section below.

## Conventions every CLI in this library follows

- `--json` auto-on when stdout is piped
- Typed exit codes: `0` ok / `2` user error / `3` auth / `4` not found / `5` conflict / `7` rate-limit
- `sync` / `search` / `sql` subcommands on every CLI
- Local SQLite at `~/.infra-press/<cli>.db`
- `--data-source live|local|auto`
- `--dry-run` on every mutation
- `--select field1,field2` for field projection
- `--compact` drops to high-gravity fields only

Full spec: [`docs/conventions.md`](docs/conventions.md).

## Auth resolution

- **AWS:** flag → env (`AWS_PROFILE`) → SDK default chain
- **GCP:** flag → `GOOGLE_APPLICATION_CREDENTIALS` → ADC
- **Azure:** flag → `AZURE_TENANT_ID` + `AZURE_CLIENT_ID` + `AZURE_CLIENT_SECRET` → DefaultAzureCredential

Run the relevant doctor script before first use:
```bash
./scripts/doctor-aws.sh
./scripts/doctor-gcp.sh
./scripts/doctor-azure.sh
```

## Known fallbacks to vendor CLIs

- **Azure Functions** needs a Storage Account first:
  `az storage account create -n <name> -g <rg> -l <region> --sku Standard_LRS`
- **AWS Lambda** first-time use needs an IAM execution role:
  use `lambda functions create --auto-role` (gap patch) or pre-create
  via `aws iam create-role ...`

## Cross-CLI compound recipes

See [`docs/compound-recipes.md`](docs/compound-recipes.md) for SQL/scripts that
join data across the SQLite stores (stale services across all clouds, public
endpoints without auth, idle-compute spend triangulation, etc.).
