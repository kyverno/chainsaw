package main

import (
	"os"

	"github.com/kyverno/chainsaw/pkg/commands"
)

func main() {
	root := commands.RootCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
