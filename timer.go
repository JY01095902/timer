package timer

import (
	"sort"
	"time"

	"github.com/robfig/cron/v3"
)

type TimedTask struct {
	Id          int
	EntryId     cron.EntryID
	Prev        time.Time
	Next        time.Time
	Name        string
	Description string
}

type Timer struct {
	cronTimer *cron.Cron
	tasks     map[int]TimedTask
}

func NewTimer() Timer {
	timer := Timer{
		cronTimer: cron.New(),
		tasks:     map[int]TimedTask{},
	}

	return timer
}

func (timer Timer) AddDisposableTask(schedule Schedule, job cron.Job, name, description string) int {
	disJob, callback := NewCallbackJob(job)
	id := timer.AddTask(schedule, disJob, name, description)
	callback(func() {
		task := timer.GetTask(id)
		task.Next = time.Time{}
		timer.tasks[id] = task
		timer.cronTimer.Remove(task.EntryId)
	})

	return id
}

func (timer Timer) AddTask(schedule Schedule, job cron.Job, name, description string) int {
	if schedule.DoFirst {
		job.Run()
	}

	entryId := timer.cronTimer.Schedule(schedule, job)

	return timer.addTask(name, description, entryId)
}

func (timer Timer) getNextId() int {
	tasks := timer.GetTasks()
	if len(tasks) == 0 {
		return 1
	}

	sort.Sort(byId(tasks))

	return tasks[len(tasks)-1].Id + 1
}

func (timer Timer) addTask(name, description string, entryId cron.EntryID) int {
	task := TimedTask{
		Id:          timer.getNextId(),
		EntryId:     entryId,
		Name:        name,
		Description: description,
	}
	timer.refreshTask(task.Id)
	timer.tasks[task.Id] = task

	return task.Id
}

func (timer Timer) refreshTask(id int) {
	task, exist := timer.tasks[id]
	if !exist {
		return
	}

	entry := timer.cronTimer.Entry(task.EntryId)
	if int(entry.ID) == 0 {
		return
	}

	task.Prev = entry.Prev
	task.Next = entry.Next
	if entry.Next.IsZero() && entry.Schedule != nil {
		task.Next = entry.Schedule.Next(time.Now())
	}
	timer.tasks[id] = task
}

func (timer Timer) GetTask(id int) TimedTask {
	timer.refreshTask(id)

	return timer.tasks[id]
}

func (timer Timer) GetTasks() []TimedTask {
	tasks := []TimedTask{}
	for _, t := range timer.tasks {
		tasks = append(tasks, timer.GetTask(t.Id))
	}

	sort.Sort(byId(tasks))

	return tasks
}

func (timer *Timer) Remove(id int) {
	delete(timer.tasks, id)

	if task, exist := timer.tasks[id]; exist {
		timer.cronTimer.Remove(task.EntryId)
	}
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

type byId []TimedTask

func (s byId) Len() int      { return len(s) }
func (s byId) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byId) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}
