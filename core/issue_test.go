package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeIssue(t *testing.T) {
	author, err := DecodeSignature("Bob <bob@test.com> 1719340759 -0700")
	require.NoError(t, err)

	actual := Issue{
		Author:    *author,
		Assignees: []string{"Alice <alice@test.com>"},
		Labels:    []string{"bug", "feature"},
		Status:    IssueStatusOpen,
		Message:   "fix bug\n\nfix the bug in the code",
	}

	expect, err := DecodeIssue(actual.Encode())
	require.NoError(t, err)

	assert.Equal(t, expect.Author, actual.Author)
	assert.ElementsMatch(t, expect.Assignees, actual.Assignees)
	assert.ElementsMatch(t, expect.Labels, actual.Labels)
	assert.Equal(t, expect.Status, actual.Status)
	assert.Equal(t, expect.Message, actual.Message)
}
