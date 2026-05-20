package discovery

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type samDetector struct{}

func (s *samDetector) Name() string { return "sam" }

func (s *samDetector) Matches(_, path string, _ fs.FileInfo) bool {
	base := filepath.Base(path)
	return base == "template.yaml" || base == "template.yml" || base == "template.json"
}

func (s *samDetector) Extract(_, path string, content []byte) []Evidence {
	if !strings.Contains(string(content), "Transform: AWS::Serverless") {
		return nil
	}
	return []Evidence{{
		Path:    path,
		Kind:    "aws-sam",
		Snippet: firstCodeLine(content),
	}}
}
