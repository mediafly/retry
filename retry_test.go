package retry

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestContinue(t *testing.T) {
	err := errors.New("TestContinue")
	invoked := 0
	options := Options{
		MaxDelay: 0,
		Do: func() Result {
			invoked++
			if invoked == 2 {
				return Stop()
			} else {
				return Continue(err)
			}
		},
	}
	assert.NoError(t, Do(options))
	assert.Equal(t, 2, invoked)
}

func TestStopWithNoErrorStops(t *testing.T) {
	invoked := 0
	options := Options{
		Do: func() Result {
			invoked++
			return Stop()
		},
	}
	assert.NoError(t, Do(options))
	assert.Equal(t, 1, invoked)
}

func TestStopWithErrorStops(t *testing.T) {
	err := errors.New("TestStopWithErrorStops")
	invoked := 0
	options := Options{
		Do: func() Result {
			invoked++
			return Stop(err)
		},
	}
	assert.Equal(t, err, Do(options))
	assert.Equal(t, 1, invoked)
}

func TestMaxAttempts(t *testing.T) {
	err := errors.New("TestMaxAttempts")
	invoked := 0
	options := Options{
		MaxAttempts: 10,
		MaxDelay:    0,
		Do: func() Result {
			invoked++
			return Continue(err)
		},
	}
	assert.Equal(t, err, Do(options))
	assert.Equal(t, 10, invoked)
}

func TestDeadline(t *testing.T) {
	err := errors.New("TestDeadline")
	invoked := 0
	options := Options{
		Deadline: time.Now().UTC().Add(time.Millisecond),
		MaxDelay: 0,
		Do: func() Result {
			invoked++
			return Continue(err)
		},
	}
	assert.Equal(t, err, Do(options).(*DeadlineError).Err)
	assert.True(t, invoked > 0)
}

func TestMaxDelay(t *testing.T) {
	err := errors.New("TestMaxDelay")
	invoked := 0
	options := Options{
		MaxDelay: time.Millisecond,
		Deadline: time.Now().UTC().Add(time.Millisecond * 10),
		Do: func() Result {
			invoked++
			if invoked == 3 {
				return Stop(err)
			} else {
				return Continue(err)
			}
		},
	}
	assert.Equal(t, err, Do(options))
	assert.Equal(t, 3, invoked)
}

func TestCancelled(t *testing.T) {
	err := errors.New("TestCancelled")
	ctx, cancel := context.WithCancel(context.Background())
	invoked := 0
	options := Options{
		MaxDelay: 0,
		Context:  ctx,
		Do: func() Result {
			invoked++
			return Continue(err)
		},
	}

	go func() {
		<-time.After(time.Microsecond)
		cancel()
	}()

	assert.Equal(t, err, Do(options).(*CancelledError).Err)
	assert.True(t, invoked > 0)
}
