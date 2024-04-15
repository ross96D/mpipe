package mpipe

import (
	"io"
	"os"
	"os/exec"
	"time"
)

type MpipeOptions func(*Mpipe)

func WithStdin(r io.Reader) MpipeOptions {
	return func(m *Mpipe) {
		m.stdin = r
	}
}

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

func WithTimeOut(timeout time.Duration) MpipeOptions {
	return func(m *Mpipe) {
		m.timeout = timeout
	}
}

type Mpipe struct {
	cmd               *exec.Cmd
	stdoutTransformer Transformer
	stderrTransformer Transformer

	stdin io.Reader

	readerOut transformerReader
	readerErr transformerReader

	stdoutCh chan struct{}
	stderrCh chan struct{}
	timeout  time.Duration
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

func (m *Mpipe) Start() (err error) {
	var thereader io.Reader
	if m.stdin != nil {
		thereader = m.stdin
	} else {
		thereader = os.Stdin
	}
	m.cmd.Stdin = thereader

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

	go func() {
		io.Copy(os.Stdout, m.readerOut)
		m.stdoutCh <- struct{}{}
	}()
	go func() {
		io.Copy(os.Stderr, m.readerErr)
		m.stderrCh <- struct{}{}
	}()

	return m.cmd.Start()
}

func (m *Mpipe) Wait() error {
	defer m.Cancel()
	err := m.cmd.Wait()
	tout := time.NewTimer(m.timeout)
	count := 2
	for {
		if count == 0 {
			return err
		}
		select {
		case <-tout.C:
			return err

		case <-m.stderrCh:
			count--

		case <-m.stdoutCh:
			count--
		}
	}
}

func (m *Mpipe) Cancel() bool {
	er := m.readerErr.cancel()
	out := m.readerOut.cancel()
	return er && out
}

func CommandWithOptions(cmd *exec.Cmd, opts ...MpipeOptions) *Mpipe {
	c := &Mpipe{
		cmd:      cmd,
		timeout:  20 * time.Millisecond,
		stdoutCh: make(chan struct{}),
		stderrCh: make(chan struct{}),
	}
	if opts != nil {
		for i := 0; i < len(opts); i++ {
			opts[i](c)
		}
	}
	c.checkTransfromers()
	return c
}

func (m *Mpipe) checkTransfromers() {
	if m.stdoutTransformer == nil {
		m.stdoutTransformer = NoTransform
	}
	if m.stderrTransformer == nil {
		m.stderrTransformer = NoTransform
	}
}
