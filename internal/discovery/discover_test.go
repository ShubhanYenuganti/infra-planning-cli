package discovery

import (
	"sort"
	"testing"
)

func TestDiscoverRepoTerraformFixture(t *testing.T) {
	ctx, err := DiscoverRepo("../../testdata/repos/aws-terraform")
	if err != nil {
		t.Fatalf("DiscoverRepo error: %v", err)
	}
	assertDetected(t, ctx, "terraform")
	assertEvidencePath(t, ctx, "main.tf")
}

func TestDiscoverRepoCloudFormationFixture(t *testing.T) {
	ctx, err := DiscoverRepo("../../testdata/repos/aws-cloudformation")
	if err != nil {
		t.Fatalf("DiscoverRepo error: %v", err)
	}
	assertDetected(t, ctx, "cloudformation")
	assertEvidencePath(t, ctx, "template.yaml")
}

func TestDiscoverRepoSAMFixture(t *testing.T) {
	ctx, err := DiscoverRepo("../../testdata/repos/aws-sam")
	if err != nil {
		t.Fatalf("DiscoverRepo error: %v", err)
	}
	// SAM templates are CloudFormation templates with a Transform; both detectors fire.
	tools := append([]string{}, ctx.DetectedTools...)
	sort.Strings(tools)
	want := []string{"cloudformation", "sam"}
	if len(tools) != 2 || tools[0] != want[0] || tools[1] != want[1] {
		t.Fatalf("DetectedTools = %v, want %v", tools, want)
	}
}

func TestDiscoverRepoCDKFixture(t *testing.T) {
	ctx, err := DiscoverRepo("../../testdata/repos/aws-cdk")
	if err != nil {
		t.Fatalf("DiscoverRepo error: %v", err)
	}
	assertDetected(t, ctx, "cdk")
	assertEvidencePath(t, ctx, "cdk.json")
}

func TestDiscoverRepoServerlessFixture(t *testing.T) {
	ctx, err := DiscoverRepo("../../testdata/repos/aws-serverless")
	if err != nil {
		t.Fatalf("DiscoverRepo error: %v", err)
	}
	assertDetected(t, ctx, "serverless")
	assertEvidencePath(t, ctx, "serverless.yml")
}

func assertDetected(t *testing.T, ctx RepoContext, want string) {
	t.Helper()
	for _, n := range ctx.DetectedTools {
		if n == want {
			return
		}
	}
	t.Fatalf("DetectedTools = %v, want to include %q", ctx.DetectedTools, want)
}

func assertEvidencePath(t *testing.T, ctx RepoContext, wantPath string) {
	t.Helper()
	for _, ev := range ctx.Evidence {
		if ev.Path == wantPath {
			return
		}
	}
	t.Fatalf("no Evidence with Path=%q in %+v", wantPath, ctx.Evidence)
}
