package discovery

import "testing"

func TestSAMDetector(t *testing.T) {
	d := &samDetector{}

	samContent := "AWSTemplateFormatVersion: '2010-09-09'\nTransform: AWS::Serverless-2016-10-31\nResources:\n  Fn:\n    Type: AWS::Serverless::Function"
	cfnOnly := "AWSTemplateFormatVersion: '2010-09-09'\nResources: {}"

	if !d.Matches("", "template.yaml", nil) {
		t.Fatalf("expected Matches=true for template.yaml")
	}
	if d.Matches("", "stack.yaml", nil) {
		t.Fatalf("expected Matches=false for stack.yaml")
	}
	if ev := d.Extract("", "template.yaml", []byte(samContent)); len(ev) == 0 {
		t.Fatalf("expected SAM evidence")
	}
	if ev := d.Extract("", "template.yaml", []byte(cfnOnly)); len(ev) != 0 {
		t.Fatalf("expected no SAM evidence for plain CFN")
	}
}
