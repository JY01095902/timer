package timer

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Schedule struct {
	schedule cron.Schedule
	DoFirst  bool
}

func NewSchedule(schedule cron.Schedule) Schedule {
	return Schedule{
		schedule: schedule,
	}
}

func (s Schedule) Next(t time.Time) time.Time {
	return s.schedule.Next(t)
}
