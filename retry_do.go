package retry

import (
	"time"
)

// Perform a Retry operation based on Options.
func Do(options Options) error {
	return New(options).Do()
}

// Perform a Retry func() until completed successfully.
func DoFunc(callback func() Result) error {
	return New(Options{Do: callback}).Do()
}

// Perform a retry func() until completed successfully or deadline passed.
func DoDeadline(deadline time.Time, callback func() Result) error {
	return New(Options{Do: callback, Deadline: deadline}).Do()
}

// Perform a retry func() until completed successfully or max attempts passed.
func DoMaxAttempts(maxAttempts uint32, callback func() Result) error {
	return New(Options{Do: callback, MaxAttempts: maxAttempts}).Do()
}
