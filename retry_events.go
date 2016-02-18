package retry

/////////////////////////////////////////////////////////////////////////////
// RetryEvent
/////////////////////////////////////////////////////////////////////////////
type RetryEvent interface {
	retryEvent()
}

func (e *RetryStartedEvent) retryEvent()   {}
func (e *RetryCompletedEvent) retryEvent() {}
func (e *RetryFailedEvent) retryEvent()    {}
func (e *RetryDoFailedEvent) retryEvent()  {}

/////////////////////////////////////////////////////////////////////////////
// RetryStartedEvent
/////////////////////////////////////////////////////////////////////////////
type RetryStartedEvent struct {
	Retry Retry
	Tag   interface{}
}

/////////////////////////////////////////////////////////////////////////////
// RetryCompletedEvent
/////////////////////////////////////////////////////////////////////////////
type RetryCompletedEvent struct {
	Retry Retry
	Tag   interface{}
}

/////////////////////////////////////////////////////////////////////////////
// RetryFailedEvent
/////////////////////////////////////////////////////////////////////////////
type RetryFailedEvent struct {
	Retry Retry
	Tag   interface{}
	Error error
}

/////////////////////////////////////////////////////////////////////////////
// RetryDoFailedEvent
/////////////////////////////////////////////////////////////////////////////
type RetryDoFailedEvent struct {
	Retry Retry
	Tag   interface{}
	Error error
}
