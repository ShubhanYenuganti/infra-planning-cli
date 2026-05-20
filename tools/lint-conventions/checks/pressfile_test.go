package checks

import "testing"

func TestPressfileCheck_Pass(t *testing.T) {
	check := &PressfileCheck{}
	if err := check.Run("testdata/pressfile-valid"); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestPressfileCheck_MissingFile(t *testing.T) {
	check := &PressfileCheck{}
	if err := check.Run("testdata/pressfile-missing"); err == nil {
		t.Error("expected error for missing pressfile.yaml")
	}
}

func TestPressfileCheck_MalformedYAML(t *testing.T) {
	check := &PressfileCheck{}
	if err := check.Run("testdata/pressfile-malformed"); err == nil {
		t.Error("expected error for malformed yaml")
	}
}
