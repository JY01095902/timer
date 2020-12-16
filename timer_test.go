package timer

import (
	"fmt"
	"testing"
	"time"

	"github.com/robfig/cron/v3"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAddDisposableTask(t *testing.T) {
	timer := NewTimer()
	dur := 2 * time.Second
	now := time.Now()
	id := timer.AddDisposableTask(NewSchedule(cron.Every(dur)), cron.FuncJob(func() {}), "", "")
	timer.Start()

	Convey("刷新添加一次性任务", t, func() {
		Convey(fmt.Sprintf("新添加的任务，上次执行时间应该是：%s，下次执行时间应该是：%s", time.Time{}.Format(time.RFC3339), now.Add(dur).Format(time.RFC3339)), func() {
			task := timer.GetTask(id)
			So(task.Prev.IsZero(), ShouldBeTrue)
			So(task.Next.Format(time.RFC3339), ShouldEqual, now.Add(dur).Format(time.RFC3339))
		})

		tr := time.NewTimer(dur)
		<-tr.C

		Convey(fmt.Sprintf("执行后的任务，上次执行时间应该是：%s，下次执行时间应该是：%s", now.Add(dur).Format(time.RFC3339), time.Time{}.Format(time.RFC3339)), func() {
			task := timer.GetTask(id)
			So(task.Prev.Format(time.RFC3339), ShouldEqual, now.Add(dur).Format(time.RFC3339))
			So(task.Next.IsZero(), ShouldBeTrue)
		})
	})
}

func TestRefreshTask(t *testing.T) {
	timer := NewTimer()
	dur := 2 * time.Second
	now := time.Now()
	cbJob, cb := NewCallbackJob(cron.FuncJob(func() {}))
	id := timer.AddTask(NewSchedule(cron.Every(dur)), cbJob, "", "")
	cb(func() {
		Convey("刷新定时任务信息", t, func() {
			Convey(fmt.Sprintf("执行后的任务，上次执行时间应该是：%s，下次执行时间应该是：%s", now.Add(dur).Format(time.RFC3339), now.Add(dur).Add(dur).Format(time.RFC3339)), func() {
				task := timer.GetTask(id)
				So(task.Prev.Format(time.RFC3339), ShouldEqual, now.Add(dur).Format(time.RFC3339))
				So(task.Next.Format(time.RFC3339), ShouldEqual, now.Add(dur).Add(dur).Format(time.RFC3339))
			})
		})
	})
	timer.Start()

	Convey("刷新定时任务信息", t, func() {
		Convey(fmt.Sprintf("新添加的任务，上次执行时间应该是：%s，下次执行时间应该是：%s", time.Time{}.Format(time.RFC3339), now.Add(dur).Format(time.RFC3339)), func() {
			task := timer.GetTask(id)
			So(task.Prev.IsZero(), ShouldBeTrue)
			So(task.Next.Format(time.RFC3339), ShouldEqual, now.Add(dur).Format(time.RFC3339))
		})
	})
	tr := time.NewTimer(dur)
	<-tr.C
}
