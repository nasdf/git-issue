package main

import (
	"fmt"
	"os"

	"github.com/nasdf/git-issue/command"
)

func main() {
	err := command.NewRootCommand().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
