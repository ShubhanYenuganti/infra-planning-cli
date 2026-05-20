package checks

import (
	"fmt"
	"path/filepath"
	"strings"
)

type FlagsCheck struct{}

func (*FlagsCheck) Name() string { return "required-flags" }

func (*FlagsCheck) Run(cliDir string) error {
	required := []string{"--json", "--dry-run", "--select", "--data-source", "--compact"}
	var contents string
	for _, dir := range []string{"cmd", "internal/cmd"} {
		c, err := concatGoSources(filepath.Join(cliDir, dir))
		if err == nil {
			contents += c
		}
	}
	if contents == "" {
		return fmt.Errorf("no Go source found in cmd/ or internal/cmd/")
	}
	for _, flag := range required {
		bare := strings.TrimPrefix(flag, "--")
		if !strings.Contains(contents, `"`+flag+`"`) &&
			!strings.Contains(contents, `"`+bare+`"`) {
			return fmt.Errorf("missing required flag: %s", flag)
		}
	}
	return nil
}
