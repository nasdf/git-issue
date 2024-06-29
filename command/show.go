package command

import (
	"fmt"
	"io"

	"github.com/nasdf/git-issue/core"
	"github.com/nasdf/git-issue/git"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show [<issue>]",
		Short: "Show issue contents",
		Long:  "",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch {
			case len(args) == 0:
				iter, err := core.ListIssues(cmd.Context())
				if err != nil {
					return err
				}
				if !iter.HasNext() {
					return nil
				}
				pr, pw := io.Pipe()
				go func() {
					err := iter.ForEach(cmd.Context(), func(i *core.Issue) error {
						_, err := fmt.Fprintln(pw, i.String())
						return err
					})
					pw.CloseWithError(err)
				}()
				// pipe output to pager program
				ok, err := git.Pager(cmd.Context(), pr)
				if err != nil || ok {
					return err
				}
				// print to stdout if no pager available
				_, err = io.Copy(cmd.OutOrStdout(), pr)
				return err

			default:
				issue, err := core.GetIssue(cmd.Context(), args[0])
				if err != nil {
					return err
				}
				_, err = fmt.Fprint(cmd.OutOrStdout(), issue.String())
				return err
			}
		},
	}
}
