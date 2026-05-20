package discovery

import "testing"

func TestTerraformDetector(t *testing.T) {
	d := &terraformDetector{}

	cases := []struct {
		name    string
		path    string
		content string
		wantHit bool
	}{
		{"aws tf", "main.tf", `resource "aws_vpc" "main" {}`, true},
		{"non-aws tf", "main.tf", `resource "google_compute_instance" "x" {}`, false},
		{"non-tf ext", "main.yaml", `resource "aws_vpc" "main" {}`, false},
		{"tfvars aws", "vars.tfvars", "region = \"us-east-1\"\naws_role = \"arn:aws:iam::...\"", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			matches := d.Matches("", tc.path, nil)
			if matches != (filepathExtIsTfOrTfvars(tc.path)) {
				t.Fatalf("Matches = %v", matches)
			}
			// Only call Extract when Matches returns true (mirrors walker behavior).
			if !matches {
				if tc.wantHit {
					t.Fatalf("expected wantHit=true but Matches=false")
				}
				return
			}
			ev := d.Extract("", tc.path, []byte(tc.content))
			if tc.wantHit && len(ev) == 0 {
				t.Fatalf("expected evidence, got none")
			}
			if !tc.wantHit && len(ev) > 0 {
				t.Fatalf("expected no evidence, got %+v", ev)
			}
		})
	}
}

func filepathExtIsTfOrTfvars(p string) bool {
	return len(p) > 3 && (p[len(p)-3:] == ".tf" || (len(p) > 7 && p[len(p)-7:] == ".tfvars"))
}
