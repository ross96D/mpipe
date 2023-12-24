package mpipe

import (
	"context"
	"os"
	"os/exec"
)

type MpipeOptions func(*Mpipe)

func WithStdoutTransformer(t Transformer) MpipeOptions {
	return func(m *Mpipe) {
		m.stdoutTransformer = t
	}
}

func WithStderrTransformer(t Transformer) MpipeOptions {
	return func(m *Mpipe) {
		m.stderrTransformer = t
	}
}

func WithStdinTransformer(t Transformer) MpipeOptions {
	return func(m *Mpipe) {
		m.stdinTransformer = t
	}
}

type Mpipe struct {
	cmd               *exec.Cmd
	stdoutTransformer Transformer
	stderrTransformer Transformer
	stdinTransformer  Transformer
}

func (m *Mpipe) checkTransfromers() {
	if m.stdoutTransformer == nil {
		m.stdoutTransformer = NoTransform
	}
	if m.stderrTransformer == nil {
		m.stderrTransformer = NoTransform
	}
	if m.stdinTransformer == nil {
		m.stdinTransformer = NoTransform
	}
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
	stdin, err := m.cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := m.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := m.cmd.StderrPipe()
	if err != nil {
		return err
	}

	go transform(os.Stdout, stdout, m.stdoutTransformer)
	go transform(os.Stderr, stderr, m.stderrTransformer)
	go transform(stdin, os.Stdin, m.stdinTransformer)

	return m.cmd.Start()
}

func (m *Mpipe) Wait() error {
	return m.cmd.Wait()
}

func Command(name string, args ...string) *Mpipe {
	return CommandWithOptions(nil, name, args...)
}

func CommandWithOptions(opts []MpipeOptions, name string, args ...string) *Mpipe {
	c := &Mpipe{
		cmd: exec.Command(name, args...),
	}
	if opts != nil {
		for i := 0; i < len(opts); i++ {
			opts[i](c)
		}
	}
	c.checkTransfromers()
	return c
}

func CommandContext(ctx context.Context, name string, args ...string) *Mpipe {
	return CommandContextWithOptions(ctx, nil, name, args...)
}

func CommandContextWithOptions(ctx context.Context, opts []MpipeOptions, name string, args ...string) *Mpipe {
	c := &Mpipe{
		cmd: exec.CommandContext(ctx, name, args...),
	}
	if opts != nil {
		for i := 0; i < len(opts); i++ {
			opts[i](c)
		}
	}
	c.checkTransfromers()
	return c
}
