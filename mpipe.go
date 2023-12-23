package mpipe

import (
	"context"
	"os/exec"
)

type MpipeOptions struct {
}

type Mpipe struct {
	cmd *exec.Cmd
}

func (m *Mpipe) String() string {
	return m.cmd.String()
}

func (m *Mpipe) Run() error {
	if err := m.Start(); err != nil {
		return err
	}
	return m.Wait()
}

func (m *Mpipe) Start() error {
	return m.cmd.Start()
}

func (m *Mpipe) Wait() error {
	return m.cmd.Wait()
}

func Command(name string, args ...string) *Mpipe {
	return &Mpipe{
		cmd: exec.Command(name, args...),
	}
}

func CommandContext(ctx context.Context, name string, arg ...string) *Mpipe {
	return &Mpipe{
		cmd: exec.CommandContext(ctx, name, arg...),
	}
}
