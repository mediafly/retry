/*
Retry is a utility package to to facilitate operations that need to be retried
when they fail.
*/
package retry

import (
	"fmt"
	"github.com/mediafly/math"
	"github.com/mediafly/math/constants"
	"golang.org/x/net/context"
	"time"
)

// Retry represents an operation that needs to be retried until is finishes
// successfully.
type Retry interface {
	Do() error
}

type retry struct {
	Options
}

// Returns a new Retry operation based on Options.
func New(options Options) Retry {
	return &retry{options}
}

func (r *retry) String() string {
	return fmt.Sprintf("[%T Tag=%v]", r, r.Tag)
}

// Perform the operation until it completes successfully or retry thresholds
// are crossed.
func (r *retry) Do() error {
	debugLog := r.DebugLog

	if debugLog == nil {
		debugLog = defaultDebugLogger
	}

	errorLog := r.ErrorLog

	if errorLog == nil {
		errorLog = defaultErrorLogger
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
			if debugLog != nil {
				debugLog(r, "completed")
			}

			return result.Err
		}

		{
			result := result.(*continueResult)

			if attempt >= maxAttempts {
				if errorLog != nil {
					errorLog(r, "failed, max attempts:", result.Err)
				}
				return result.Err
			}

			if attempt == 0 {
				if debugLog != nil {
					debugLog(r, "started")
				}
			} else {
				if errorLog != nil {
					errorLog(r, "retrying attempt", attempt+1, "of", maxAttempts, ":", result.Err)
				}
			}

			delay := math.MinDuration(maxDelay, time.Second*time.Duration(attempt*attempt))

			select {
			case <-ctx.Done():
				if errorLog != nil {
					errorLog(r, "cancelled:", result.Err)
				}
				return &CancelledError{r, result.Err}

			case <-time.After(deadline.Sub(time.Now().UTC())):
				if errorLog != nil {
					errorLog(r, "passed deadline:", result.Err)
				}
				return &DeadlineError{r, result.Err}

			case <-time.After(delay):
			}
		}

		result = r.Options.Do()

		attempt++
	}
}
