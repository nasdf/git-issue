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
			var issues []*core.Issue
			switch {
			case len(args) == 0:
				list, err := core.ListIssues(cmd.Context())
				if err != nil {
					return err
				}
				issues = append(issues, list...)

			default:
				issue, err := core.GetIssue(cmd.Context(), args[0])
				if err != nil {
					return err
				}
				issues = append(issues, issue)
			}

			if len(issues) == 0 {
				return nil
			}

			pr, pw := io.Pipe()
			go func() {
				for _, i := range issues {
					fmt.Fprintln(pw, i.String())
				}
				pw.Close()
			}()

			switch {
			case len(issues) > 3:
				// pipe output to pager program
				ok, err := git.Pager(cmd.Context(), pr)
				if err != nil || ok {
					return err
				}
				fallthrough

			default:
				// print issues to stdout
				_, err := io.Copy(cmd.OutOrStdout(), pr)
				return err
			}
		},
	}
}
