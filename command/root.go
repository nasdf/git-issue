package command

import (
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "git-issue",
	Short: "Create, edit, or list issues",
	Long:  "",
}

func init() {
	rootCommand.AddCommand(createCommand)
	rootCommand.AddCommand(showCommand)
}

func Execute() error {
	return rootCommand.Execute()
}
