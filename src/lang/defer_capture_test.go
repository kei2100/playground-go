package lang

import (
	"fmt"
	"testing"
)

type myCloser struct {
	message string
}

func (c *myCloser) Close() {
	fmt.Println(c.message)
}

func TestDeferCaptureVariable(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		c := &myCloser{message: "a"}
		defer c.Close()
		c = &myCloser{message: "aaa"}
		// Prints `a`
		// defer に続く関数の評価時にキャプチャされる
	})
	t.Run("2", func(t *testing.T) {
		c := &myCloser{message: "b"}
		defer func() {
			c.Close()
		}()
		c = &myCloser{message: "bbb"}
		// Prints `bbb`
		// 無名関数の評価時では c はキャプチャ対象にならない
	})
}
