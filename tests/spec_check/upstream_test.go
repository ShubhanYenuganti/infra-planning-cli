package spec_check

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

type finding struct {
	cli      string
	priority string // "high" if op present in vendored, "info" if upstream-only.
	message  string
}

func TestUpstreamAdvisory(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	cat, err := LoadCatalog(root)
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}

	var (
		mu       sync.Mutex
		findings []finding
		wg       sync.WaitGroup
	)
	client := &http.Client{Timeout: 5 * time.Second}

	for _, cli := range cat.CLIs {
		cli := cli
		wg.Add(1)
		go func() {
			defer wg.Done()

			pf, err := LoadPressfile(root, cli.Path)
			if err != nil {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "info",
					"could not load pressfile: " + err.Error()})
				mu.Unlock()
				return
			}
			if pf.SpecURL == "" {
				return
			}

			// Fetch upstream with timeout.
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			req, _ := http.NewRequestWithContext(ctx, "GET", pf.SpecURL, nil)
			resp, err := client.Do(req)
			if err != nil {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "info",
					"upstream unreachable: " + err.Error() + " (advisory only)"})
				mu.Unlock()
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == 404 {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "high",
					"spec URL returned 404 — consider re-pinning"})
				mu.Unlock()
				return
			}
			if resp.StatusCode != 200 {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "info",
					"upstream HTTP " + resp.Status + " (advisory only)"})
				mu.Unlock()
				return
			}

			// Load both specs through kin-openapi for operation-id awareness.
			loader := openapi3.NewLoader()
			loader.IsExternalRefsAllowed = true

			liveDoc, err := loader.LoadFromIoReader(resp.Body)
			if err != nil {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "info",
					"upstream spec did not parse as OpenAPI 3.x (advisory only)"})
				mu.Unlock()
				return
			}

			vendoredBytes, err := os.ReadFile(filepath.Join(root, pf.VendoredSpec))
			if err != nil {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "info",
					"vendored spec missing: " + err.Error()})
				mu.Unlock()
				return
			}
			vendoredDoc, err := loader.LoadFromData(vendoredBytes)
			if err != nil {
				mu.Lock()
				findings = append(findings, finding{cli.Name, "info",
					"vendored spec did not parse (this is a separate bug)"})
				mu.Unlock()
				return
			}

			// Compute operation-set diff.
			liveOps := operationSet(liveDoc)
			vendoredOps := operationSet(vendoredDoc)

			var removed, added []string
			for op := range vendoredOps {
				if _, ok := liveOps[op]; !ok {
					removed = append(removed, op)
				}
			}
			for op := range liveOps {
				if _, ok := vendoredOps[op]; !ok {
					added = append(added, op)
				}
			}
			sort.Strings(removed)
			sort.Strings(added)

			mu.Lock()
			defer mu.Unlock()
			for _, op := range removed {
				findings = append(findings, finding{cli.Name, "high",
					"operation removed upstream: " + op + " (CLI exposes it; consider regen)"})
			}
			if len(added) > 0 {
				findings = append(findings, finding{cli.Name, "info",
					"new endpoints available upstream (not exposed): " +
						strings.Join(added, ", ")})
			}
		}()
	}
	wg.Wait()

	// Sort: high-priority findings first, then info.
	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].priority != findings[j].priority {
			return findings[i].priority == "high"
		}
		return findings[i].cli < findings[j].cli
	})

	if len(findings) == 0 {
		t.Log("upstream advisory: no findings (clean)")
		return
	}

	t.Log("upstream advisory findings (ranked):")
	for _, f := range findings {
		t.Logf("  [%s] %s: %s", f.priority, f.cli, f.message)
	}
	// L7 NEVER fails the PR.
}

// operationSet builds a "METHOD path::operationId" set from an OpenAPI 3 doc.
func operationSet(doc *openapi3.T) map[string]struct{} {
	out := map[string]struct{}{}
	if doc == nil || doc.Paths == nil {
		return out
	}
	for path, item := range doc.Paths.Map() {
		for method, op := range item.Operations() {
			id := op.OperationID
			if id == "" {
				id = "<no-op-id>"
			}
			out[method+" "+path+"::"+id] = struct{}{}
		}
	}
	return out
}
