package discovery

import (
	"io/fs"
	"path/filepath"
	"regexp"
)

type serverlessDetector struct{}

var serverlessProviderRE = regexp.MustCompile(`(?m)^\s*provider:\s*aws\s*$|^\s*name:\s*aws\s*$`)

func (s *serverlessDetector) Name() string { return "serverless" }

func (s *serverlessDetector) Matches(_, path string, _ fs.FileInfo) bool {
	base := filepath.Base(path)
	return base == "serverless.yml" || base == "serverless.yaml"
}

func (s *serverlessDetector) Extract(_, path string, content []byte) []Evidence {
	if !serverlessProviderRE.Match(content) {
		return nil
	}
	return []Evidence{{
		Path:    path,
		Kind:    "aws-serverless",
		Snippet: firstCodeLine(content),
	}}
}
