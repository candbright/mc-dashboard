package core

import (
	"fmt"
	"github.com/candbright/go-server/internal/mc-server/utils"
	"github.com/pkg/errors"
	"path"
)

type Process struct {
	processId string
	rootDir   string
	screen    Tmux
}

type ProcessConfig struct {
	RootDir string
}

func NewProcess(cfg ProcessConfig) *Process {
	randomId := utils.RandomString(8)
	p := &Process{
		processId: randomId,
		rootDir:   cfg.RootDir,
		screen:    Tmux{Name: fmt.Sprintf("mc-%s", randomId)},
	}

	return p
}

func (p *Process) Active() bool {
	return p.screen.Exists()
}

func (p *Process) ExecFile() string {
	return path.Join(p.rootDir, "bedrock_server")
}

func (p *Process) ScreenName() string {
	return fmt.Sprintf("mc-%s", p.processId)
}

func (p *Process) Start() error {
	if !p.Active() {
		err := p.screen.Create()
		if err != nil {
			return err
		}
		err = p.screen.ExecCmd("LD_LIBRARY_PATH=.", p.ExecFile())
		if err != nil {
			return err
		}
	}
	return errors.New("server process is already running")
}

func (p *Process) Stop() error {
	if p.Active() {
		return p.screen.Exit()
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
	return p.screen.ExecCmd(arg...)
}
