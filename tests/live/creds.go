//go:build live

package live

import "os"

// HasCreds returns true if env vars indicating OIDC creds for the given
// cloud are present. The vars are set by the canonical GitHub Actions
// configure-creds actions for each cloud.
func HasCreds(cloud string) bool {
	switch cloud {
	case "aws":
		return os.Getenv("AWS_ROLE_ARN") != "" && os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE") != ""
	case "gcp":
		return os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" ||
			os.Getenv("GCP_WORKLOAD_IDENTITY_PROVIDER") != ""
	case "azure":
		return os.Getenv("AZURE_CLIENT_ID") != "" && os.Getenv("AZURE_TENANT_ID") != ""
	default:
		return false
	}
}
