package mpipe

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type transformer func(string) string

var noTransform transformer = func(s string) string {
	return s
}

func transform(w io.Writer, r io.Reader, trasn transformer) error {
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
