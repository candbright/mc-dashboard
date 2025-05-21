package process

import (
	"fmt"
	"os/exec"
	"strings"
)

type Tmux struct {
	Name string
}

func (s Tmux) Create() error {
	return exec.Command("tmux", "new", "-s", s.Name).Run()
}

func (s Tmux) Active() bool {
	err := exec.Command("tmux", "has-session", "-t", s.Name).Run()
	return err == nil
}

func (s Tmux) Exit() error {
	return exec.Command("tmux", "kill-session", "-t", s.Name).Run()
}

func (s Tmux) ExecCmd(arg ...string) error {
	return exec.Command("tmux", "send-keys", "-t", s.Name,
		fmt.Sprintf("\"%s\"", strings.Join(arg, " ")), "Enter").Run()
}
