package checks

import "testing"

func TestSQLitePathCheck_Pass(t *testing.T) {
	check := &SQLitePathCheck{}
	if err := check.Run("testdata/sqlite-ok"); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestSQLitePathCheck_WrongPath(t *testing.T) {
	check := &SQLitePathCheck{}
	if err := check.Run("testdata/sqlite-wrong"); err == nil {
		t.Error("expected error for non-conventional path")
	}
}
