package mpipe

import (
	"errors"

	"github.com/ross96D/cancelreader"
)

type Transformer func([]byte) []byte

var NoTransform Transformer = func(b []byte) []byte {
	return b
}

type transformerReader struct {
	// reader          io.Reader
	cancelable      cancelreader.CancelReader
	transformerFunc Transformer
}

func (t transformerReader) cancel() bool {
	return t.cancelable.Cancel()
}

func (t transformerReader) applyTransformerFunc(p []byte, n int) (int, error) {
	writted := t.transformerFunc(p[:n])
	if len(writted) > cap(p) {
		return 0, errors.New("mpipe: transformer function result exceeds buffer capacity")
	}
	copy(p, writted)
	return len(writted), nil
}

func (t transformerReader) Read(p []byte) (int, error) {
	n, err := t.cancelable.Read(p)
	if err != nil {
		return n, err
	}

	if t.transformerFunc != nil {
		n, err = t.applyTransformerFunc(p, n)
		if err != nil {
			return 0, err
		}
	}
	return n, nil
}
