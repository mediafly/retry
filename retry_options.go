package retry

import (
	"time"
)

type Options struct {
	Tag         interface{}
	Deadline    time.Time
	MaxAttempts uint32
	MaxDelay    time.Duration
	Do          Callback
	Events      chan<- interface{}
	Cancel      func() error
}
