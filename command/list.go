package command

import (
	"fmt"
	"io"

	"github.com/nasdf/git-issue/core"
	"github.com/nasdf/git-issue/git"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			issues, err := core.ListIssues(cmd.Context())
			if err != nil {
				return err
			}
			// buffer output by writing to a pipe
			pr, pw := io.Pipe()
			go func() {
				for _, i := range issues {
					fmt.Fprintln(pw, i.String())
				}
				pw.Close()
			}()
			ok, err := git.Pager(cmd.Context(), pr)
			if err != nil || ok {
				return err
			}
			_, err = io.Copy(cmd.OutOrStdout(), pr)
			return err
		},
	}
	return cmd
}
