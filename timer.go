package timer

import (
	"fmt"

	"github.com/jy01095902/timer/logger"
)

type Timer interface {
	Run(fns ...func() error)
}

type TimedFunc func() error

func (fn TimedFunc) hideError(timerName string) func() {
	return func() {
		err := fn()
		if err != nil {
			logger.Error(fmt.Sprintf("an error occurred when %s executing function", timerName), "error", err.Error())
		}
	}
}
