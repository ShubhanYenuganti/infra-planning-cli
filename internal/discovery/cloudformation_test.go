package discovery

import "testing"

func TestCloudFormationDetector(t *testing.T) {
	d := &cloudformationDetector{}

	yamlHit := "AWSTemplateFormatVersion: '2010-09-09'\nResources:\n  VPC:\n    Type: AWS::EC2::VPC"
	yamlMiss := "version: 2\njobs: {}"
	cases := []struct {
		name    string
		path    string
		content string
		wantHit bool
	}{
		{"cfn yaml", "template.yaml", yamlHit, true},
		{"cfn yml", "template.yml", yamlHit, true},
		{"non-cfn yaml", "ci.yaml", yamlMiss, false},
		{"non-yaml ext", "main.tf", yamlHit, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			matches := d.Matches("", tc.path, nil)
			if !matches && tc.wantHit {
				t.Fatalf("expected Matches=true")
			}
			// Only call Extract when Matches returns true (mirrors walker behavior).
			if !matches {
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
