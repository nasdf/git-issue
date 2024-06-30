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

var (
	createMessage      []string
	createFile         string
	createStripSpace   bool
	createNoStripSpace bool
	createAssign       []string
	createLabel        []string
	createStatus       string
)

func init() {
	createCommand.Flags().StringVarP(&createFile, "file", "F", "", "Take the issue message from the given file. Use - to read the issue message from the standard input. Lines starting with # and empty lines other than a single line between paragraphs will be stripped out. If you wish to keep them verbatim, use --no-stripspace.")
	createCommand.Flags().StringArrayVarP(&createMessage, "message", "m", []string{}, "Use the given issue message (instead of prompting). If multiple -m options are given, their values are concatenated as separate paragraphs. Lines starting with # and empty lines other than a single line between paragraphs will be stripped out. If you wish to keep them verbatim, use --no-stripspace.")
	createCommand.Flags().BoolVar(&createNoStripSpace, "no-stripspace", false, "")
	createCommand.Flags().BoolVar(&createStripSpace, "stripspace", true, "Strip leading and trailing whitespace from the issue message. Also strip out empty lines other than a single line between paragraphs.")
	createCommand.Flags().StringArrayVarP(&createAssign, "assignee", "a", []string{}, "Add an assignee to the issue. Multiple users can be assigned with multiple -a options.")
	createCommand.Flags().StringArrayVarP(&createLabel, "label", "l", []string{}, "Add a label to the issue. Multiple labels can be added with multiple -l options.")
	createCommand.Flags().StringVarP(&createStatus, "status", "s", "open", "Set the issue status. Defaults to open.")
	createCommand.MarkFlagsMutuallyExclusive("message", "file")
}

var createCommand = &cobra.Command{
	Use:   "create [--[no]-stripspace] [-F <file> | -m <message>] [-l <label>] [-a <user>] [-s <status>]",
	Short: "Create a new issue",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		issue, err := core.CreateIssue(cmd.Context(), createStatus, createAssign, createLabel)
		if err != nil {
			return err
		}
		var message io.Reader
		switch {
		case createFile == "-":
			// read message from stdin
			message = os.Stdin

		case createFile != "":
			// read message from file
			file, err := os.ReadFile(createFile)
			if err != nil {
				return err
			}
			message = bytes.NewBuffer(file)

		case len(createMessage) > 0:
			// read message from flags
			message = bytes.NewBufferString(strings.Join(createMessage, "\n\n"))

		default:
			// read message interactively
			edit, err := core.EditIssueMessage(cmd.Context(), issue)
			if err != nil {
				return err
			}
			message = bytes.NewBuffer(edit)
		}
		if createStripSpace && !createNoStripSpace {
			// strip comments from message
			clean, err := git.Exec(cmd.Context(), message, "stripspace", "--strip-comments")
			if err != nil {
				return err
			}
			issue.Message = string(clean)
		}
		err = core.AddIssueNote(cmd.Context(), issue)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), issue.Hash)
		return err
	},
}
