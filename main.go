package main

import (
	"os"

	"github.com/rarimo/voting-relayer/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
