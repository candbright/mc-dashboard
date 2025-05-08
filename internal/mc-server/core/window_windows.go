package core

import (
	"fmt"
	"github.com/pkg/errors"
	"syscall"
	"unsafe"
)

// 定义 Windows API 所需的常量、结构体和函数
const (
	TH32CS_SNAPPROCESS = 0x00000002
	MAX_PATH           = 260
	PROCESS_ALL_ACCESS = 0x1F0FFF
	MEM_COMMIT         = 0x1000
	MEM_RESERVE        = 0x2000
	PAGE_READWRITE     = 0x04
	PROCESS_TERMINATE  = 0x0001
)

type ProcessEntry32 struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [MAX_PATH]uint16
}

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	loadLibraryA                 = kernel32.NewProc("LoadLibraryA")
	procCreateToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = kernel32.NewProc("Process32FirstW")
	procProcess32Next            = kernel32.NewProc("Process32NextW")
	procOpenProcess              = kernel32.NewProc("OpenProcess")
	procVirtualAllocEx           = kernel32.NewProc("VirtualAllocEx")
	procWriteProcessMemory       = kernel32.NewProc("WriteProcessMemory")
	procCreateRemoteThread       = kernel32.NewProc("CreateRemoteThread")
	procTerminateProcess         = kernel32.NewProc("TerminateProcess")
)

type Window struct {
	title        string
	processEntry *ProcessEntry32
}

func (w *Window) ProcessEntry() (*ProcessEntry32, error) {
	if w.processEntry == nil {
		// 创建进程快照
		snapshot, _, err := procCreateToolhelp32Snapshot.Call(
			TH32CS_SNAPPROCESS,
			0,
		)
		if snapshot == uintptr(syscall.InvalidHandle) {
			return nil, errors.WithStack(err)
		}
		defer syscall.CloseHandle(syscall.Handle(snapshot))

		var entry ProcessEntry32
		entry.Size = uint32(unsafe.Sizeof(entry))

		// 遍历进程列表
		ret, _, err := procProcess32First.Call(snapshot, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			return nil, errors.WithStack(err)
		}

		for {
			// 提取进程名并比较
			name := syscall.UTF16ToString(entry.ExeFile[:])
			if name == w.title {
				w.processEntry = &entry
				break
			}

			// 继续下一个进程
			ret, _, _ = procProcess32Next.Call(snapshot, uintptr(unsafe.Pointer(&entry)))
			if ret == 0 {
				break
			}
		}
		if w.processEntry == nil {
			return nil, errors.New("process entry not found")
		}
	}

	return w.processEntry, nil
}

func (w *Window) IsRunning() bool {
	processEntry, _ := w.ProcessEntry()
	return processEntry != nil
}

func (w *Window) ExecuteCommand(command string) error {
	processEntry, err := w.ProcessEntry()
	if err != nil {
		return err
	}
	// 4. 如果找到进程且需要执行命令
	if command != "" {
		err = w.executeCommand(processEntry.ProcessID, command)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	return nil
}

// 向进程执行命令（通过内存注入）
func (w *Window) executeCommand(pid uint32, command string) error {
	// 1. 打开目标进程
	hProcess, _, err := procOpenProcess.Call(
		PROCESS_ALL_ACCESS,
		0,
		uintptr(pid),
	)
	if hProcess == 0 {
		return fmt.Errorf("打开进程失败: %v", err)
	}
	defer syscall.CloseHandle(syscall.Handle(hProcess))

	// 2. 在目标进程分配内存
	remoteMem, _, err := procVirtualAllocEx.Call(
		hProcess,
		0,
		uintptr(len(command)+1),
		MEM_COMMIT|MEM_RESERVE,
		PAGE_READWRITE,
	)
	if remoteMem == 0 {
		return fmt.Errorf("内存分配失败: %v", err)
	}

	// 3. 写入命令到目标进程
	commandBytes := append([]byte(command), 0) // 添加 null 终止符
	var bytesWritten uintptr
	_, _, err = procWriteProcessMemory.Call(
		hProcess,
		remoteMem,
		uintptr(unsafe.Pointer(&commandBytes[0])),
		uintptr(len(commandBytes)),
		uintptr(unsafe.Pointer(&bytesWritten)),
	)
	if bytesWritten != uintptr(len(commandBytes)) {
		return fmt.Errorf("写入内存失败: %v", err)
	}

	// 4. 创建远程线程执行命令

	var threadId uint32
	hThread, _, err := procCreateRemoteThread.Call(
		hProcess,
		0,
		0,
		loadLibraryA.Addr(),
		remoteMem,
		0,
		uintptr(unsafe.Pointer(&threadId)),
	)
	if hThread == 0 {
		return fmt.Errorf("创建远程线程失败: %v", err)
	}
	defer syscall.CloseHandle(syscall.Handle(hThread))

	return nil
}

// 终止进程
func (w *Window) terminateProcess(pid uint32) error {
	hProcess, _, err := procOpenProcess.Call(
		PROCESS_TERMINATE,
		0,
		uintptr(pid),
	)
	if hProcess == 0 {
		return fmt.Errorf("打开进程失败: %v", err)
	}
	defer syscall.CloseHandle(syscall.Handle(hProcess))

	ret, _, err := procTerminateProcess.Call(
		hProcess,
		1, // 退出代码
	)
	if ret == 0 {
		return fmt.Errorf("终止进程失败: %v", err)
	}

	return nil
}

func (w *Window) Close() error {
	processEntry, err := w.ProcessEntry()
	if err != nil {
		return err
	}
	err = w.terminateProcess(processEntry.ProcessID)
	if err != nil {
		return err
	}
	w.processEntry = nil
	return nil
}
