package main

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/kyverno/chainsaw/pkg/commands"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func main() {
	log.SetLogger(logr.Discard())
	root := commands.RootCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
