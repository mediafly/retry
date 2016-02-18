package retry

import (
	"fmt"
)

type Result int

const (
	Stop     Result = iota
	Continue Result = iota
)

func (r Result) String() string {
	switch r {
	case Continue:
		return "continue"
	case Stop:
		return "stop"
	default:
		panic(fmt.Sprint("invalid result:", r))
	}
}
