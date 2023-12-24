package mpipe

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Transformer func(string) string

var NoTransform Transformer = func(s string) string {
	return s
}

func transform(w io.Writer, r io.Reader, trasn Transformer) error {
	s := bufio.NewScanner(r)

	for s.Scan() {
		t := s.Text()
		t = stripNewLines(t)
		t = trasn(t)

		_, err := io.Copy(w, bytes.NewBuffer([]byte(t)))
		if err != nil {
			return err
		}
	}

	return nil
}

func stripNewLines(s string) string {
	var ok bool
	if s, ok = strings.CutSuffix(s, "\n"); ok {
		s, _ = strings.CutSuffix(s, "\r")
	}
	return s
}
