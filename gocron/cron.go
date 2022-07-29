package gocron

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
)

var (
	c      *cron.Cron
	nameID = sync.Map{} // 任务名称和id映射，用在对任务的增删改查
	idName = sync.Map{} // id和任务名称映射，用在日志打印
)

// Task 定时任务
type Task struct {
	// 秒(0-59) 分(0-59) 时(0-23) 日(1-31) 月(1-12) 星期(0-6)
	// "*/5 * * * * *"  表示每隔5秒执行
	// "0 15,45 9-12 * * * "  表示每天上午9点到12点的第15和第45分钟执行
	TimeSpec string

	Name string // 任务名称
	Fn   func() // 任务
}

// Init 初始化和启动定时任务
func Init(opts ...Option) error {
	o := defaultOptions()
	o.apply(opts...)

	log := &myLog{zapLog: o.zapLog}
	cronOpts := []cron.Option{
		cron.WithSeconds(), // 秒级粒度，默认是分钟级别粒度
		cron.WithLogger(log),
		cron.WithChain(
			cron.Recover(log), // or use cron.DefaultLogger
		),
	}

	c = cron.New(cronOpts...) // 实例化
	c.Start()                 // 启动定时器

	return nil
}

// Run 添加新的任务
func Run(tasks ...*Task) error {
	if c == nil {
		return errors.New("cron is not initialized")
	}

	var errs []string
	for _, task := range tasks {
		if IsRunningTask(task.Name) {
			errs = append(errs, fmt.Sprintf("task '%s' is already exists", task.Name))
			continue
		}

		id, err := c.AddFunc(task.TimeSpec, task.Fn)
		if err != nil {
			errs = append(errs, fmt.Sprintf("run task '%s' error: %v", task.Name, err))
			continue
		}
		idName.Store(id, task.Name)
		nameID.Store(task.Name, id)
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, " || "))
	}

	return nil
}

// IsRunningTask 判断任务是运行
func IsRunningTask(name string) bool {
	_, ok := nameID.Load(name)
	return ok
}

// GetRunningTasks 获取正在运行的任务名称列表
func GetRunningTasks() []string {
	var names []string
	nameID.Range(func(key, value interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}

// DeleteTask 删除任务
func DeleteTask(name string) {
	if id, ok := nameID.Load(name); ok {
		entryID := id.(cron.EntryID)
		c.Remove(entryID) // 从定时器中删除
		nameID.Delete(name)
		idName.Delete(entryID)
	}
}

// Stop 停止定时任务
func Stop() {
	if c != nil {
		c.Stop()
	}
}

// EverySecond 每隔size秒执行(1~59)
func EverySecond(size int) string {
	return fmt.Sprintf("@every %ds", size)
}

// EveryMinute 每隔size分钟执行(1~59)
func EveryMinute(size int) string {
	return fmt.Sprintf("@every %dm", size)
}

// EveryHour 每隔size小时执行(1~23)
func EveryHour(size int) string {
	return fmt.Sprintf("@every %dh", size)
}

// EveryDay 每隔size天执行(1~31)
func EveryDay(size int) string {
	return fmt.Sprintf("@every %dd", size)
}
