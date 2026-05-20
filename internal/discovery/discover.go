package discovery

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const maxFileBytes = 1 << 20 // 1 MiB

var skipDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	".terraform":   true,
	"cdk.out":      true,
}

func DiscoverRepo(root string) (RepoContext, error) {
	return DiscoverRepoWith(root, defaultDetectors())
}

func DiscoverRepoWith(root string, detectors []Detector) (RepoContext, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return RepoContext{}, err
	}

	ctx := RepoContext{Root: abs}
	fired := map[string]bool{}

	walkErr := filepath.WalkDir(abs, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if skipDirs[entry.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		if info.Size() > maxFileBytes {
			return nil
		}

		var content []byte
		for _, d := range detectors {
			if !d.Matches(abs, path, info) {
				continue
			}
			if content == nil {
				content, err = os.ReadFile(path)
				if err != nil {
					return nil
				}
			}
			records := d.Extract(abs, path, content)
			if len(records) == 0 {
				continue
			}
			if !fired[d.Name()] {
				fired[d.Name()] = true
				ctx.DetectedTools = append(ctx.DetectedTools, d.Name())
			}
			for _, ev := range records {
				if rel, err := filepath.Rel(abs, ev.Path); err == nil && !strings.HasPrefix(rel, "..") {
					ev.Path = rel
				}
				ctx.Evidence = append(ctx.Evidence, ev)
			}
		}
		return nil
	})
	if walkErr != nil {
		return RepoContext{}, walkErr
	}

	return ctx, nil
}

func defaultDetectors() []Detector {
	return []Detector{
		&terraformDetector{},
		&cloudformationDetector{},
		&samDetector{},
		&cdkDetector{},
		&serverlessDetector{},
	}
}

func firstCodeLine(content []byte) string {
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			continue
		}
		if len(trimmed) > 120 {
			trimmed = trimmed[:120]
		}
		return trimmed
	}
	return ""
}
