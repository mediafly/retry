package retry

import (
	"log"
	"os"
)

/////////////////////////////////////////////////////////////////////////////
// Logger
/////////////////////////////////////////////////////////////////////////////
type Logger interface {
	IsDebugEnabled() bool
	Debug(...interface{})
	IsErrorEnabled() bool
	Error(...interface{})
}

/////////////////////////////////////////////////////////////////////////////
// Default Logger
/////////////////////////////////////////////////////////////////////////////
var defaultLogger Logger = &StandardLogger{log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)}

func SetDefaultLogger(logger Logger) {
	if logger == nil {
		panic("logger is nil")
	}

	defaultLogger = logger
}

/////////////////////////////////////////////////////////////////////////////
// StandardLogger
/////////////////////////////////////////////////////////////////////////////
type StandardLogger struct {
	*log.Logger
}

func (l *StandardLogger) IsDebugEnabled() bool {
	return false
}

func (l *StandardLogger) Debug(v ...interface{}) {
}

func (l *StandardLogger) IsErrorEnabled() bool {
	return true
}

func (l *StandardLogger) Error(v ...interface{}) {
	l.Println(v...)
}

/////////////////////////////////////////////////////////////////////////////
// NullLogger
/////////////////////////////////////////////////////////////////////////////
type NullLogger struct {
}

func (l *NullLogger) IsDebugEnabled() bool {
	return false
}

func (l *NullLogger) Debug(v ...interface{}) {
}

func (l *NullLogger) IsErrorEnabled() bool {
	return false
}
func (l *NullLogger) Error(v ...interface{}) {
}
