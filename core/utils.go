package core

import (
	"fmt"
	"strconv"
	"time"
)

// parseTimezone parses the timezone from the given text.
func parseTimezone(text string) (*time.Location, error) {
	if len(text) != 5 {
		return nil, fmt.Errorf("invalid timezone format")
	}
	hours, err := strconv.ParseInt(text[0:3], 10, 64)
	if err != nil {
		return nil, err
	}
	minutes, err := strconv.ParseInt(text[3:], 10, 64)
	if err != nil {
		return nil, err
	}
	if hours < 0 {
		minutes *= -1
	}
	return time.FixedZone("", int(hours*60*60+minutes*60)), nil
}
