package io

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
