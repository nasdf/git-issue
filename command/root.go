package command

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git-issue",
		Short: "Create, edit, or list issues",
		Long:  "",
	}
	cmd.AddCommand(NewCreateCommand())
	cmd.AddCommand(NewShowCommand())
	return cmd
}
