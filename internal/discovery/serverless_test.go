package discovery

import "testing"

func TestServerlessDetector(t *testing.T) {
	d := &serverlessDetector{}

	flat := "service: hello\nprovider: aws\nfunctions: {}"
	nested := "service: hello\nprovider:\n  name: aws\n  runtime: nodejs20.x"
	googleProvider := "service: hello\nprovider:\n  name: google\n  runtime: nodejs20.x"

	if !d.Matches("", "serverless.yml", nil) {
		t.Fatalf("expected Matches=true")
	}
	if d.Matches("", "compose.yml", nil) {
		t.Fatalf("expected Matches=false")
	}
	if ev := d.Extract("", "serverless.yml", []byte(flat)); len(ev) == 0 {
		t.Fatalf("expected evidence for flat provider")
	}
	if ev := d.Extract("", "serverless.yml", []byte(nested)); len(ev) == 0 {
		t.Fatalf("expected evidence for nested provider")
	}
	if ev := d.Extract("", "serverless.yml", []byte(googleProvider)); len(ev) != 0 {
		t.Fatalf("expected no evidence for google provider")
	}
}
