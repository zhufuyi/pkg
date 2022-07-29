## gocron

在[cron v3](github.com/robfig/cron)基础上封装的定时任务库。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/gocron

<br>

## 使用示例

```go
package main

import (
    "fmt"
    "time"

    "github.com/zhufuyi/pkg/gocron"
)

var task1 = func() {
     fmt.Println("this is task1")
     fmt.Println("running task list:", gocron.GetRunningTasks()) // 当前运行的任务
}

var taskOnce = func() {
	taskName := "taskOnce"
    fmt.Println("this is taskOnce")
    gocron.DeleteTask(taskName)  // 执行完删除任务
}

func main() {
    // 初始化
    err := gocron.Init()
    if err != nil {
        panic(err)
    }

    // 运行定时任务
    gocron.Run([]*gocron.Task{
        {
            Name:     "task1",
            TimeSpec: "@every 2s",
            Fn:       task1,
        },
        {
            Name:     "taskOnce",
            TimeSpec: "@every 5s",
            Fn:       taskOnce,
        },
    }...)

    time.Sleep(time.Minute)
}
```
