package command

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		SilenceErrors: true,
		SilenceUsage:  true,
		Use:           "git-issue",
		Short:         "Open, close, and edit issues",
		Long:          "",
	}
	cmd.AddCommand(NewOpenCommand())
	cmd.AddCommand(NewShowCommand())
	return cmd
}
