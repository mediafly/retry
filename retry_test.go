package retry

import (
	"errors"
	"testing"
)

func TestRetryCompleted(t *testing.T) {

}

func TestRetryFailed(t *testing.T) {

}

func TestRetryCancelled(t *testing.T) {
	do := make(chan bool, 1)
	cancelled := make(chan interface{})
	done := make(chan interface{})

	o := Options{
		Do: func() (Result, error) {
			do <- true
			return Continue, errors.New("error")
		},
		Cancel: func() error {
			close(cancelled)
			return nil
		},
	}

	r := New(o)

	go func() {
		if err := r.Do(); err != nil && err != Cancelled {
			t.Fatal(err)
		}
		close(done)
	}()

	// cancel after first failure that way we know do was called once and we
	// can see if cancel is in fact short circuiting the delay
	<-do

	if err := r.Cancel(); err != nil && err != Cancelled {
		t.Fatal(err)
	}

	<-done
	<-cancelled
}
