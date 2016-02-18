package retry

import (
	"errors"
	"math"
	"sync"
	"time"
)

/////////////////////////////////////////////////////////////////////////////
// Errors
/////////////////////////////////////////////////////////////////////////////
var Cancelled error = errors.New("cancelled")

/////////////////////////////////////////////////////////////////////////////
// Callback
/////////////////////////////////////////////////////////////////////////////
type Callback func() (Result, error)

/////////////////////////////////////////////////////////////////////////////
// Retry
/////////////////////////////////////////////////////////////////////////////
type Retry interface {
	Do() error
	Cancel() error
}

/////////////////////////////////////////////////////////////////////////////
// Do
/////////////////////////////////////////////////////////////////////////////
func Do(callback Callback) error {
	return New(Options{Do: callback}).Do()
}

/////////////////////////////////////////////////////////////////////////////
// DoUntil
/////////////////////////////////////////////////////////////////////////////
func DoUntil(deadline time.Time, callback Callback) error {
	return New(Options{Do: callback, Deadline: deadline}).Do()
}

/////////////////////////////////////////////////////////////////////////////
// retry
/////////////////////////////////////////////////////////////////////////////
type retry struct {
	options   Options
	err       error
	cancelled chan interface{}
	done      chan interface{}
	events    chan<- RetryEvent
	once      sync.Once
}

func New(options Options) Retry {
	return &retry{
		options:   options,
		cancelled: make(chan interface{}),
		done:      make(chan interface{}),
	}
}

func (r *retry) Do() error {
	defer func() {
		if r.err != nil {
			r.raise(&RetryFailedEvent{r, r.options.Tag, r.err})
		} else {
			r.raise(&RetryCompletedEvent{r, r.options.Tag})
		}

		if r.events != nil {
			close(r.events)
		}

		close(r.done)
	}()

	r.raise(&RetryStartedEvent{r, r.options.Tag})

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
		delay := minDuration(maxDelay, time.Second*time.Duration(attempt*attempt))

		// TODO: allow Cancel to run concurrently with Do
		select {
		case <-r.cancelled:
			r.err = Cancelled
			if r.options.Cancel != nil {
				if err := r.options.Cancel(); err != nil {
					r.err = err
				}
			}
			return r.err

		case <-time.After(deadline.Sub(time.Now().UTC())):
			if r.err == nil {
				r.err = errors.New("passed deadline")
			}
			return r.err

		case <-time.After(delay):
		}

		cont, err := r.options.Do()

		r.err = err

		if r.err != nil {
			r.raise(&RetryDoFailedEvent{r, r.options.Tag, r.err})
		}

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

func (r *retry) raise(event RetryEvent) {
	if r.options.Events != nil {
		r.options.Events <- event
	}
}
