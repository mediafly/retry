package retry

import (
	"golang.org/x/net/context"
	"time"
)

/////////////////////////////////////////////////////////////////////////////
// Options
/////////////////////////////////////////////////////////////////////////////
type Options struct {
	Do          func() Result
	Context     context.Context
	Tag         interface{}
	Deadline    time.Time
	MaxAttempts uint32
	MaxDelay    time.Duration
	Log         Logger
}
