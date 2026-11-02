package main

import (
	"os"

	"github.com/scryer-labs/scryer-java-engine/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
