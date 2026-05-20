package discovery

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type cloudformationDetector struct{}

func (c *cloudformationDetector) Name() string { return "cloudformation" }

func (c *cloudformationDetector) Matches(_, path string, _ fs.FileInfo) bool {
	switch filepath.Ext(path) {
	case ".yaml", ".yml", ".json":
		return true
	}
	return false
}

func (c *cloudformationDetector) Extract(_, path string, content []byte) []Evidence {
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "AWSTemplateFormatVersion") {
			return []Evidence{{
				Path:    path,
				Kind:    "aws-cloudformation",
				Snippet: firstCodeLine(content),
			}}
		}
	}
	return nil
}
