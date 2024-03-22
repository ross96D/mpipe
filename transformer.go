package mpipe

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
