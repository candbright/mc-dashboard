package process

import (
	"os/exec"
)

type Screen struct {
	Name string
}

func (s Screen) Create() error {
	return exec.Command("screen", "-dmS", s.Name).Run()
}

func (s Screen) Active() bool {
	err := exec.Command("screen", "-ls", s.Name).Run()
	return err == nil
}

func (s Screen) Exit() error {
	return exec.Command("screen", "-X", "-S", s.Name, "quit").Run()
}

func (s Screen) ExecCmd(arg ...string) error {
	arg = append([]string{"-X", s.Name}, arg...)
	return exec.Command("screen", arg...).Run()
}
