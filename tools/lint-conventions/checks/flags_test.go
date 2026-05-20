package checks

import "testing"

func TestFlagsCheck_Pass(t *testing.T) {
	check := &FlagsCheck{}
	if err := check.Run("testdata/flags-ok"); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestFlagsCheck_MissingDryRun(t *testing.T) {
	check := &FlagsCheck{}
	if err := check.Run("testdata/flags-missing"); err == nil {
		t.Error("expected error when --dry-run missing")
	}
}
