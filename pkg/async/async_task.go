package async

import (
	"fmt"
	"sync"
)

// TaskStatus 任务状态
type TaskStatus struct {
	TotalTasks     int     `json:"total_tasks"`     // 总任务数
	CompletedTasks int     `json:"completed_tasks"` // 已完成任务数
	Percentage     float64 `json:"percentage"`      // 完成百分比
	CurrentTask    string  `json:"current_task"`    // 当前执行的任务
	IsRunning      bool    `json:"is_running"`      // 是否正在运行
	Err            error   `json:"err,omitempty"`   // 错误信息
}

// ErrorCallback 错误回调函数类型
type ErrorCallback func(task Task, err error)

// Task 任务接口
type Task interface {
	Execute() error
	Name() string
}

// AsyncTaskExecutor 异步任务执行器
type AsyncTaskExecutor struct {
	tasks         []Task
	status        TaskStatus
	statusMutex   sync.Mutex
	statusChan    chan TaskStatus
	stopChan      chan struct{}
	errorCallback ErrorCallback
}

// NewAsyncTaskExecutor 创建新的异步任务执行器
func NewAsyncTaskExecutor() *AsyncTaskExecutor {
	return &AsyncTaskExecutor{
		statusChan: make(chan TaskStatus, 10),
		stopChan:   make(chan struct{}),
	}
}

// SetErrorCallback 设置错误回调函数
func (e *AsyncTaskExecutor) SetErrorCallback(callback ErrorCallback) {
	e.errorCallback = callback
}

// AddTask 添加任务
func (e *AsyncTaskExecutor) AddTask(task Task) {
	e.tasks = append(e.tasks, task)
	e.statusMutex.Lock()
	e.status.TotalTasks = len(e.tasks)
	e.statusMutex.Unlock()
}

// Start 开始执行任务
func (e *AsyncTaskExecutor) Start() {
	go func() {
		defer close(e.statusChan)
		e.updateStatus(TaskStatus{IsRunning: true})

		for i, task := range e.tasks {
			select {
			case <-e.stopChan:
				e.updateStatus(TaskStatus{
					TotalTasks:     e.status.TotalTasks,
					CompletedTasks: i,
					CurrentTask:    "已停止",
					IsRunning:      false,
					Err:            fmt.Errorf("任务被停止"),
				})
				return
			default:
				// 更新当前任务状态
				e.updateStatus(TaskStatus{
					TotalTasks:     e.status.TotalTasks,
					CompletedTasks: i,
					CurrentTask:    task.Name(),
					IsRunning:      true,
				})

				// 同步执行任务
				if err := task.Execute(); err != nil {
					status := TaskStatus{
						TotalTasks:     e.status.TotalTasks,
						CompletedTasks: i,
						CurrentTask:    task.Name(),
						IsRunning:      false,
						Err:            err,
					}
					e.updateStatus(status)

					// 调用错误回调函数
					if e.errorCallback != nil {
						e.errorCallback(task, err)
					}
					return
				}

				// 任务完成后更新状态
				e.updateStatus(TaskStatus{
					TotalTasks:     e.status.TotalTasks,
					CompletedTasks: i + 1,
					CurrentTask:    task.Name(),
					IsRunning:      true,
				})
			}
		}

		// 所有任务执行完成
		e.updateStatus(TaskStatus{
			TotalTasks:     e.status.TotalTasks,
			CompletedTasks: len(e.tasks),
			CurrentTask:    "完成",
			IsRunning:      false,
		})
	}()
}

// Stop 停止执行任务
func (e *AsyncTaskExecutor) Stop() {
	close(e.stopChan)
}

// Status 获取状态通道
func (e *AsyncTaskExecutor) Status() <-chan TaskStatus {
	return e.statusChan
}

// GetCurrentStatus 获取当前状态
func (e *AsyncTaskExecutor) GetCurrentStatus() TaskStatus {
	e.statusMutex.Lock()
	defer e.statusMutex.Unlock()
	return e.status
}

// updateStatus 更新状态
func (e *AsyncTaskExecutor) updateStatus(status TaskStatus) {
	e.statusMutex.Lock()
	defer e.statusMutex.Unlock()

	// 保留原有的 TotalTasks
	if e.status.TotalTasks > 0 {
		status.TotalTasks = e.status.TotalTasks
	}

	// 计算百分比
	if status.TotalTasks > 0 {
		status.Percentage = float64(status.CompletedTasks) / float64(status.TotalTasks) * 100
	}

	e.status = status
	select {
	case e.statusChan <- status:
	default: // 避免阻塞
	}
}

// SimpleTask 简单任务实现
type SimpleTask struct {
	name     string
	execFunc func() error
}

// NewSimpleTask 创建简单任务
func NewSimpleTask(name string, execFunc func() error) *SimpleTask {
	return &SimpleTask{
		name:     name,
		execFunc: execFunc,
	}
}

func (t *SimpleTask) Execute() error {
	return t.execFunc()
}

func (t *SimpleTask) Name() string {
	return t.name
}
