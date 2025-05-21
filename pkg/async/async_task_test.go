package async

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// 测试正常任务执行
func TestAsyncTaskExecutor_NormalExecution(t *testing.T) {
	executor := NewAsyncTaskExecutor()

	// 添加测试任务
	taskCount := 3
	for i := 0; i < taskCount; i++ {
		taskNum := i + 1
		executor.AddTask(NewSimpleTask(
			fmt.Sprintf("任务%d", taskNum),
			func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			},
		))
	}

	// 开始执行
	executor.Start()

	// 检查执行状态
	var finalStatus TaskStatus
	for status := range executor.Status() {
		if status.Err != nil {
			t.Fatalf("执行出错: %v", status.Err)
		}
		finalStatus = status

		if !status.IsRunning && status.CompletedTasks == taskCount {
			break
		}
	}

	// 验证最终状态
	if finalStatus.Percentage != 100 {
		t.Errorf("期望完成度100%%, 实际得到 %.2f%%", finalStatus.Percentage)
	}
	if finalStatus.CompletedTasks != taskCount {
		t.Errorf("期望完成 %d 个任务, 实际完成 %d 个", taskCount, finalStatus.CompletedTasks)
	}
}

// 测试任务执行错误
func TestAsyncTaskExecutor_ExecutionError(t *testing.T) {
	executor := NewAsyncTaskExecutor()
	expectedErr := errors.New("测试错误")

	// 添加会失败的任务
	executor.AddTask(NewSimpleTask("失败任务", func() error {
		return expectedErr
	}))

	// 添加不会执行的任务
	executor.AddTask(NewSimpleTask("未执行任务", func() error {
		return nil
	}))

	// 开始执行
	executor.Start()

	// 检查执行状态
	var gotError bool
	for status := range executor.Status() {
		if status.Err != nil {
			gotError = true
			if status.Err != expectedErr {
				t.Errorf("期望错误 %v, 实际得到 %v", expectedErr, status.Err)
			}
			break
		}
	}

	if !gotError {
		t.Error("期望得到错误，但执行成功完成")
	}
}

// 测试任务停止
func TestAsyncTaskExecutor_StopExecution(t *testing.T) {
	executor := NewAsyncTaskExecutor()

	// 添加长时间运行的任务
	executor.AddTask(NewSimpleTask("长时间任务", func() error {
		time.Sleep(1 * time.Second)
		return nil
	}))

	// 开始执行
	executor.Start()

	// 等待一段时间后停止
	time.Sleep(100 * time.Millisecond)
	executor.Stop()

	// 检查是否收到停止信号
	var gotStop bool
	for status := range executor.Status() {
		if status.Err != nil && status.CurrentTask == "已停止" {
			gotStop = true
			break
		}
	}

	if !gotStop {
		t.Error("未能成功停止任务")
	}
}

// 测试进度报告
func TestAsyncTaskExecutor_ProgressReporting(t *testing.T) {
	executor := NewAsyncTaskExecutor()

	// 添加多个任务
	taskCount := 5
	for i := 0; i < taskCount; i++ {
		taskNum := i + 1
		executor.AddTask(NewSimpleTask(
			fmt.Sprintf("任务%d", taskNum),
			func() error {
				time.Sleep(50 * time.Millisecond)
				return nil
			},
		))
	}

	// 开始执行
	executor.Start()

	// 收集进度报告
	var progressReports []TaskStatus
	for status := range executor.Status() {
		progressReports = append(progressReports, status)
		if !status.IsRunning && status.CompletedTasks == taskCount {
			break
		}
	}

	// 检查是否收到多个进度更新
	if len(progressReports) < 2 {
		t.Errorf("期望至少2个进度报告, 实际得到 %d", len(progressReports))
	}

	// 检查进度是否递增
	var lastCompleted int
	for i, report := range progressReports {
		if i > 0 && report.CompletedTasks < lastCompleted {
			t.Errorf("已完成任务数不应减少: 前一个 %d, 当前 %d", lastCompleted, report.CompletedTasks)
		}
		lastCompleted = report.CompletedTasks
	}

	// 检查最终完成状态
	lastReport := progressReports[len(progressReports)-1]
	if lastReport.Percentage != 100 {
		t.Errorf("期望最终进度100%%, 实际得到 %.2f%%", lastReport.Percentage)
	}
}

// 测试并发状态访问
func TestAsyncTaskExecutor_ConcurrentStatusAccess(t *testing.T) {
	executor := NewAsyncTaskExecutor()

	// 添加任务
	executor.AddTask(NewSimpleTask("测试任务", func() error {
		time.Sleep(100 * time.Millisecond)
		return nil
	}))

	// 开始执行
	executor.Start()

	// 并发获取状态
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				status := executor.GetCurrentStatus()
				if status.TotalTasks != 1 {
					t.Errorf("期望总任务数1, 实际得到 %d", status.TotalTasks)
				}
			}
			done <- struct{}{}
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 等待任务完成
	for status := range executor.Status() {
		if !status.IsRunning {
			break
		}
	}
}
