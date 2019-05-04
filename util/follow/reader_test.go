package follow

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kei2100/playground-go/util/follow/file"
	"github.com/kei2100/playground-go/util/follow/posfile"
)

func TestNoPositionFile(t *testing.T) {
	t.Run("Glow", func(t *testing.T) {
		t.Parallel()

		ds, teardown := setup()
		defer teardown()

		ds.logFile.WriteString("foo")
		wantRead(t, ds.reader, "fo")
		wantPositionFile(t, ds.reader, ds.logFile, 2)

		wantRead(t, ds.reader, "o")
		wantReadAll(t, ds.reader, "")
		wantPositionFile(t, ds.reader, ds.logFile, 3)

		ds.logFile.WriteString("bar")
		wantReadAll(t, ds.reader, "bar")
		wantPositionFile(t, ds.reader, ds.logFile, 6)
	})

	t.Run("Follow Rotate", func(t *testing.T) {
		t.Parallel()

		ds, teardown := setup(WithWatchRotateInterval(10*time.Millisecond), WithDetectRotateDelay(0))
		defer teardown()

		rotateLogFile(ds.logFile)
		wantDetectRotate(t, ds.reader, 500*time.Millisecond)

		ds.logFile.WriteString("foo")
		wantReadAll(t, ds.reader, "foo")
		wantPositionFile(t, ds.reader, ds.logFile, 3)
	})

	t.Run("No Follow Rotate", func(t *testing.T) {
		t.Parallel()

		ds, teardown := setup(WithWatchRotateInterval(10*time.Millisecond), WithDetectRotateDelay(0), WithFollowRotate(false))
		defer teardown()

		bkLogFile, err := file.Open(ds.logFile.Name())
		if err != nil {
			t.Errorf("failed to open %v", ds.logFile.Name())
			return
		}
		defer bkLogFile.Close()

		rotateLogFile(ds.logFile)
		wantNoDetectRotate(t, ds.reader, 500*time.Millisecond)

		ds.logFile.WriteString("foo")
		wantReadAll(t, ds.reader, "")
		wantPositionFile(t, ds.reader, bkLogFile, 0)
	})

	t.Run("Follow Rotate DetectRotateDelay", func(t *testing.T) {
		t.Parallel()

		ds, teardown := setup(WithWatchRotateInterval(10*time.Millisecond), WithDetectRotateDelay(500*time.Millisecond))
		defer teardown()

		bkLogFile, err := file.Open(ds.logFile.Name())
		if err != nil {
			t.Errorf("failed to open %v", ds.logFile.Name())
			return
		}
		defer bkLogFile.Close()

		ds.logFile.WriteString("foo")
		rotateLogFile(ds.logFile)
		wantReadAll(t, ds.reader, "foo")
		wantPositionFile(t, ds.reader, bkLogFile, 3)

		wantDetectRotate(t, ds.reader, time.Second)
		ds.logFile.WriteString("bar")
		wantReadAll(t, ds.reader, "bar")
		wantPositionFile(t, ds.reader, ds.logFile, 3)
	})
}

func TestWithPositionFile(t *testing.T) {
	t.Run("Works", func(t *testing.T) {
		t.Parallel()

		logFile, fileInfo := createLogFile()
		logFile.WriteString("bar")
		positionFile := posfile.NewMemoryPositionFile(fileInfo, 2)
		ds, teardown := setupWithLogFile(logFile, WithPositionFile(positionFile))
		defer teardown()

		wantReadAll(t, ds.reader, "r")
		wantPositionFile(t, ds.reader, ds.logFile, 3)

		ds.logFile.WriteString("baz")
		wantReadAll(t, ds.reader, "baz")
		wantPositionFile(t, ds.reader, ds.logFile, 6)
	})

	t.Run("Incorrect offset", func(t *testing.T) {
		t.Parallel()

		logFile, fileInfo := createLogFile()
		logFile.WriteString("bar")
		positionFile := posfile.NewMemoryPositionFile(fileInfo, 4)
		ds, teardown := setupWithLogFile(logFile, WithPositionFile(positionFile))
		defer teardown()

		wantReadAll(t, ds.reader, "")
		wantPositionFile(t, ds.reader, ds.logFile, 3)

		ds.logFile.WriteString("baz")
		wantReadAll(t, ds.reader, "baz")
		wantPositionFile(t, ds.reader, ds.logFile, 6)
	})

	t.Run("Same file not found", func(t *testing.T) {
		t.Parallel()

		logFile, fileInfo := createLogFile()
		logFile.WriteString("bar")
		rotateLogFile(logFile)

		positionFile := posfile.NewMemoryPositionFile(fileInfo, 2)
		newLogFile, err := os.OpenFile(filepath.Join(filepath.Dir(logFile.Name()), fileInfo.Name()), os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		ds, teardown := setupWithLogFile(newLogFile, WithPositionFile(positionFile))
		defer teardown()

		wantReadAll(t, ds.reader, "")
		wantPositionFile(t, ds.reader, ds.logFile, 0)

		ds.logFile.WriteString("baz")
		wantReadAll(t, ds.reader, "baz")
		wantPositionFile(t, ds.reader, ds.logFile, 3)
	})
}

// use saved positionfile

type dataset struct {
	logFile *os.File
	reader  *reader
}

func createLogFile() (*os.File, os.FileInfo) {
	tempDir, err := ioutil.TempDir("", "follow-")
	if err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile(filepath.Join(tempDir, "test.log"), os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	fileInfo, err := logFile.Stat()
	if err != nil {
		panic(err)
	}
	return logFile, fileInfo
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
	logFile, _ := createLogFile()
	return setupWithLogFile(logFile, opts...)
}

func setupWithLogFile(logFile *os.File, opts ...OptionFunc) (ds *dataset, teardown func()) {
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
		os.Remove(logFile.Name())
		os.Remove(logFile.Name() + ".1") // See rotateLogFile(*os.File)
	}
	return &dataset{logFile: logFile, reader: reader}, teardown
}

func wantPositionFile(t *testing.T, reader *reader, wantFileInfoFile *os.File, wantOffset int64) {
	t.Helper()

	wantFileInfo, err := wantFileInfoFile.Stat()
	if err != nil {
		t.Errorf("failed to get fileInfo: %v", err)
		return
	}
	if !os.SameFile(reader.positionFile.FileInfo(), wantFileInfo) {
		t.Errorf("fileInfo not same")
	}
	if g, w := reader.positionFile.Offset(), wantOffset; g != w {
		t.Errorf("offset got %v, want %v", g, w)
	}
}

func wantRead(t *testing.T, reader *reader, want string) {
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

func wantNoDetectRotate(t *testing.T, reader *reader, wait time.Duration) {
	t.Helper()

	select {
	case <-reader.rotated:
		t.Errorf("detect rotate. want not detect")
	case <-time.After(wait):
		return
	}
}
