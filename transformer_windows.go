//go:build windows
// +build windows

package mpipe

import (
	"io"
)

type transformerReader struct {
	reader io.Reader
	// cancelable      cancelreader.CancelReader
	transformerFunc Transformer
}

func (t transformerReader) cancel() bool {
	return false
}

func newTransformerReader(reader io.Reader, f Transformer) (transformerReader, error) {
	return transformerReader{reader: reader, transformerFunc: f}, nil
}
