package core

import (
	"path"
	"strings"

	"github.com/pkg/errors"
)

type Process struct {
	rootDir string
	window  Window
}

type ProcessConfig struct {
	RootDir string
}

func NewProcess(cfg ProcessConfig) *Process {
	p := &Process{
		rootDir: cfg.RootDir,
		window:  Window{title: "bedrock_server.exe"},
	}
	return p
}

func (p *Process) Active() bool {
	return p.window.IsRunning()
}

func (p *Process) ExecFile() string {
	return path.Join(p.rootDir, "bedrock_server.exe")
}

func (p *Process) Start() error {
	if !p.Active() {
		err := Command(p.ExecFile()).Run()
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return errors.New("server process is already running")
}

func (p *Process) Stop() error {
	if p.Active() {
		return p.window.Close()
	}
	return errors.New("server process is not running")
}

func (p *Process) Restart() error {
	err := p.Stop()
	if err != nil {
		return err
	}
	err = p.Start()
	if err != nil {
		return err
	}
	return nil
}

func (p *Process) ExecCmd(arg ...string) error {
	if !p.Active() {
		return errors.New("server process is not running")
	}
	return p.window.ExecuteCommand(strings.Join(arg, " "))
}
