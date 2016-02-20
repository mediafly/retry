package retry

/////////////////////////////////////////////////////////////////////////////
// CancelledError
/////////////////////////////////////////////////////////////////////////////
type CancelledError struct {
	Retry  Retry
	Source error
}

func (e *CancelledError) Error() string {
	return "cancelled"
}

/////////////////////////////////////////////////////////////////////////////
// DeadlineError
/////////////////////////////////////////////////////////////////////////////
type DeadlineError struct {
	Retry Retry
}

func (e *DeadlineError) Error() string {
	return "deadline"
}
