package timer

import (
	"time"

	"github.com/robfig/cron/v3"
)

type TimedTask struct {
	Id          cron.EntryID
	Prev        time.Time
	Next        time.Time
	Name        string
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

func (timer Timer) AddDisposableTask(schedule Schedule, job cron.Job, name, description string) cron.EntryID {
	disJob, callback := NewCallbackJob(job)
	entryId := timer.AddTask(schedule, disJob, name, description)
	callback(func() {
		timer.cronTimer.Remove(entryId)
	})

	return entryId
}

func (timer Timer) AddTask(schedule Schedule, job cron.Job, name, description string) cron.EntryID {
	if schedule.DoFirst {
		job.Run()
	}

	entryId := timer.cronTimer.Schedule(schedule, job)
	task := TimedTask{
		Id:          entryId,
		Name:        name,
		Description: description,
	}
	timer.tasks[entryId] = task

	return entryId
}

func (timer Timer) GetTask(id cron.EntryID) TimedTask {
	entry := timer.cronTimer.Entry(id)
	task := timer.tasks[id]
	task.Prev = entry.Prev
	task.Next = entry.Next
	if entry.Next.IsZero() {
		task.Next = entry.Schedule.Next(time.Now())
	}

	return task
}

func (timer Timer) GetTasks() []TimedTask {
	tasks := []TimedTask{}
	for _, t := range timer.tasks {
		tasks = append(tasks, timer.GetTask(t.Id))
	}

	return tasks
}

func (timer *Timer) Remove(id cron.EntryID) {
	delete(timer.tasks, id)
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
