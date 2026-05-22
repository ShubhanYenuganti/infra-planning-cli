//go:build live

package live

import (
	"encoding/json"
	"os/exec"
	"testing"
)

func TestLiveCloud(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	cat, err := LoadCatalog(root)
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}
	for _, cli := range cat.CLIs {
		cli := cli
		t.Run(cli.Name, func(t *testing.T) {
			if !HasCreds(cli.Cloud) {
				t.Skipf("no OIDC creds for %s in env", cli.Cloud)
			}
			bin := BinaryPath(root, cli)

			t.Run("doctor-reaches-cloud", func(t *testing.T) {
				out, err := exec.Command(bin, "doctor", "--json").Output()
				if err != nil {
					t.Fatalf("doctor --json failed: %v\noutput: %s", err, out)
				}
				var report map[string]interface{}
				if err := json.Unmarshal(out, &report); err != nil {
					t.Fatalf("doctor --json not parseable: %v\noutput: %s", err, out)
				}
				// Convention: doctor reports "status": "ok" on full reachability.
				if status, _ := report["status"].(string); status != "ok" {
					t.Fatalf("doctor reports non-ok status: %v", report)
				}
			})

			if cli.LiveSmoke.ListSubcommand == "" {
				t.Log("no list_subcommand declared; skipping list-endpoint assertion")
				return
			}

			t.Run("list-subcommand-returns-json", func(t *testing.T) {
				args := append([]string{cli.LiveSmoke.ListSubcommand, "--json"}, cli.LiveSmoke.Args...)
				out, err := exec.Command(bin, args...).Output()
				if err != nil {
					t.Fatalf("%s failed: %v\noutput: %s", cli.LiveSmoke.ListSubcommand, err, out)
				}
				var v interface{}
				if err := json.Unmarshal(out, &v); err != nil {
					t.Fatalf("%s output not parseable JSON: %v", cli.LiveSmoke.ListSubcommand, err)
				}
			})
		})
	}
}
