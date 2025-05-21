package process

import (
	"runtime"
)

type Process interface {
	Exist() bool
	Active() bool
	Start() error
	Stop() error
	Restart() error
	ExecCmd(arg ...string) error
	ScanLog(line int) (string, error)
}

type ProcessConfig struct {
	OS         string
	DriverType string
	ID         string
	RootDir    string
	LogFile    string
}

func NewProcess(cfg ProcessConfig) Process {
	if cfg.OS == "" {
		cfg.OS = runtime.GOOS
	}
	switch cfg.OS {
	case "linux":
		return NewLinuxProcess(cfg)
	case "windows":
		return NewWindowsProcess(cfg)
	default:
		return NewLinuxProcess(cfg)
	}
}
