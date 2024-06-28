package command

import (
	"fmt"

	"github.com/nasdf/git-issue/core"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show <issue>",
		Short: "Show issue contents",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			issue, err := core.GetIssue(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), issue.String())
			return err
		},
	}
}
