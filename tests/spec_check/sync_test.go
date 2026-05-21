package spec_check

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSpecSync(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	cat, err := LoadCatalog(root)
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}
	if len(cat.CLIs) == 0 {
		t.Fatal("catalog has no CLIs")
	}
	for _, cli := range cat.CLIs {
		cli := cli
		t.Run(cli.Name, func(t *testing.T) {
			pf, err := LoadPressfile(root, cli.Path)
			if err != nil {
				t.Fatalf("load pressfile: %v", err)
			}
			if pf.SpecSHA256 == "" {
				t.Fatal("pressfile missing spec_sha256")
			}
			if pf.VendoredSpec == "" {
				t.Fatal("pressfile missing vendored_spec")
			}

			specBytes, err := os.ReadFile(filepath.Join(root, pf.VendoredSpec))
			if err != nil {
				t.Fatalf("read vendored spec: %v", err)
			}

			// 1) SHA-256 matches recorded value
			sum := sha256.Sum256(specBytes)
			actualHex := hex.EncodeToString(sum[:])
			if actualHex != pf.SpecSHA256 {
				t.Fatalf("spec SHA mismatch:\n  vendored file: %s\n  pressfile:     %s",
					actualHex, pf.SpecSHA256)
			}

			// 2) Vendored spec parses as JSON
			var doc map[string]interface{}
			if err := json.Unmarshal(specBytes, &doc); err != nil {
				t.Fatalf("vendored spec not parseable JSON: %v", err)
			}

			// 3) Spec advertises OpenAPI 3.x
			v, ok := doc["openapi"].(string)
			if !ok {
				t.Fatal("vendored spec missing 'openapi' field")
			}
			if !strings.HasPrefix(v, "3.") {
				t.Fatalf("vendored spec not OpenAPI 3.x (got openapi=%q)", v)
			}
		})
	}
}
