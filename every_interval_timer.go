package timer

import (
	"fmt"
	"time"

	"github.com/jy01095902/timer/logger"
	"github.com/robfig/cron/v3"
)

type EveryIntervalTimer struct {
	invlDur  time.Duration
	DoAtOnce bool
}

// 间隔以秒为单位
func NewEveryIntervalTimer(sec int) EveryIntervalTimer {
	d, _ := time.ParseDuration(fmt.Sprintf("%ds", sec))
	timer := EveryIntervalTimer{
		invlDur: d,
	}

	return timer
}

func (timer EveryIntervalTimer) Run(fns ...TimedFunc) {
	if len(fns) == 0 {
		return
	}

	name := "-=every interval timer=-"
	c := cron.New()
	for _, fn := range fns {
		if timer.DoAtOnce {
			fn.hideError(name)()
		}
		spec := "@every " + timer.invlDur.String()
		entryId, err := c.AddFunc(spec, fn.hideError(name))
		if err != nil {
			logger.Error(fmt.Sprintf("an error occurred when %s executing function", name), "error", err.Error())
		}

		logger.Info(fmt.Sprintf("%s tasks has been created", name), "entry id", entryId)
	}
	c.Start()
	logger.Info(fmt.Sprintf("%s started", name))
}
