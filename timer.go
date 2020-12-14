package timer

import (
	"fmt"
	"time"

	"github.com/jy01095902/timer/logger"
	"github.com/robfig/cron/v3"
)

// type Timer interface {
// 	Run(fns ...func() error)
// }

type TimedFunc func() error

func (fn TimedFunc) hideError(timerName string) func() {
	return func() {
		err := fn()
		if err != nil {
			logger.Error(fmt.Sprintf("an error occurred when %s executing function", timerName), "error", err.Error())
		}
	}
}

type TimedTask struct {
	Id          cron.EntryID
	Prev        time.Time
	Next        time.Time
	Description string
}

type Timer struct {
	cronTimer *cron.Cron
	tasks     map[cron.EntryID]TimedTask
}

func NewTimer() Timer {
	timer := Timer{
		cronTimer: cron.New(),
		tasks:     map[cron.EntryID]TimedTask{},
	}

	return timer
}

func (timer Timer) AddTask(schedule cron.Schedule, job func(), description string) {
	entryId := timer.cronTimer.Schedule(schedule, cron.FuncJob(job))
	task := TimedTask{
		Id:          entryId,
		Description: description,
	}
	timer.tasks[entryId] = task
}

func (timer Timer) GetTask(id cron.EntryID) TimedTask {
	entry := timer.cronTimer.Entry(id)
	task := timer.tasks[id]
	task.Prev = entry.Prev
	task.Next = entry.Next

	return task
}

func (timer Timer) GetTasks() []TimedTask {
	tasks := []TimedTask{}
	for _, t := range timer.tasks {
		tasks = append(tasks, timer.GetTask(t.Id))
	}

	return tasks
}

func (timer Timer) Remove(id cron.EntryID) {
	timer.cronTimer.Remove(id)
}

func (timer Timer) Start() {
	timer.cronTimer.Start()
}

func (timer Timer) Run() {
	timer.cronTimer.Run()
}

func (timer Timer) Stop() {
	timer.cronTimer.Stop()
}
