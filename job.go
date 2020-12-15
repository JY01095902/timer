package timer

import (
	"github.com/robfig/cron/v3"
)

type CallbackJob struct {
	job      cron.Job
	callback func()
}

func NewCallbackJob(job cron.Job) (*CallbackJob, func(callback func())) {
	cbJob := CallbackJob{
		job:      job,
		callback: func() {},
	}

	return &cbJob, func(callback func()) {
		cbJob.callback = callback
	}
}

func (j *CallbackJob) Run() {
	j.job.Run()

	if j.callback != nil {
		j.callback()
	}
}
