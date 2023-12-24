package mpipe

import (
	"context"
	"os"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
)

type MpipeOptions func(*Mpipe)

func WithStyleOut(s *lipgloss.Style) func(*Mpipe) {
	return func(m *Mpipe) {
		m.styleOut = s
	}
}

func WithStyleErr(s *lipgloss.Style) func(*Mpipe) {
	return func(m *Mpipe) {
		m.styleErr = s
	}
}

type Mpipe struct {
	cmd      *exec.Cmd
	styleOut *lipgloss.Style
	styleErr *lipgloss.Style
}

func (m *Mpipe) StdoutTransformer() transformer {
	if m.styleOut == nil {
		return noTransform
	}
	return func(s string) string {
		return m.styleOut.Render(s)
	}
}

func (m *Mpipe) StderrTransformer() transformer {
	if m.styleErr == nil {
		return noTransform
	}
	return func(s string) string {
		return m.styleErr.Render(s)
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

	go transform(os.Stdout, stdout, m.StdoutTransformer())
	go transform(os.Stderr, stderr, m.StderrTransformer())
	go transform(stdin, os.Stdin, noTransform)

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

func CommandWithOptions(opts []MpipeOptions, name string, args ...string) *Mpipe {
	c := Command(name, args...)
	if opts != nil {
		for i := 0; i < len(opts); i++ {
			opts[i](c)
		}
	}
	return c
}

func CommandContext(ctx context.Context, name string, arg ...string) *Mpipe {
	return &Mpipe{
		cmd: exec.CommandContext(ctx, name, arg...),
	}
}

func CommandContextWithOptions(ctx context.Context, opts []MpipeOptions, name string, args ...string) *Mpipe {
	c := CommandContext(ctx, name, args...)
	if opts != nil {
		for i := 0; i < len(opts); i++ {
			opts[i](c)
		}
	}
	return c
}
