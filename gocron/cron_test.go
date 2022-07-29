package gocron

import (
	"fmt"
	"testing"
	"time"
)

func TestInitAndRun(t *testing.T) {
	count := 0
	task1 := func() {
		fmt.Println("running task list:", GetRunningTasks()) // 当前运行的任务
	}
	task2 := func() { // 如果执行时间超过定时时间，不会影响下一个新定时任务的执行
		time.Sleep(time.Second * 12)
	}
	task3 := func() {
		count++
		if count%3 == 0 {
			panic("触发panic")
		}
	}

	tasks := []*Task{
		{
			Name:     "task1",
			TimeSpec: "@every 5s",
			Fn:       task1,
		},
		{
			Name:     "task2",
			TimeSpec: "@every 10s",
			Fn:       task2,
		},
		{
			Name:     "task3",
			TimeSpec: "@every 15s",
			Fn:       task3,
		},
	}

	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	err = Run(tasks...)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Minute)
}

func TestRunOnce(t *testing.T) {
	myTask := func() {
		taskName := "myTask"
		fmt.Println("running task list:", GetRunningTasks()) // 当前运行的任务
		fmt.Printf("the task '%s' is executed only once\n", taskName)
		DeleteTask(taskName)                                 // 执行完删除任务
		fmt.Println("running task list:", GetRunningTasks()) // 当前运行的任务
	}

	tasks := []*Task{
		{
			Name:     "myTask",
			Fn:       myTask,
			TimeSpec: "@every 2s",
		},
	}

	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	err = Run(tasks...)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 10)
}

func TestEvery(t *testing.T) {
	task1 := func() {
		fmt.Println("this is task1")
		fmt.Println("running task list:", GetRunningTasks()) // 当前运行的任务
	}
	task2 := func() { // 如果执行时间超过定时时间，不会影响下一个新定时任务的执行
		fmt.Println("this is task2")
	}
	task3 := func() {
		fmt.Println("this is task3")
	}

	tasks := []*Task{
		{
			TimeSpec: EverySecond(5),
			Name:     "task1",
			Fn:       task1,
		},
		{
			TimeSpec: EveryMinute(1),
			Name:     "task2",
			Fn:       task2,
		},
		{
			TimeSpec: EveryHour(1),
			Name:     "task3",
			Fn:       task3,
		},
	}

	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	err = Run(tasks...)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Minute)
}
