package posfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/kei2100/playground-go/util/follow/stat"
)

func TestOpenUpdate(t *testing.T) {
	dir, err := ioutil.TempDir("", "follow-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	someFile, err := os.OpenFile(filepath.Join(dir, "somefile"), os.O_CREATE, 0600)
	if err != nil {
		t.Error(err)
		return
	}
	someFileStat, err := stat.Stat(someFile)
	if err != nil {
		t.Error(err)
		return
	}

	pf, err := Open(filepath.Join(dir, "posfile"))
	if err != nil {
		t.Errorf("failed to open posfile: %+v", err)
		return
	}
	pf = &onceClose{PositionFile: pf}
	defer pf.Close()

	pf.Update(someFileStat, 0)
	pf.IncreaseOffset(2)
	if !stat.SameFile(pf.FileStat(), someFileStat) {
		t.Errorf("not same fileStat\ngot: \n%+v\nwant: \n%+v", pf.FileStat(), someFileStat)
	}
	if g, w := pf.Offset(), int64(2); g != w {
		t.Errorf("offset got %v, want %v", g, w)
	}
	if err := pf.Close(); err != nil {
		t.Errorf("failed to close: %+v", err)
		return
	}

	pf2, err := Open(filepath.Join(dir, "posfile"))
	if err != nil {
		t.Errorf("failed to open posfile: %+v", err)
		return
	}
	pf2 = &onceClose{PositionFile: pf2}
	defer pf2.Close()

	if !stat.SameFile(pf2.FileStat(), someFileStat) {
		t.Errorf("not same fileStat\ngot: \n%+v\nwant: \n%+v", pf2.FileStat(), someFileStat)
	}
	if g, w := pf2.Offset(), int64(2); g != w {
		t.Errorf("offset got %v, want %v", g, w)
	}
	if err := pf2.Close(); err != nil {
		t.Errorf("failed to close: %+v", err)
		return
	}
}

type onceClose struct {
	once sync.Once
	PositionFile
}

func (oc *onceClose) Close() error {
	var err error
	oc.once.Do(func() {
		err = oc.PositionFile.Close()
	})
	return err
}
