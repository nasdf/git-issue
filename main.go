//go:build !docs

package main

import (
	"os"

	"github.com/nasdf/git-issue/command"
)

func main() {
	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
