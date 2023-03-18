package server

import (
	"time"
)

func ParseQueryOffset(queryOffset string) (int64, error) {
	now := time.Now()
	if queryOffset == "now" {
		return now.Unix(), nil
	}

	duration := queryOffset[4:]

	timeOffset, err := time.ParseDuration(duration)

	if err != nil {
		return 0, err
	}

	return now.Add(-timeOffset).Unix(), nil
}
