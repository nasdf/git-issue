package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeSignature(t *testing.T) {
	expect := &Signature{
		Name:  "Bob",
		Email: "bob@test.com",
		When:  time.Now(),
	}

	actual, err := DecodeSignature(expect.Encode())
	require.NoError(t, err)

	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Email, actual.Email)
	assert.Equal(t, expect.When.Unix(), actual.When.Unix())
}
