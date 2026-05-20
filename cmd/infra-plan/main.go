package main

import (
	"os"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
