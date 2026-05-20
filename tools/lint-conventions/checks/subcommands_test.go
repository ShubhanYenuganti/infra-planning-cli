package checks

import "testing"

func TestSubcommandsCheck_Pass(t *testing.T) {
	check := &SubcommandsCheck{}
	if err := check.Run("testdata/subcommands-ok"); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestSubcommandsCheck_MissingSync(t *testing.T) {
	check := &SubcommandsCheck{}
	if err := check.Run("testdata/subcommands-missing"); err == nil {
		t.Error("expected error when sync command missing")
	}
}
