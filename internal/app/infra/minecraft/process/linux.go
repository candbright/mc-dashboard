package process

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/candbright/mc-dashboard/internal/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Driver interface {
	Create() error
	Active() bool
	Exit() error
	ExecCmd(arg ...string) error
}

type LinuxProcess struct {
	execFile     string
	runFile      string
	checkRunFile bool
	logFile      string
	cmd          *exec.Cmd
	stdin        io.WriteCloser
	pid          int
}

func NewLinuxProcess(cfg ProcessConfig) Process {
	p := &LinuxProcess{
		execFile:     path.Join(cfg.RootDir, "bedrock_server"),
		runFile:      path.Join(cfg.RootDir, fmt.Sprintf("bedrock_server_%s", cfg.ID)),
		checkRunFile: true,
		logFile:      cfg.LogFile,
	}

	// 尝试查找并接管已存在的进程
	if pid := p.findExistingProcess(); pid != 0 {
		p.pid = pid
		p.attachToProcess()
	}

	return p
}

// findExistingProcess 查找已存在的服务器进程
func (p *LinuxProcess) findExistingProcess() int {
	cmd := exec.Command("ps", "-ef")
	output, err := cmd.Output()
	if err != nil {
		logrus.WithError(err).Debug("Failed to list processes")
		return 0
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, p.runFile) {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				pid, err := strconv.Atoi(fields[1])
				if err != nil {
					continue
				}
				return pid
			}
		}
	}
	return 0
}

// attachToProcess 接管已存在的进程
func (p *LinuxProcess) attachToProcess() {
	if p.pid == 0 {
		return
	}

	process, err := os.FindProcess(p.pid)
	if err != nil {
		logrus.WithError(err).Error("Failed to find process")
		return
	}

	// 创建新的命令对象
	p.cmd = exec.Command(p.runFile)
	p.cmd.Process = process

	// 尝试获取标准输入
	stdin, err := os.OpenFile(fmt.Sprintf("/proc/%d/fd/0", p.pid), os.O_WRONLY, 0)
	if err == nil {
		p.stdin = stdin
	} else {
		logrus.WithError(err).Warn("Failed to attach to process stdin")
	}

	logrus.WithField("pid", p.pid).Info("Successfully attached to existing process")
}

func (p *LinuxProcess) Exist() bool {
	return utils.FileExist(p.execFile)
}

// checkProcessActive 检查进程是否在运行
func (p *LinuxProcess) checkProcessActive(pid int) bool {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid))
	err := cmd.Run()
	return err == nil
}

func (p *LinuxProcess) Active() bool {
	// 如果进程ID为0，说明进程未启动
	if p.pid == 0 {
		return false
	}

	// 如果是新启动的进程，检查 cmd 对象
	if p.cmd != nil && p.cmd.Process != nil {
		// 检查进程状态
		if p.cmd.ProcessState != nil {
			return false
		}

		// 检查进程是否真的在运行
		err := p.cmd.Process.Signal(syscall.Signal(0))
		if err == nil {
			return true
		}
	}

	// 如果是接管的进程或 cmd 对象无效，直接检查进程ID
	return p.checkProcessActive(p.pid)
}

func (p *LinuxProcess) Start() error {
	if p.Active() {
		return errors.New("server process is already running")
	}

	if p.checkRunFile && p.Exist() && !utils.FileExist(p.runFile) {
		if err := utils.CopyFile(p.execFile, p.runFile); err != nil {
			logrus.WithError(err).Error("Failed to copy server executable")
		}
		p.checkRunFile = false
	}

	// 创建命令
	p.cmd = exec.Command(p.runFile)

	// 设置进程组，这样可以在停止时终止整个进程组
	p.cmd.SysProcAttr = &syscall.SysProcAttr{
		// 使用默认值，让系统自动处理进程组
	}

	// 创建日志文件
	logFile, err := os.OpenFile(p.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// 设置标准输出和错误输出到日志文件
	p.cmd.Stdout = logFile
	p.cmd.Stderr = logFile

	// 获取标准输入管道
	p.stdin, err = p.cmd.StdinPipe()
	if err != nil {
		logFile.Close()
		return fmt.Errorf("failed to get stdin pipe: %v", err)
	}

	if err := p.cmd.Start(); err != nil {
		logFile.Close()
		logrus.WithError(err).Error("Failed to start server process")
		return err
	}

	// 保存进程ID
	p.pid = p.cmd.Process.Pid
	logrus.WithField("pid", p.pid).Info("Server process started")

	return nil
}

func (p *LinuxProcess) Stop() error {
	if !p.Active() {
		return errors.New("server process is not running")
	}

	// 关闭标准输入
	if p.stdin != nil {
		p.stdin.Close()
	}

	// 如果是通过 cmd 启动的进程
	if p.cmd != nil && p.cmd.Process != nil {
		// 发送SIGTERM信号到进程
		if err := p.cmd.Process.Signal(syscall.SIGTERM); err != nil {
			// 如果发送信号失败，强制终止进程
			return p.cmd.Process.Kill()
		}

		// 等待进程结束
		done := make(chan error, 1)
		go func() {
			done <- p.cmd.Wait()
		}()

		// 设置超时时间
		select {
		case err := <-done:
			if err != nil {
				logrus.WithError(err).Warn("Server process exited with error")
			}
		case <-time.After(10 * time.Second):
			// 如果超时，强制终止进程
			if err := p.cmd.Process.Kill(); err != nil {
				return fmt.Errorf("failed to kill process: %v", err)
			}
		}
	} else {
		// 如果是接管的进程，直接使用 kill 命令终止
		cmd := exec.Command("kill", "-9", strconv.Itoa(p.pid))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to kill process: %v", err)
		}
	}

	// 清除进程ID
	p.pid = 0
	return nil
}

func (p *LinuxProcess) Restart() error {
	err := p.Stop()
	if err != nil {
		return err
	}
	return p.Start()
}

func (p *LinuxProcess) ExecCmd(arg ...string) error {
	if !p.Active() {
		return errors.New("server process is not running")
	}

	// 构建命令字符串
	cmd := strings.Join(arg, " ") + "\n"

	// 写入命令
	_, err := p.stdin.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to write command: %v", err)
	}

	return nil
}

func (p *LinuxProcess) ScanLog(line int) (string, error) {
	// 参数验证
	if line <= 0 {
		return "", fmt.Errorf("invalid line number: %d", line)
	}

	// 打开日志文件
	file, err := os.OpenFile(p.logFile, os.O_RDONLY, 0666)
	if err != nil {
		return "", fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}
	fileSize := fileInfo.Size()

	// 创建缓冲区
	buffer := make([]byte, fileSize)

	// 读取整个文件
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// 将文件内容按行分割
	lines := strings.Split(string(buffer), "\n")

	// 如果请求的行数大于文件总行数，返回所有行
	if line >= len(lines) {
		return strings.Join(lines, "\n"), nil
	}

	// 返回最后 line 行
	return strings.Join(lines[len(lines)-line:], "\n"), nil
}
