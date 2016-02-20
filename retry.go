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
	options     Options
	doError     error
	cancelError error
	cancelled   chan interface{}
	done        chan interface{}
	once        sync.Once
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
		if r.doError != nil {
			if log.IsErrorEnabled() {
				log.Error(r, "failed:", r.doError)
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

	if r.options.MaxDelay >= 0 {
		maxDelay = r.options.MaxDelay
	}

	deadline := time.Now().UTC().Add(time.Duration(int64(math.MaxInt64)))

	if r.options.Deadline.After(time.Now().UTC()) {
		deadline = r.options.Deadline
	}

	for attempt := uint32(0); attempt < maxAttempts; attempt++ {
		if r.doError != nil {
			if log.IsErrorEnabled() {
				log.Error(r, "attempt", attempt, "/", maxAttempts, "failed, retrying:", r.doError)
			}
		}

		delay := minDuration(maxDelay, time.Second*time.Duration(attempt*attempt))

		// TODO: allow Cancel to run concurrently with Do
		select {
		case <-r.cancelled:
			if log.IsDebugEnabled() {
				log.Debug(r, "cancelled")
			}

			if r.options.Cancel != nil {
				if err := r.options.Cancel(); err != nil {
					r.cancelError = err
				}
			}

			if r.doError == nil {
				r.doError = &CancelledError{r, r.doError}
			}

			return r.doError

		case <-time.After(deadline.Sub(time.Now().UTC())):
			if log.IsDebugEnabled() {
				log.Debug(r, "passed deadline")
			}

			if r.doError == nil {
				r.doError = &DeadlineError{r}
			}

			return r.doError

		case <-time.After(delay):
		}

		if log.IsDebugEnabled() {
			log.Debug(r, "do")
		}

		cont, err := r.options.Do()

		r.doError = err

		if r.doError == nil || cont == Stop {
			return r.doError
		}
	}

	return r.doError
}

func (r *retry) Cancel() error {
	r.once.Do(func() {
		close(r.cancelled)
	})

	<-r.done

	return r.cancelError
}
