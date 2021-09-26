package filepath

import (
	"path/filepath"
	"testing"
)

func TestCurDir(t *testing.T) {
	g := CurDir()
	w, _ := filepath.Abs("testdata/..")
	if g != w {
		t.Errorf("CurDir got %s, want %s", g, w)
	}
}
