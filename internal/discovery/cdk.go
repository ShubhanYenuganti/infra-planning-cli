package discovery

import (
	"io/fs"
	"path/filepath"
)

type cdkDetector struct{}

func (c *cdkDetector) Name() string { return "cdk" }

func (c *cdkDetector) Matches(_, path string, _ fs.FileInfo) bool {
	return filepath.Base(path) == "cdk.json"
}

func (c *cdkDetector) Extract(_, path string, content []byte) []Evidence {
	return []Evidence{{
		Path:    path,
		Kind:    "aws-cdk",
		Snippet: firstCodeLine(content),
	}}
}
