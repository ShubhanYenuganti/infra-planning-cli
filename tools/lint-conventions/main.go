package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ShubhanYenuganti/infra-press/tools/lint-conventions/checks"
)

func main() {
	ci := flag.Bool("ci", false, "Exit non-zero on any violation (suitable for CI)")
	flag.Parse()

	target := flag.Arg(0)
	if target == "" {
		fmt.Fprintln(os.Stderr, "usage: lint-conventions [--ci] <path>")
		os.Exit(2)
	}

	cliDirs, err := discoverCLIDirs(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "discovery failed: %v\n", err)
		os.Exit(1)
	}

	violations := 0
	for _, dir := range cliDirs {
		for _, check := range checks.Registry() {
			if err := check.Run(dir); err != nil {
				fmt.Printf("FAIL  %s  %s: %v\n", dir, check.Name(), err)
				violations++
			} else {
				fmt.Printf("PASS  %s  %s\n", dir, check.Name())
			}
		}
	}

	if violations > 0 && *ci {
		os.Exit(1)
	}
}

// discoverCLIDirs returns CLI dirs under target. If target itself is a
// CLI dir (has pressfile.yaml), returns just it; otherwise walks library/.
func discoverCLIDirs(target string) ([]string, error) {
	if _, err := os.Stat(filepath.Join(target, "pressfile.yaml")); err == nil {
		return []string{target}, nil
	}
	var dirs []string
	err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if _, err := os.Stat(filepath.Join(path, "pressfile.yaml")); err == nil {
				dirs = append(dirs, path)
			}
		}
		return nil
	})
	return dirs, err
}
