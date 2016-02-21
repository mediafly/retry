package retry

import (
	"fmt"
	"github.com/mediafly/math"
	"github.com/mediafly/math/constants"
	"golang.org/x/net/context"
	"time"
)

/////////////////////////////////////////////////////////////////////////////
// Retry
/////////////////////////////////////////////////////////////////////////////
type Retry interface {
	Do() error
}

/////////////////////////////////////////////////////////////////////////////
// retry
/////////////////////////////////////////////////////////////////////////////
type retry struct {
	Options
}

func New(options Options) Retry {
	return &retry{options}
}

func (r *retry) String() string {
	return fmt.Sprintf("[%T Tag=%v]", r, r.Tag)
}

func (r *retry) Do() error {
	log := r.Log

	if log == nil {
		log = defaultLogger
	}

	maxAttempts := constants.MaxUint32

	if r.MaxAttempts > 0 {
		maxAttempts = r.MaxAttempts
	}

	maxDelay := time.Second * 30

	if r.MaxDelay >= 0 {
		maxDelay = r.MaxDelay
	}

	deadline := constants.MaxTime

	if r.Deadline.After(time.Now().UTC()) {
		deadline = r.Deadline
	}

	ctx := r.Context

	if ctx == nil {
		ctx = context.Background()
	}

	var attempt uint32 = 0
	var result Result = &continueResult{nil}

	for {
		if result, ok := result.(*stopResult); ok {
			if log.IsDebugEnabled() {
				log.Debug(r, "completed")
			}

			return result.Err
		}

		{
			result := result.(*continueResult)

			if attempt >= maxAttempts {
				if log.IsErrorEnabled() {
					log.Error(r, "failed, max attempts:", result.Err)
				}
				return result.Err
			}

			if attempt == 0 {
				if log.IsDebugEnabled() {
					log.Debug(r, "started")
				}
			} else {
				if log.IsErrorEnabled() {
					log.Error(r, "retrying attempt", attempt+1, "of", maxAttempts, ":", result.Err)
				}
			}

			delay := math.MinDuration(maxDelay, time.Second*time.Duration(attempt*attempt))

			select {
			case <-ctx.Done():
				if log.IsErrorEnabled() {
					log.Error(r, "cancelled:", result.Err)
				}
				return &CancelledError{r, result.Err}

			case <-time.After(deadline.Sub(time.Now().UTC())):
				if log.IsErrorEnabled() {
					log.Error(r, "passed deadline:", result.Err)
				}
				return &DeadlineError{r, result.Err}

			case <-time.After(delay):
			}
		}

		result = r.Options.Do()

		attempt++
	}
}
