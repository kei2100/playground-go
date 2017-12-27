package errors

import (
	"bytes"
	"io"
	"testing"
)

type HoldErrorWriter struct {
	w   io.Writer
	err error
}

func (hw *HoldErrorWriter) Write(p []byte) int {
	if hw.err != nil {
		return 0
	}
	n, err := hw.w.Write(p)
	hw.err = err
	return n
}

func (hw *HoldErrorWriter) Err() error {
	return hw.err
}

func TestHoldErrorWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	hw := &HoldErrorWriter{
		w: buf,
	}

	hw.Write([]byte("foo"))
	hw.Write([]byte("bar"))

	if g, w := buf.String(), "foobar"; g != w {
		t.Errorf(" got %v, want %v", g, w)
	}
	if err := hw.Err(); err != nil {
		t.Errorf(" got %v, want nil", err)
	}
}
