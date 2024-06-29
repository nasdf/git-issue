//go:build docs

package main

import (
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/nasdf/git-issue/command"
)

const docsDir = "docs"

func main() {
	err := os.RemoveAll(docsDir)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(docsDir, 0755)
	if err != nil {
		panic(err)
	}
	err = doc.GenMarkdownTree(command.NewRootCommand(), docsDir)
	if err != nil {
		panic(err)
	}
}
