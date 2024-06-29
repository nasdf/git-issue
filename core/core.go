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
func CreateIssue(ctx context.Context, status string, assignees []string, labels []string) (*Issue, error) {
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
	blob := bytes.NewBufferString(fmt.Sprintf("issue %s", author.Encode()))
	hash, err := git.Exec(ctx, blob, "hash-object", "--stdin")
	if err != nil {
		return nil, err
	}
	return &Issue{
		Hash:      string(hash),
		Author:    author,
		Status:    status,
		Assignees: assignees,
		Labels:    labels,
	}, nil
}

// AddIssueNote adds a note that will contain the issue data.
func AddIssueNote(ctx context.Context, issue *Issue) error {
	blob := bytes.NewBufferString(fmt.Sprintf("issue %s", issue.Author.Encode()))
	hash, err := git.Exec(ctx, blob, "hash-object", "-w", "--stdin")
	if err != nil {
		return err
	}
	stdin := bytes.NewBufferString(issue.Encode())
	_, err = git.Exec(ctx, stdin, "notes", "--ref", issuesRef, "add", string(hash), "-F", "-")
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

// ListIssues returns a list containing all issues.
func ListIssues(ctx context.Context) (*IssueIterator, error) {
	notes, err := git.Exec(ctx, nil, "notes", "--ref", issuesRef, "list")
	if err != nil {
		return nil, err
	}
	var hashes []string
	for _, v := range bytes.Split(notes, []byte("\n")) {
		parts := bytes.Fields(v)
		if len(parts) != 2 {
			continue
		}
		hashes = append(hashes, string(parts[1]))
	}
	return &IssueIterator{hashes}, nil
}

// GetIssue returns the issue anchored to the object with the given hash.
func GetIssue(ctx context.Context, hash string) (*Issue, error) {
	data, err := git.Exec(ctx, nil, "notes", "--ref", issuesRef, "show", hash)
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
