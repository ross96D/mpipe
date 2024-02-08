package mpipe

import (
	"os/exec"
	"testing"
)

func TestMpipe(t *testing.T) {
	cmd := exec.Command("psql", "-U", "postgres")
	command := CommandWithOptions(cmd, WithStderrTransformer(func(s []byte) []byte {
		return append([]byte("preferr: "), s...)
	}), WithStdoutTransformer(func(s []byte) []byte {
		return append([]byte("prefout: "), s...)
	}))

	err := command.Run()
	if err != nil {
		t.Error(err)
	}

}
