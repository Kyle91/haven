package routine

import (
	"fmt"
	"github.com/Kyle91/haven/log"
	"sync"
)

// Routine 管理goroutine的执行
type Routine struct {
	taskQueue chan func()
	wg        sync.WaitGroup
	maxTasks  int
}

var instance *Routine
var once sync.Once

// NewRoutine 创建一个新的 Routine 实例
func NewRoutine(maxTasks int) *Routine {
	once.Do(func() {
		instance = &Routine{
			taskQueue: make(chan func(), maxTasks),
			maxTasks:  maxTasks,
		}
		go instance.run()
	})
	return instance
}

// run 处理任务队列中的任务
func (r *Routine) run() {
	for task := range r.taskQueue {
		r.wg.Add(1)
		go func(task func()) {
			defer r.wg.Done()
			defer catchPanic()
			task()
		}(task)
	}
}

// Go 启动一个新的goroutine
func Go(logic func()) {
	r := getInstance()
	r.taskQueue <- logic
}

// Wait 等待所有goroutine完成
// 调用该方法后，将不再接受新的任务
func Wait() {
	gm := getInstance()
	close(gm.taskQueue)
	gm.wg.Wait()
}

// getInstance 获取GoManager的单例实例
func getInstance() *Routine {
	if instance == nil {
		NewRoutine(3000) // 默认任务数量
	}
	return instance
}

// catchPanic 捕获并处理goroutine中的panic
func catchPanic() {
	if err := recover(); err != nil {
		log.Info(fmt.Sprintf("panic recovered: %v", err))
	}
}
