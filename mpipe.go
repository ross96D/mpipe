package mpipe

import (
	"io"
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
	pipeStdout := transformerReader{reader: stdout, transformerFunc: m.stdoutTransformer}

	stderr, err := m.cmd.StderrPipe()
	if err != nil {
		return err
	}
	pipeStderr := transformerReader{reader: stderr, transformerFunc: m.stderrTransformer}

	pipeOsStdin := transformerReader{reader: os.Stdin, transformerFunc: m.stdinTransformer}

	go io.Copy(os.Stdout, pipeStdout)
	go io.Copy(os.Stderr, pipeStderr)
	go io.Copy(stdin, pipeOsStdin)

	return m.cmd.Start()
}

func (m *Mpipe) Wait() error {
	return m.cmd.Wait()
}

func CommandWithOptions(cmd *exec.Cmd, opts ...MpipeOptions) *Mpipe {
	c := &Mpipe{
		cmd: cmd,
	}
	if opts != nil {
		for i := 0; i < len(opts); i++ {
			opts[i](c)
		}
	}
	io.Pipe()
	c.checkTransfromers()
	return c
}
