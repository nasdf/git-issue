package core

import (
	"context"
	"fmt"
	"io"
	"strings"
)

const (
	IssueStatusOpen   = "open"
	IssueStatusClosed = "closed"
)

type Issue struct {
	// Hash is the hash of the object the issue note is anchored to
	Hash string
	// Author is the creator of the issue
	Author Signature
	// Assignees is a list of assigned users
	Assignees []string
	// Labels is a list of labels used to filter issues
	Labels []string
	// Status is the issue status (open or closed)
	Status string
	// Message is the issue description message
	Message string
}

// DecodeIssue decodes the given text into an issue.
func DecodeIssue(text string) (*Issue, error) {
	parts := strings.SplitN(text, "\n\n", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid issue format")
	}
	issue := Issue{
		Message: parts[1],
	}

	for _, v := range strings.Split(parts[0], "\n") {
		switch {
		case strings.HasPrefix(v, "author"):
			author, err := DecodeSignature(strings.TrimSpace(v[6:]))
			if err != nil {
				return nil, err
			}
			issue.Author = *author

		case strings.HasPrefix(v, "status"):
			issue.Status = strings.TrimSpace(v[6:])

		case strings.HasPrefix(v, "assignee"):
			issue.Assignees = append(issue.Assignees, strings.TrimSpace(v[8:]))

		case strings.HasPrefix(v, "label"):
			issue.Labels = append(issue.Labels, strings.TrimSpace(v[5:]))

		default:
			return nil, fmt.Errorf("invalid issue format")
		}
	}

	return &issue, nil
}

// Encode returns the encoded issue headers and message.
func (i *Issue) Encode() string {
	var headers []string
	headers = append(headers, fmt.Sprintf("author %s", i.Author.Encode()))
	headers = append(headers, fmt.Sprintf("status %s", i.Status))

	assignees := make(map[string]string)
	for _, v := range i.Assignees {
		assignees[v] = fmt.Sprintf("assignee %s", v)
	}
	for _, v := range assignees {
		headers = append(headers, v)
	}

	labels := make(map[string]string)
	for _, v := range i.Labels {
		labels[v] = fmt.Sprintf("label %s", v)
	}
	for _, v := range labels {
		headers = append(headers, v)
	}

	return fmt.Sprintf("%s\n\n%s", strings.Join(headers, "\n"), i.Message)
}

// String returns the issue encoded as a human friendly string.
func (i *Issue) String() string {
	var headers []string
	headers = append(headers, fmt.Sprintf("issue %s", i.Hash))
	headers = append(headers, fmt.Sprintf("Author: %s", i.Author.String()))
	headers = append(headers, fmt.Sprintf("Time:   %s", i.Author.When.Format("Mon Jan 02 15:04:05 2006")))
	headers = append(headers, fmt.Sprintf("Status: %s", i.Status))
	oneline := strings.SplitN(i.Message, "\n", 2)[0]
	return fmt.Sprintf("%s\n\n    %s\n", strings.Join(headers, "\n"), oneline)
}

type IssueIterator struct {
	hashes []string
}

func (i *IssueIterator) HasNext() bool {
	return len(i.hashes) > 0
}

func (i *IssueIterator) Next(ctx context.Context) (*Issue, error) {
	if len(i.hashes) == 0 {
		return nil, io.EOF
	}
	issue, err := GetIssue(ctx, i.hashes[0])
	if err != nil {
		return nil, err
	}
	i.hashes = i.hashes[1:]
	return issue, nil
}

func (i *IssueIterator) ForEach(ctx context.Context, fn func(*Issue) error) error {
	for {
		issue, err := i.Next(ctx)
		if err != nil {
			return err
		}
		if err := fn(issue); err != nil {
			return err
		}
	}
}
