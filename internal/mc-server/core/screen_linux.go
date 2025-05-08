package core

import (
	"github.com/pkg/errors"
)

type Screen struct {
	Name string
}

func (s Screen) Create() error {
	return errors.WithStack(Command("screen", "-dmS", s.Name).Run())
}

func (s Screen) Exists() bool {
	err := Command("screen", "-ls", s.Name).Run()
	return err == nil
}

func (s Screen) Exit() error {
	return errors.WithStack(Command("screen", "-X", "-S", s.Name, "quit").Run())
}

func (s Screen) ExecCmd(arg ...string) error {
	arg = append([]string{"-X", s.Name}, arg...)
	return errors.WithStack(Command("screen", arg...).Run())
}
