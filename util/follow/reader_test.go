package follow

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNoPositionFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "tailf")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(tempDir)

	logFile, err := os.OpenFile(filepath.Join(tempDir, "test.log"), os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Error(err)
		return
	}

	fwi, err := Open(logFile.Name())
	if err != nil {
		t.Error(err)
		return
	}
	fw := fwi.(*reader)
	fw.detectRotateDelay = 100 * time.Millisecond

	testRead(t, fw, logFile, 0)

	t.Run("logFile was rotated", func(t *testing.T) {
		logFile.Close()
		os.Rename(logFile.Name(), logFile.Name()+".1")
		logFile, err = os.OpenFile(filepath.Join(tempDir, "test.log"), os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			t.Error(err)
			return
		}
		testDetectRotate(t, fw, time.Second)
		testRead(t, fw, logFile, 0)
	})
}

func testDetectRotate(t *testing.T, follower *reader, timeout time.Duration) {
	t.Helper()

	select {
	case <-follower.rotated:
		return
	case <-time.After(timeout):
		t.Errorf("%s timeout while waiting for detect rotate", timeout)
	}
}

func testRead(t *testing.T, follower *reader, logFile *os.File, offset int64) {
	t.Helper()

	// write foo
	_, err := logFile.WriteString("foo")
	if err != nil {
		t.Errorf("write foo: %v", err)
	}

	// read 2 bytes
	b := make([]byte, 2)
	n, err := follower.Read(b)
	if err != nil {
		t.Errorf("read fo: %v", err)
	}
	if g, w := n, 2; g != w {
		t.Errorf("read fo: n read bytes got %v, want %v", g, w)
	}
	if g, w := follower.positionFile.Offset(), offset+2; g != w {
		t.Errorf("read fo: offset got %v, want %v", g, w)
	}
	if g, w := string(b), "fo"; g != w {
		t.Errorf("read fo: byteString got %v, want %v", g, w)
	}

	// read 2 bytes
	b = make([]byte, 2)
	n, err = follower.Read(b)
	if err != nil {
		t.Errorf("read o: %v", err)
	}
	if g, w := n, 1; g != w {
		t.Errorf("read o: n read bytes got %v, want %v", g, w)
	}
	if g, w := follower.positionFile.Offset(), offset+3; g != w {
		t.Errorf("read o: offset got %v, want %v", g, w)
	}
	if g, w := string(b[:1]), "o"; g != w {
		t.Errorf("read o: byteString got %v, want %v", g, w)
	}

	// append bar
	_, err = logFile.WriteString("bar")
	if err != nil {
		t.Errorf("write bar: %v", err)
	}

	// read 2 bytes
	b = make([]byte, 2)
	n, err = follower.Read(b)
	if err != nil {
		t.Errorf("read ba: %v", err)
	}
	if g, w := n, 2; g != w {
		t.Errorf("read ba: n read bytes got %v, want %v", g, w)
	}
	if g, w := follower.positionFile.Offset(), offset+5; g != w {
		t.Errorf("read ba: offset got %v, want %v", g, w)
	}
	if g, w := string(b), "ba"; g != w {
		t.Errorf("read ba: byteString got %v, want %v", g, w)
	}

	// read all
	b, err = ioutil.ReadAll(follower)
	if err != nil {
		t.Errorf("read all: %v", err)
	}
	if g, w := follower.positionFile.Offset(), offset+6; g != w {
		t.Errorf("read all: offset got %v, want %v", g, w)
	}
	if g, w := string(b), "r"; g != w {
		t.Errorf("read all: byteString got %v, want %v", g, w)
	}
}
