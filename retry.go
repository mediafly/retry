package retry

import (
	"fmt"
	"math"
	"sync"
	"time"
)

/////////////////////////////////////////////////////////////////////////////
// Retry
/////////////////////////////////////////////////////////////////////////////
type Retry interface {
	Do() error
	Cancel() error
}

/////////////////////////////////////////////////////////////////////////////
// retry
/////////////////////////////////////////////////////////////////////////////
type retry struct {
	options   Options
	err       error
	cancelled chan interface{}
	done      chan interface{}
	once      sync.Once
}

func New(options Options) Retry {
	return &retry{
		options:   options,
		cancelled: make(chan interface{}),
		done:      make(chan interface{}),
	}
}

func (r *retry) String() string {
	return fmt.Sprintf("[%T Tag=%v]", r, r.options.Tag)
}

func (r *retry) Do() error {
	log := r.options.Log

	if log == nil {
		log = defaultLogger
	}

	defer func() {
		if r.err != nil {
			if log.IsErrorEnabled() {
				log.Error(r, "failed:", r.err)
			}
		} else {
			if log.IsDebugEnabled() {
				log.Debug(r, "completed")
			}
		}

		close(r.done)
	}()

	if log.IsDebugEnabled() {
		log.Debug(r, "started")
	}

	maxAttempts := uint32(math.MaxUint32)

	if r.options.MaxAttempts > 0 {
		maxAttempts = r.options.MaxAttempts
	}

	maxDelay := time.Second * 30

	if r.options.MaxDelay > 0 {
		maxDelay = r.options.MaxDelay
	}

	deadline := time.Now().UTC().Add(time.Duration(int64(math.MaxInt64)))

	if r.options.Deadline.After(time.Now().UTC()) {
		deadline = r.options.Deadline
	}

	for attempt := uint32(0); attempt < maxAttempts; attempt++ {
		if r.err != nil {
			if log.IsErrorEnabled() {
				log.Error(r, "attempt", attempt, "/", maxAttempts, "failed, retrying:", r.err)
			}
		}

		delay := minDuration(maxDelay, time.Second*time.Duration(attempt*attempt))

		// TODO: allow Cancel to run concurrently with Do
		select {
		case <-r.cancelled:
			if log.IsDebugEnabled() {
				log.Debug(r, "cancelled")
			}
			r.err = &CancelledError{r}
			if r.options.Cancel != nil {
				if err := r.options.Cancel(); err != nil {
					r.err = err
				}
			}
			return r.err

		case <-time.After(deadline.Sub(time.Now().UTC())):
			if log.IsDebugEnabled() {
				log.Debug(r, "passed deadline")
			}
			if r.err == nil {
				r.err = &DeadlineError{r}
			}
			return r.err

		case <-time.After(delay):
		}

		if log.IsDebugEnabled() {
			log.Debug(r, "do")
		}

		cont, err := r.options.Do()

		r.err = err

		if r.err == nil || cont == Stop {
			return r.err
		}
	}

	return r.err
}

func (r *retry) Cancel() error {
	r.once.Do(func() {
		close(r.cancelled)
	})

	<-r.done

	return r.err
}
