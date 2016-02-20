package retry

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestContinueWithNoErrorStops(t *testing.T) {
	invoked := 0
	options := Options{
		Do: func() (Result, error) {
			invoked++
			return Continue, nil
		},
		Log: &NullLogger{},
	}
	assert.NoError(t, Do(options))
	assert.Equal(t, 1, invoked)
}

func TestContinueWithErrorContinues(t *testing.T) {
	invoked := 0
	options := Options{
		MaxDelay: 0,
		Do: func() (Result, error) {
			invoked++
			if invoked == 2 {
				return Stop, nil
			} else {
				return Continue, errors.New("failed")
			}
		},
		Log: &NullLogger{},
	}
	assert.NoError(t, Do(options))
	assert.Equal(t, 2, invoked)
}

func TestStopWithNoErrorStops(t *testing.T) {
	invoked := 0
	options := Options{
		Do: func() (Result, error) {
			invoked++
			return Stop, nil
		},
		Log: &NullLogger{},
	}
	assert.NoError(t, Do(options))
	assert.Equal(t, 1, invoked)
}

func TestStopWithErrorStops(t *testing.T) {
	err := errors.New("error")
	invoked := 0
	options := Options{
		Do: func() (Result, error) {
			invoked++
			return Stop, err
		},
		Log: &NullLogger{},
	}
	assert.Equal(t, err, Do(options))
	assert.Equal(t, 1, invoked)
}

func TestMaxAttempts(t *testing.T) {
	err := errors.New("error")
	invoked := 0
	options := Options{
		MaxAttempts: 10,
		MaxDelay:    0,
		Do: func() (Result, error) {
			invoked++
			return Continue, err
		},
		Log: &NullLogger{},
	}
	assert.Equal(t, err, Do(options))
	assert.Equal(t, 10, invoked)
}

func TestDeadline(t *testing.T) {
	err := errors.New("error")
	invoked := 0
	options := Options{
		Deadline: time.Now().UTC().Add(time.Millisecond),
		MaxDelay: 0,
		Do: func() (Result, error) {
			invoked++
			return Continue, err
		},
		Log: &NullLogger{},
	}
	assert.Equal(t, err, Do(options))
	assert.True(t, invoked > 0)
}

func TestMaxDelay(t *testing.T) {
	err := errors.New("error")
	invoked := 0
	options := Options{
		MaxDelay: time.Millisecond,
		Deadline: time.Now().UTC().Add(time.Millisecond * 10),
		Do: func() (Result, error) {
			invoked++
			if invoked == 3 {
				return Stop, err
			} else {
				return Continue, err
			}
		},
		Log: &NullLogger{},
	}
	assert.Equal(t, err, Do(options))
	assert.Equal(t, 3, invoked)
}

func TestCancelledWithNoCancelError(t *testing.T) {
	err := errors.New("TestCancelledWithNoCancelError")
	cancelled := make(chan interface{})
	options := Options{
		Do: func() (Result, error) {
			<-cancelled
			return Stop, err
		},
		Cancel: func() error {
			close(cancelled)
			return nil
		},
		Log: &NullLogger{},
	}

	retry := New(options)

	go func() {
		assert.IsType(t, &CancelledError{}, retry.Do())
	}()

	assert.NoError(t, retry.Cancel())
}

func TestCancelledWithCancelError(t *testing.T) {
	err := errors.New("TestCancelledWithCancelError")
	cancelError := errors.New("Cancel Failed")
	cancelled := make(chan interface{})
	options := Options{
		Do: func() (Result, error) {
			<-cancelled
			return Stop, err
		},
		Cancel: func() error {
			close(cancelled)
			return cancelError
		},
		Log: &NullLogger{},
	}

	retry := New(options)

	go func() {
		assert.IsType(t, &CancelledError{}, retry.Do())
	}()

	assert.Equal(t, cancelError, retry.Cancel())
}
