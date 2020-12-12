package timer

import (
	"fmt"
	"time"

	"github.com/jy01095902/timer/logger"
	"github.com/robfig/cron/v3"
)

type SpecifiedDateTimeTimer struct {
	dateTime time.Time
}

func NewSpecifiedDateTimeTimer(dateTime time.Time) SpecifiedDateTimeTimer {
	timer := SpecifiedDateTimeTimer{
		dateTime: dateTime,
	}

	return timer
}

func (timer SpecifiedDateTimeTimer) Run(fns ...TimedFunc) {
	if timer.dateTime.Before(time.Now()) {
		return
	}

	if len(fns) == 0 {
		return
	}

	name := "-=specified date time timer=-"
	c := cron.New()
	for _, fn := range fns {

		spec := "@every " + time.Until(timer.dateTime).String()
		entryId, err := c.AddFunc(spec, fn.hideError(name))
		if err != nil {
			logger.Error(fmt.Sprintf("an error occurred when %s executing function", name), "error", err.Error())
		}

		logger.Info(fmt.Sprintf("%s tasks has been created", name), "entry id", entryId)
	}

	// 此定时任务只执行一次，所以执行完后停止
	go func() {
		time.AfterFunc(time.Until(timer.dateTime), func() {
			c.Stop()
			logger.Info(fmt.Sprintf("%s stopped", name))
		})
	}()

	c.Start()
	logger.Info(fmt.Sprintf("%s started", name))
}
