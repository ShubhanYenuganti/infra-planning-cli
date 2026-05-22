package smoke

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

var semverRE = regexp.MustCompile(`^v\d+\.\d+\.\d+`)

var requiredSubcommands = []string{
	"api", "auth", "doctor", "export", "import",
	"sync", "search", "sql", "version",
}

func TestCLISmoke(t *testing.T) {
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
			bin := BinaryPath(root, cli)

			t.Run("version-is-semver", func(t *testing.T) {
				out, err := exec.Command(bin, "--version").Output()
				if err != nil {
					t.Fatalf("--version failed: %v", err)
				}
				if !semverRE.Match(out) {
					t.Fatalf("--version output not SemVer: %q", out)
				}
			})

			t.Run("help-mentions-required-subcommands", func(t *testing.T) {
				out, err := exec.Command(bin, "--help").Output()
				if err != nil {
					t.Fatalf("--help failed: %v", err)
				}
				body := string(out)
				for _, sub := range requiredSubcommands {
					if !strings.Contains(body, sub) {
						t.Errorf("--help missing required subcommand %q", sub)
					}
				}
			})

			t.Run("unknown-command-exits-2", func(t *testing.T) {
				cmd := exec.Command(bin, "asdf-not-a-real-command")
				err := cmd.Run()
				exitErr, ok := err.(*exec.ExitError)
				if !ok {
					t.Fatalf("expected non-zero exit, got %v", err)
				}
				if exitErr.ExitCode() != 2 {
					t.Fatalf("expected exit code 2, got %d", exitErr.ExitCode())
				}
			})

			t.Run("doctor-json-parseable", func(t *testing.T) {
				// Accept any exit code; assert stdout parses as JSON.
				out, _ := exec.Command(bin, "doctor", "--json").Output()
				if len(out) == 0 {
					t.Fatal("doctor --json produced no stdout")
				}
				var v interface{}
				if err := json.Unmarshal(out, &v); err != nil {
					t.Fatalf("doctor --json output not parseable JSON: %v\noutput: %s", err, out)
				}
			})

			t.Run("agent-flag-honored", func(t *testing.T) {
				cmd := exec.Command(bin, "--agent", "doctor")
				err := cmd.Run()
				if exitErr, ok := err.(*exec.ExitError); ok {
					// Exit 2 = unknown flag/command; means --agent was rejected.
					if exitErr.ExitCode() == 2 {
						t.Fatalf("--agent treated as unknown (exit 2)")
					}
				}
			})
		})
	}
}
