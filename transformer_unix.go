//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || aix
// +build darwin dragonfly freebsd linux netbsd openbsd solaris aix

package mpipe

import (
	"io"

	"github.com/ross96D/cancelreader"
)

type transformerReader struct {
	reader cancelreader.CancelReader
	// cancelable      cancelreader.CancelReader
	transformerFunc Transformer
}

func (t transformerReader) cancel() bool {
	return t.reader.Cancel()
}

func newTransformerReader(reader io.Reader, f Transformer) (transformerReader, error) {
	r, err := cancelreader.NewReader(reader)
	if err != nil {
		return transformerReader{}, err
	}
	return transformerReader{reader: r, transformerFunc: f}, nil
}
