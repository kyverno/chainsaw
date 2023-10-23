package main

import (
	"errors"
	"os"

	"github.com/kyverno/chainsaw/pkg/commands"
	"github.com/spf13/cobra/doc"
)

const path = "./website/docs/commands"

func identity(s string) string {
	return s
}

func empty(s string) string {
	return ""
}

func main() {
	root := commands.RootCommand()
	prepender := empty
	linkHandler := identity
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path, os.ModeDir|os.ModePerm); err != nil {
			panic(err)
		}
	}
	root.DisableAutoGenTag = true
	if err := doc.GenMarkdownTreeCustom(root, path, prepender, linkHandler); err != nil {
		panic(err)
	}
}
