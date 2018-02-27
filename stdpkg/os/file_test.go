package os

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileMode(t *testing.T) {
	t.Run("IsRegular()", func(t *testing.T) {
		info, err := os.Stdout.Stat()
		if err != nil {
			t.Fatal(err)
		}
		if g, w := info.Mode().IsRegular(), false; g != w {
			t.Errorf("Stdout IsRegular got %v, want %v", g, w)
		}

		f, err := ioutil.TempFile("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		info, err = f.Stat()
		if err != nil {
			t.Fatal(err)
		}
		if g, w := info.Mode().IsRegular(), true; g != w {
			t.Errorf("TempFile IsRegular got %v, want %v", g, w)
		}
	})
}
