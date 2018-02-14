package io

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestWriterDecorator(t *testing.T) {
	t.Parallel()

	// File Writer
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	// Buffered Writer (buffer size is default)
	bw := bufio.NewWriter(f)
	defer bw.Flush()

	// Multiple Writer
	mw := io.MultiWriter(bw, os.Stdout)

	// dump http request
	req, err := http.NewRequest("GET", "https://www.google.co.jp", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Write(mw)
}

func TestReaderDecorator(t *testing.T) {
	t.Parallel()

	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(f.Name())
	}()
	if _, err := f.Write([]byte("foo bar ")); err != nil {
		t.Fatal(err)
	}
	f.Close()

	f, err = os.Open(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Buffered Reader (buffer size is default)
	br := bufio.NewReader(f)

	// Concat reader's
	mr := io.MultiReader(br, strings.NewReader("baz"))

	// tee
	teebf := new(bytes.Buffer)
	tr := io.TeeReader(mr, teebf)

	result, err := ioutil.ReadAll(tr)
	if err != nil {
		t.Fatal(err)
	}
	if g, w := string(result), "foo bar baz"; g != w {
		t.Errorf("result got %v, want %v", g, w)
	}
	if g, w := teebf.String(), string(result); g != w {
		t.Errorf(" got %v, want %v", g, w)
	}
}
