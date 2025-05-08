package core

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type Tmux struct {
	Name string
}

func (s Tmux) Create() error {
	return errors.WithStack(Command("tmux", "new", "-s", s.Name).Run())
}

func (s Tmux) Exists() bool {
	err := Command("tmux", "has-session", "-t", s.Name).Run()
	return err == nil
}

func (s Tmux) Exit() error {
	return errors.WithStack(Command("tmux", "kill-session", "-t", s.Name).Run())
}

func (s Tmux) ExecCmd(arg ...string) error {
	return errors.WithStack(Command("tmux", "send-keys", "-t", s.Name,
		fmt.Sprintf("\"%s\"", strings.Join(arg, " ")), "Enter").Run())
}
