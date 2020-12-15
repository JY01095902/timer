package timer

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(args ...interface{})
}

func NewLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.CallerKey = ""
	config.EncoderConfig.StacktraceKey = ""
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	l, _ := config.Build()

	return l.Sugar()
}

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
	logger    Logger
}

func NewTimer() Timer {

	timer := Timer{
		cronTimer: cron.New(),
		tasks:     map[cron.EntryID]TimedTask{},
		logger:    NewLogger(),
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
		if !schedule.IsSilent {
			timer.logger.Info(fmt.Sprintf("-=%s(%s)=- execution finished, next execution time: %s", name, description, schedule.Next(time.Now()).Format(time.RFC3339)))
		}
	}

	cbJob, callback := NewCallbackJob(job)
	entryId := timer.cronTimer.Schedule(schedule, cbJob)
	task := TimedTask{
		Id:          entryId,
		Name:        name,
		Description: description,
	}
	timer.tasks[entryId] = task

	if !schedule.IsSilent {
		callback(func() {
			task := timer.GetTask(entryId)
			timer.logger.Info(fmt.Sprintf("-=%s(%s)=- execution finished, next execution time: %s", task.Name, task.Description, task.Next.Format(time.RFC3339)))
		})
	}

	return entryId
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
