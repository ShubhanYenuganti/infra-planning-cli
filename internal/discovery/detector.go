package discovery

import "io/fs"

type Detector interface {
	Name() string
	Matches(repoRoot, path string, info fs.FileInfo) bool
	Extract(repoRoot, path string, content []byte) []Evidence
}
