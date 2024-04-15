package mpipe

import (
	"fmt"
	"io"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMpipe(t *testing.T) {
	cmd := exec.Command("psql", "--help")
	command := CommandWithOptions(cmd, WithStderrTransformer(func(s []byte) []byte {
		return append([]byte("preferr: "), s...)
	}), WithStdoutTransformer(func(s []byte) []byte {
		return append([]byte("prefout: "), s...)
	}))

	err := command.Run()
	require.Equalf(t, nil, err, "error: %s", err)
}

func TestEcho(t *testing.T) {
	cmd := exec.Command("echo", "Hello World")
	command := CommandWithOptions(cmd, WithStderrTransformer(func(s []byte) []byte {
		return append([]byte("preferr: "), s...)
	}), WithStdoutTransformer(func(s []byte) []byte {
		return append([]byte("prefout: "), s...)
	}))

	err := command.Run()
	require.Equal(t, nil, err)
}

func TestMpipePy(t *testing.T) {
	pr, pw := io.Pipe()

	cmd := exec.Command("python3", "./test.py")
	command := CommandWithOptions(cmd,
		WithStderrTransformer(func(s []byte) []byte {
			return append([]byte("err: "), s...)
		}),
		WithStdoutTransformer(func(s []byte) []byte {
			return append([]byte("out: "), s...)
		}),
		WithStdin(pr),
	)

	channelCmd := make(chan error)
	channelInput := make(chan error)
	timeout := time.NewTimer(1 * time.Second)
	go func() {
		err := command.Run()
		channelCmd <- err
	}()

	go func() {
		for i := 0; i < 50; i++ {
			_, err := pw.Write([]byte(fmt.Sprintf("%d\n", i)))
			if err != nil {
				channelInput <- err
			}
		}
		_, err := pw.Write([]byte("close\n"))
		pw.Close()
		channelInput <- err
	}()

loop:
	for {
		select {
		case inputErr := <-channelInput:
			require.Equal(t, nil, inputErr)
			fmt.Printf("input success\n")

		case cmdErr := <-channelCmd:
			require.Equal(t, nil, cmdErr)
			fmt.Printf("cmd success\n")
			break loop

		case <-timeout.C:
			require.FailNow(t, "timeout")
		}
	}
}
