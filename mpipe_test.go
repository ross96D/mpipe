package mpipe

import (
	"os/exec"
	"testing"
)

func TestMpipe(t *testing.T) {
	cmd := exec.Command("ls")
	command := CommandWithOptions(cmd, WithStderrTransformer(func(s []byte) []byte {
		return append([]byte("preferr: "), s...)
	}), WithStdoutTransformer(func(s []byte) []byte {
		return append([]byte("prefout: "), s...)
	}))
	command.Run()
}
