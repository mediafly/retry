package retry

import (
	"time"
)

func minDuration(a time.Duration, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func maxDuration(a time.Duration, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
