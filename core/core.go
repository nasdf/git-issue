package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nasdf/git-issue/git"
)

const (
	issuesEditFile     = "ISSUES_EDITMSG"
	issuesRef          = "refs/notes/issues"
	issuesEditTemplate = "%s\n\n#\n# Write/edit the message for the following issue:\n#\n# %s"
)

// CreateIssue returns a new issue with the author set to the current user.
func CreateIssue(ctx context.Context) (*Issue, error) {
	authorName, err := git.Exec(ctx, nil, "config", "user.name")
	if err != nil {
		return nil, err
	}
	authorEmail, err := git.Exec(ctx, nil, "config", "user.email")
	if err != nil {
		return nil, err
	}
	author := Signature{
		Name:  string(authorName),
		Email: string(authorEmail),
		When:  time.Now(),
	}

	// create an object to anchor the issue note to
	blob := fmt.Sprintf("issue %s", author)
	hash, err := git.Exec(ctx, bytes.NewBufferString(blob), "hash-object", "-w", "--stdin")
	if err != nil {
		return nil, err
	}
	return &Issue{
		Hash:   string(hash),
		Author: author,
		Status: IssueStatusOpen,
	}, nil
}

// CreateIssueNote creates a note that will contain the issue data.
func CreateIssueNote(ctx context.Context, issue *Issue) error {
	data := issue.Encode()
	stdin := bytes.NewBufferString(data)
	_, err := git.Exec(ctx, stdin, "notes", "--ref", issuesRef, "add", issue.Hash, "-F", "-")
	return err
}

// EditIssueMessage opens the issue message in an interactive editor.
func EditIssueMessage(ctx context.Context, issue *Issue) ([]byte, error) {
	dir, err := git.Exec(ctx, nil, "rev-parse", "--git-dir")
	if err != nil {
		return nil, err
	}
	file := filepath.Join(string(dir), issuesEditFile)
	temp := fmt.Sprintf(issuesEditTemplate, issue.Message, issue.Hash)
	if err := os.WriteFile(file, []byte(temp), 0755); err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, os.Remove(file))
	}()
	if err = git.LaunchEditor(ctx, file); err != nil {
		return nil, err
	}
	return os.ReadFile(file)
}

func ListIssueNotes(ctx context.Context) ([][]byte, error) {
	notes, err := git.Exec(ctx, nil, "notes", "--ref", issuesRef, "list")
	if err != nil {
		return nil, err
	}
	return bytes.Split(notes, []byte("\n")), nil
}

func ListIssues(ctx context.Context) ([]*Issue, error) {
	notes, err := ListIssueNotes(ctx)
	if err != nil {
		return nil, err
	}
	var issues []*Issue
	for _, v := range notes {
		parts := bytes.Fields(v)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid note list format")
		}
		issue, err := GetIssue(ctx, string(parts[0]))
		if err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

func GetIssue(ctx context.Context, hash string) (*Issue, error) {
	data, err := git.Exec(ctx, nil, "cat-file", "blob", hash)
	if err != nil {
		return nil, err
	}
	issue, err := DecodeIssue(string(data))
	if err != nil {
		return nil, err
	}
	issue.Hash = hash
	return issue, nil
}
