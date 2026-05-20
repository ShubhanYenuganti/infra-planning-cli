package discovery

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type terraformDetector struct{}

func (t *terraformDetector) Name() string { return "terraform" }

func (t *terraformDetector) Matches(_, path string, _ fs.FileInfo) bool {
	ext := filepath.Ext(path)
	return ext == ".tf" || ext == ".tfvars"
}

func (t *terraformDetector) Extract(_, path string, content []byte) []Evidence {
	s := string(content)
	lower := strings.ToLower(s)
	if !strings.Contains(s, "aws_") && !strings.Contains(lower, `provider "aws"`) {
		return nil
	}
	return []Evidence{{
		Path:    path,
		Kind:    "aws-iac",
		Snippet: firstCodeLine(content),
	}}
}
