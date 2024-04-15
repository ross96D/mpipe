package mpipe

import "io"

type Transformer func([]byte) []byte

var NoTransform Transformer = func(b []byte) []byte {
	return b
}

func (t transformerReader) applyTransformerFunc(p []byte, n int) (int, error) {
	writted := t.transformerFunc(p[:n])
	copied := copy(p, writted)
	return copied, nil
}

func (t transformerReader) Read(p []byte) (int, error) {
	n, err := t.reader.Read(p)
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

type NoCancel struct {
	reader io.Reader
}

func (nc NoCancel) Read(p []byte) (n int, err error) {
	return nc.reader.Read(p)
}

func (NoCancel) Close() error { return nil }

func (NoCancel) Cancel() bool { return true }

func newTransformerReaderWithoutCancel(reader io.Reader, f Transformer) (transformerReader, error) {
	return transformerReader{reader: NoCancel{reader: reader}, transformerFunc: f}, nil
}
