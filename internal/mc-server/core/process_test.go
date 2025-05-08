package core

import (
	"testing"
)

func TestProcess_Active(t *testing.T) {
	process := testServer(t).process
	active := process.Active()
	t.Log(active)
}

func TestProcess_Start(t *testing.T) {
	process := testServer(t).process
	err := process.Start()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	active := process.Active()
	t.Log(active)
}

func TestProcess_Stop(t *testing.T) {
	process := testServer(t).process
	err := process.Stop()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	active := process.Active()
	t.Log(active)
}
