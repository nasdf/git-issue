package core

import (
	"context"
	"fmt"
	"io"
	"strings"
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
	headers = append(headers, fmt.Sprintf("issue %s", strings.Join(append([]string{i.Hash}, i.Labels...), " ")))
	headers = append(headers, fmt.Sprintf("Author:   %s", i.Author.String()))
	headers = append(headers, fmt.Sprintf("Time:     %s", i.Author.When.Format("Mon Jan 02 15:04:05 2006")))
	headers = append(headers, fmt.Sprintf("Status:   %s", i.Status))

	for _, a := range i.Assignees {
		headers = append(headers, fmt.Sprintf("Assignee: %s", a))
	}

	oneline := strings.SplitN(i.Message, "\n", 2)[0]
	return fmt.Sprintf("%s\n\n    %s\n", strings.Join(headers, "\n"), oneline)
}

type IssueFilter struct {
	assignees []string
	labels    []string
	status    []string
}

func NewIssueFilter(assignees []string, labels []string, status []string) *IssueFilter {
	return &IssueFilter{
		assignees: assignees,
		labels:    labels,
		status:    status,
	}
}

func (i *IssueFilter) Match(issue *Issue) bool {
	matchAssignee := len(i.assignees) == 0
	for _, a := range i.assignees {
		for _, b := range issue.Assignees {
			matchAssignee = matchAssignee || a == b
		}
	}
	matchLabel := len(i.labels) == 0
	for _, a := range i.labels {
		for _, b := range issue.Labels {
			matchLabel = matchLabel || a == b
		}
	}
	matchStatus := len(i.status) == 0
	for _, a := range i.status {
		matchLabel = matchLabel || a == issue.Status
	}
	return matchAssignee && matchLabel && matchStatus
}

type IssueIterator struct {
	hashes []string
	filter *IssueFilter
	issue  *Issue
}

func NewIssueIterator(ctx context.Context, hashes []string, filter *IssueFilter) (*IssueIterator, error) {
	iter := &IssueIterator{
		hashes: hashes,
		filter: filter,
	}
	if len(hashes) == 0 {
		return iter, nil
	}
	if err := iter.Next(ctx); err != nil {
		return nil, err
	}
	return iter, nil
}

func (i *IssueIterator) Value() *Issue {
	return i.issue
}

func (i *IssueIterator) HasNext() bool {
	return i.issue != nil
}

func (i *IssueIterator) Next(ctx context.Context) error {
	for {
		if len(i.hashes) == 0 {
			return io.EOF
		}
		issue, err := GetIssue(ctx, i.hashes[0])
		if err != nil {
			return err
		}
		i.hashes = i.hashes[1:]
		i.issue = issue
		if i.filter.Match(issue) {
			return nil
		}
	}
}

func (i *IssueIterator) ForEach(ctx context.Context, fn func(*Issue) error) error {
	for {
		err := i.Next(ctx)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := fn(i.Value()); err != nil {
			return err
		}
	}
}
