package discovery

import "testing"

func TestCDKDetector(t *testing.T) {
	d := &cdkDetector{}
	if !d.Matches("", "cdk.json", nil) {
		t.Fatalf("expected Matches=true for cdk.json")
	}
	if d.Matches("", "package.json", nil) {
		t.Fatalf("expected Matches=false for package.json")
	}
	ev := d.Extract("", "cdk.json", []byte(`{"app":"npx ts-node bin/app.ts"}`))
	if len(ev) != 1 || ev[0].Kind != "aws-cdk" {
		t.Fatalf("Extract result = %+v", ev)
	}
}
