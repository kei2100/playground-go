package io

import (
	"io"
	"io/ioutil"
	"testing"
)

func TestPipe(t *testing.T) {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		_, err := pw.Write([]byte("hello"))
		if err != nil {
			panic(err)
		}
	}()

	defer pr.Close()
	got, err := ioutil.ReadAll(pr)
	if err != nil {
		panic(err)
	}

	if g, w := string(got), "hello"; g != w {
		t.Errorf(" got %v, want %v", g, w)
	}
}
