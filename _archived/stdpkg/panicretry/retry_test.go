package panicretry_test

import (
	"errors"
	"testing"

	"github.com/kei2100/playground-go/stdpkg/panicretry"
)

func TestDo(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		var someerr = errors.New("someerr")
		got := panicretry.Do(func() error {
			return someerr
		})
		if g, w := got, someerr; g != w {
			t.Errorf("err got %v, want %v", g, w)
		}
	})
	t.Run("2", func(t *testing.T) {
		got := panicretry.Do(func() error {
			return nil
		})
		if g, w := got, error(nil); g != w {
			t.Errorf("err got %v, want %v", g, w)
		}
	})
	t.Run("3", func(t *testing.T) {
		counter := 0
		got := panicretry.Do(func() error {
			if counter < 3 {
				counter++
				panic("omg")
			}
			return nil
		})
		if g, w := got, error(nil); g != w {
			t.Errorf("err got %v, want %v", g, w)
		}
		if g, w := counter, 3; g != w {
			t.Errorf("couter got %v, want %v", g, w)
		}
	})
}
