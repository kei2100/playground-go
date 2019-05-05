package stat

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestSameFile(t *testing.T) {
	ds, teardown := setup()
	defer teardown()

	f := createFile(ds, "f")
	defer f.Close()
	stat := statFile(f)

	f.Close()
	os.Rename(f.Name(), f.Name()+".bk")
	renamedStat := statFile(openFile(f.Name() + ".bk"))

	newf := createFile(ds, "f")
	defer newf.Close()
	newstat := statFile(newf)

	if !SameFile(stat, renamedStat) {
		t.Errorf("stat renamedStat are the not same")
	}
	if SameFile(stat, newstat) {
		t.Errorf("stat newstat are the same")
	}
}

type dataSet struct {
	tempDir string
}

func setup() (ds *dataSet, teardown func()) {
	tempDir, err := ioutil.TempDir("", "follow-")
	if err != nil {
		panic(err)
	}
	return &dataSet{tempDir: tempDir}, func() {
		os.RemoveAll(tempDir)
	}
}

func createFile(ds *dataSet, name string) *onceCloseFile {
	f, err := os.OpenFile(filepath.Join(ds.tempDir, name), os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	return &onceCloseFile{File: f}
}

func openFile(name string) *onceCloseFile {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	return &onceCloseFile{File: f}
}

func statFile(f *onceCloseFile) *FileStat {
	s, err := Stat(f.File)
	if err != nil {
		panic(err)
	}
	return s
}

type onceCloseFile struct {
	once sync.Once
	*os.File
}

func (f *onceCloseFile) Close() error {
	var err error
	f.once.Do(func() {
		err = f.File.Close()
	})
	return err
}
