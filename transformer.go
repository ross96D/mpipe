package mpipe

import (
	"errors"
	"io"
)

type Transformer func([]byte) []byte

var NoTransform Transformer = func(b []byte) []byte {
	return b
}

type transformerReader struct {
	reader          io.Reader
	transformerFunc Transformer
}

func (t transformerReader) applyTransformerFunc(p []byte) ([]byte, error) {
	writted := make([]byte, 0)
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			break
		}
		writted = append(writted, p[i])
	}
	writted = t.transformerFunc(writted)
	if len(writted) > cap(p) {
		return nil, errors.New("mpipe: transformer function result exceeds buffer capacity")
	}
	copy(p, writted)
	return p, nil
}

func (t transformerReader) Read(p []byte) (int, error) {
	n, err := t.reader.Read(p)
	if err != nil {
		return n, err
	}

	if t.transformerFunc != nil {
		_, err := t.applyTransformerFunc(p)
		if err != nil {
			return 0, err
		}
	}
	return len(p), nil
}
