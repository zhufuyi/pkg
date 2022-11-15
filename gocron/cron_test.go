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
		time.Sleep(time.Second)
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
			TimeSpec: "@every 1s",
			Fn:       task1,
		},
		{
			Name:     "task2",
			TimeSpec: "@every 2s",
			Fn:       task2,
		},
		{
			Name:     "task3",
			TimeSpec: "@every 3s",
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

	time.Sleep(time.Second * 7)
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

	time.Sleep(time.Second * 5)
}

func TestEvery(t *testing.T) {
	task_1 := func() {
		fmt.Println("this is task_1")
		fmt.Println("running task list:", GetRunningTasks()) // 当前运行的任务
	}
	task_2 := func() { // 如果执行时间超过定时时间，不会影响下一个新定时任务的执行
		fmt.Println("this is task_2")
	}
	task_3 := func() {
		fmt.Println("this is task_3")
	}

	tasks := []*Task{
		{
			TimeSpec: EverySecond(5),
			Name:     "task_1",
			Fn:       task_1,
		},
		{
			TimeSpec: EveryMinute(1),
			Name:     "task_2",
			Fn:       task_2,
		},
		{
			TimeSpec: EveryHour(1),
			Name:     "task_3",
			Fn:       task_3,
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

	time.Sleep(time.Second * 7)
}
