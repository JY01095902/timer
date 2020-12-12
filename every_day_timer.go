package timer

import (
	"fmt"

	"github.com/jy01095902/timer/logger"
	"github.com/robfig/cron/v3"
)

type EveryDayTimer struct {
	hour     int
	min      int
	DoAtOnce bool
}

func NewEveryDayTimer(hour, min int) EveryDayTimer {
	timer := EveryDayTimer{
		hour: hour,
		min:  min,
	}

	return timer
}

func (timer EveryDayTimer) Run(fns ...TimedFunc) {
	if timer.hour < 0 || timer.hour > 24 {
		return
	}

	if timer.min < 0 || timer.min > 59 {
		return
	}

	if len(fns) == 0 {
		return
	}

	name := "-=every day timer=-"
	c := cron.New()
	for _, fn := range fns {
		if timer.DoAtOnce {
			fn.hideError(name)()
		}
		spec := fmt.Sprintf("%d %d * * ?", timer.min, timer.hour)
		entryId, err := c.AddFunc(spec, fn.hideError(name))
		if err != nil {
			logger.Error(fmt.Sprintf("an error occurred when %s executing function", name), "error", err.Error())
		}

		logger.Info(fmt.Sprintf("%s tasks has been created", name), "entry id", entryId)
	}
	c.Start()
	logger.Info(fmt.Sprintf("%s started", name))
}
