package core

import (
	"fmt"
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
	oneline := strings.SplitN(i.Message, "\n", 2)[0]
	date := i.Author.When.Format("Mon Jan 02 15:04:05 2006")
	return fmt.Sprintf("issue %s\nAuthor: %s\nTime:   %s\nStatus: %s\n\n    %s\n", i.Hash, i.Author.String(), date, i.Status, oneline)
}
