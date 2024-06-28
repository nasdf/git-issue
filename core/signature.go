package core

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var signaturePattern = regexp.MustCompile(`(?P<name>.*) <(?P<email>.*)> (?P<unix>\d+) (?P<zone>[+|-]\d\d\d\d)`)

type Signature struct {
	// Name is the name of the user.
	Name string
	// Email is the email of the user.
	Email string
	// When is the time the signature was created.
	When time.Time
}

// DecodeSignature decodes the given text into a signature.
func DecodeSignature(text string) (*Signature, error) {
	matches := signaturePattern.FindStringSubmatch(text)
	if len(matches) != 5 {
		return nil, fmt.Errorf("invalid author format")
	}
	unix, err := strconv.ParseInt(matches[3], 10, 64)
	if err != nil {
		return nil, err
	}
	zone, err := parseTimezone(matches[4])
	if err != nil {
		return nil, err
	}
	return &Signature{
		Name:  matches[1],
		Email: matches[2],
		When:  time.Unix(unix, 0).In(zone),
	}, nil
}

// Encode returns the signature header encoding.
func (s *Signature) Encode() string {
	return fmt.Sprintf("%s <%s> %d %s", s.Name, s.Email, s.When.Unix(), s.When.Format("-0700"))
}

// String returns the signature human friendly encoding.
func (s *Signature) String() string {
	return fmt.Sprintf("%s <%s>", s.Name, s.Email)
}
