package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nasdf/git-issue/core"
	"github.com/nasdf/git-issue/git"
)

func NewOpenCommand() *cobra.Command {
	var messageFlag []string
	var fileFlag string
	var stripSpaceFlag bool
	var noStripSpaceFlag bool
	cmd := &cobra.Command{
		Use:   "open [--[no]-stripspace] [-F <file> | -m <msg>]",
		Short: "Open a new issue",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			issue, err := core.CreateIssue(cmd.Context())
			if err != nil {
				return err
			}
			var message io.Reader
			switch {
			case fileFlag == "-":
				// read message from stdin
				message = os.Stdin

			case fileFlag != "":
				// read message from file
				file, err := os.ReadFile(fileFlag)
				if err != nil {
					return err
				}
				message = bytes.NewBuffer(file)

			case len(messageFlag) > 0:
				// read message from flags
				message = bytes.NewBufferString(strings.Join(messageFlag, "\n\n"))

			default:
				// read message interactively
				edit, err := core.EditIssueMessage(cmd.Context(), issue)
				if err != nil {
					return err
				}
				message = bytes.NewBuffer(edit)
			}
			if stripSpaceFlag && !noStripSpaceFlag {
				// strip comments from message
				clean, err := git.Exec(cmd.Context(), message, "stripspace", "--strip-comments")
				if err != nil {
					return err
				}
				issue.Message = string(clean)
			}
			err = core.CreateIssueNote(cmd.Context(), issue)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), issue.Hash)
			return err
		},
	}
	cmd.Flags().StringVarP(&fileFlag, "file", "F", "", "Take the issue message from the given file. Use - to read the issue message from the standard input. Lines starting with # and empty lines other than a single line between paragraphs will be stripped out. If you wish to keep them verbatim, use --no-stripspace.")
	cmd.Flags().StringArrayVarP(&messageFlag, "message", "m", []string{}, "Use the given issue message (instead of prompting). If multiple -m options are given, their values are concatenated as separate paragraphs. Lines starting with # and empty lines other than a single line between paragraphs will be stripped out. If you wish to keep them verbatim, use --no-stripspace.")
	cmd.Flags().BoolVar(&noStripSpaceFlag, "no-stripspace", false, "")
	cmd.Flags().BoolVar(&stripSpaceFlag, "stripspace", true, "Strip leading and trailing whitespace from the issue message. Also strip out empty lines other than a single line between paragraphs.")
	return cmd
}
