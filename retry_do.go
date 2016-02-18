package retry

import (
	"time"
)

/////////////////////////////////////////////////////////////////////////////
// Do
/////////////////////////////////////////////////////////////////////////////
func Do(options Options) error {
	return New(options).Do()
}

/////////////////////////////////////////////////////////////////////////////
// DoFunc
/////////////////////////////////////////////////////////////////////////////
func DoFunc(callback DoCallback) error {
	return New(Options{Do: callback}).Do()
}

/////////////////////////////////////////////////////////////////////////////
// DoDeadline
/////////////////////////////////////////////////////////////////////////////
func DoDeadline(deadline time.Time, callback DoCallback) error {
	return New(Options{Do: callback, Deadline: deadline}).Do()
}

/////////////////////////////////////////////////////////////////////////////
// DoMaxAttempts
/////////////////////////////////////////////////////////////////////////////
func DoMaxAttempts(maxAttempts uint32, callback DoCallback) error {
	return New(Options{Do: callback, MaxAttempts: maxAttempts}).Do()
}
