package mpipe

import (
	"io"
	"os"
	"os/exec"
	"time"
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

	readerOut transformerReader
	readerIn  transformerReader
	readerErr transformerReader

	stdout chan struct{}
	stderr chan struct{}
	stdin  chan struct{}
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

	if m.readerOut, err = newTransformerReader(stdout, m.stdoutTransformer); err != nil {
		return err
	}

	if m.readerErr, err = newTransformerReader(stderr, m.stderrTransformer); err != nil {
		return err
	}

	if m.readerIn, err = newTransformerReader(os.Stdin, m.stdinTransformer); err != nil {
		return err
	}

	go func() {
		io.Copy(os.Stdout, m.readerOut)
		<-m.stdout
	}()
	go func() {
		io.Copy(os.Stderr, m.readerErr)
		<-m.stderr
	}()
	go func() {
		io.Copy(stdin, m.readerIn)
		<-m.stdin
	}()

	return m.cmd.Start()
}

func (m *Mpipe) Wait() error {
	defer m.Cancel()
	err := m.cmd.Wait()
	tout := time.NewTimer(500 * time.Millisecond)
	count := 3
	for {
		if count == 0 {
			return err
		}
		select {
		case <-tout.C:
			return err

		case <-m.stderr:
			count--

		case <-m.stdout:
			count--

		case <-m.stdin:
			count--
		}
	}
}

func (m *Mpipe) Cancel() bool {
	er := m.readerErr.cancel()
	out := m.readerOut.cancel()
	in := m.readerIn.cancel()
	return er && out && in
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
