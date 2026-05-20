package docs

import "testing"

func TestAWSDocsForDomain(t *testing.T) {
	got := AWSDocsForDomain("networking")
	if len(got) == 0 {
		t.Fatalf("expected networking docs")
	}
	if got[0].Title != "Amazon VPC documentation" {
		t.Fatalf("first title = %q", got[0].Title)
	}
	if got[0].URL == "" {
		t.Fatalf("first URL was empty")
	}
	if unknown := AWSDocsForDomain("unknown"); len(unknown) != 0 {
		t.Fatalf("unknown docs length = %d, want 0", len(unknown))
	}
}
