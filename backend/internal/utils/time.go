package utils

import (
	"errors"
	"time"
)

func ParseTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New("empty time string")
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("invalid time format")
}
