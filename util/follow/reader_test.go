package follow

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNoPositionFile(t *testing.T) {
	t.Run("Glow", func(t *testing.T) {
		ds, teardown := setup()
		defer teardown()

		ds.logFile.WriteString("foo")
		wantRead(t, ds.reader, "fo")
		wantRead(t, ds.reader, "o")
		wantReadAll(t, ds.reader, "")

		ds.logFile.WriteString("bar")
		wantReadAll(t, ds.reader, "bar")
	})

	t.Run("Rotate", func(t *testing.T) {
		ds, teardown := setup(WithWatchRotateInterval(200*time.Millisecond), WithDetectRotateDelay(0))
		defer teardown()

		rotateLogFile(ds.logFile)
		wantDetectRotate(t, ds.reader, time.Second)

		ds.logFile.WriteString("foo")
		wantReadAll(t, ds.reader, "foo")
	})
}

type dataset struct {
	logFile *os.File
	reader  *reader
}

func rotateLogFile(logFile *os.File) {
	logFile.Close()
	if err := os.Rename(logFile.Name(), logFile.Name()+".1"); err != nil {
		panic(err)
	}
	newLogFile, err := os.OpenFile(logFile.Name(), os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	*logFile = *newLogFile
}

func setup(opts ...OptionFunc) (ds *dataset, teardown func()) {
	tempDir, err := ioutil.TempDir("", "follow-")
	if err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile(filepath.Join(tempDir, "test.log"), os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	r, err := Open(logFile.Name(), opts...)
	if err != nil {
		panic(err)
	}
	reader, ok := r.(*reader)
	if !ok {
		panic("failed to cast")
	}

	teardown = func() {
		logFile.Close()
		reader.Close()
		os.RemoveAll(tempDir)
	}
	return &dataset{logFile: logFile, reader: reader}, teardown
}

func wantRead(t *testing.T, reader Reader, want string) {
	t.Helper()

	b := make([]byte, len(want))
	n, err := reader.Read(b)
	if err != nil {
		t.Errorf("failed to read: %v", err)
		return
	}
	if g, w := n, len(b); g != w {
		t.Errorf("nReadBytes got %v, want %v", g, w)
	}
	if g, w := string(b), want; g != w {
		t.Errorf("byteString got %v, want %v", g, w)
	}
}

func wantReadAll(t *testing.T, reader *reader, want string) {
	t.Helper()

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Errorf("failed to read all: %v", err)
		return
	}
	if g, w := len(b), len(want); g != w {
		t.Errorf("nReadBytes got %v, want %v", g, w)
	}
	if g, w := string(b), want; g != w {
		t.Errorf("byteString got %v, want %v", g, w)
	}
}

func wantDetectRotate(t *testing.T, reader *reader, timeout time.Duration) {
	t.Helper()

	select {
	case <-reader.rotated:
		return
	case <-time.After(timeout):
		t.Errorf("%s timeout while waiting for detect rotate", timeout)
	}
}
